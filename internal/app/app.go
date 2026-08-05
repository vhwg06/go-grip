// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/internal/controller/restapi"
	cartmodule "github.com/evrone/go-clean-template/internal/module/cart"
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	"github.com/evrone/go-clean-template/internal/module/catalog/catalogbase"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	importermodule "github.com/evrone/go-clean-template/internal/module/importer"
	leadmodule "github.com/evrone/go-clean-template/internal/module/lead"
	mediamodule "github.com/evrone/go-clean-template/internal/module/media"
	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/evrone/go-clean-template/internal/repo/webapi"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type cartNotificationAdapter struct {
	uc notificationmodule.NotificationUseCase
}

func (a cartNotificationAdapter) Dispatch(ctx context.Context, channel, to, subject string) error {
	if a.uc == nil {
		return nil
	}
	return a.uc.Dispatch(ctx, notificationmodule.Notification{Channel: channel, To: to, Subject: subject})
}

type useCases struct {
	user        usermodule.UserUseCase
	catalog     catalogmodule.CatalogUseCase
	catalogBase catalogbase.UseCase
	auth        usermodule.AuthUseCase
	checkout    ordermodule.CheckoutUseCase
	orders      ordermodule.OrdersUseCase
	profile     usermodule.ProfileUseCase
	admin       usermodule.AdminUseCase
	maintenance *usermodule.MaintenanceUseCase
	wishlist    wishlistmodule.WishlistUseCase
	notify      notificationmodule.NotificationCenterUseCase
	media       mediamodule.MediaUseCase
	homepage    contentmodule.HomepageUseCase
	cart        cartmodule.CartUseCase
	lead        leadmodule.LeadUseCase
	content     contentmodule.ContentUseCase
	importer    importermodule.ImporterUseCase
}

type servers struct {
	http              *httpserver.Server
	maintenance       *usermodule.MaintenanceUseCase
	maintenanceTicker *time.Ticker
	maintenanceDone   chan struct{}
}

func initUseCases(cfg *config.Config, pg *postgres.Postgres, jwtManager *jwt.Manager) useCases {
	userRepo := persistent.NewUserRepo(pg)
	catalogRepo := persistent.NewCatalogRepo(pg)
	authRepo := persistent.NewAuthRepo(pg)
	profileRepo := persistent.NewProfileRepo(pg)
	gripCatalogRepo := persistent.NewGripCatalogRepo(pg)
	gripCheckoutRepo := persistent.NewCheckoutRepo(pg)
	gripOrderRepo := persistent.NewGripOrderRepo(pg)
	adminRepo := persistent.NewAdminRepo(pg)
	maintenanceRepo := persistent.NewMaintenanceRepo(pg)
	wishlistRepo := persistent.NewWishlistRepo(pg)
	notificationRepo := persistent.NewNotificationRepo(pg)
	mediaRepo := persistent.NewMediaRepo(pg)
	homepageRepo := persistent.NewHomepageRepo(pg)
	supportRepo := persistent.NewSupportChannelRepo(pg)
	cartRepo := persistent.NewCartRepo(pg)
	orderRepo := persistent.NewOrderRequestRepo(pg)
	leadRepo := persistent.NewLeadRepo(pg)
	contentRepo := persistent.NewContentRepo(pg)
	importRepo := persistent.NewImportRepo(pg, catalogRepo, contentRepo)
	notificationUseCase := notificationmodule.NewNotificationUseCase(cfg.Notification.Enabled)
	epayVerifier := webapi.NewEpayVerifier(cfg.Payment.SecretKey)
	adminNotifier := webapi.NewNoopAdminNotifier()

	checkoutUC := ordermodule.NewCheckoutUseCase(gripCheckoutRepo, gripOrderRepo)
	checkoutUC.SetPaymentVerifier(epayVerifier)

	var mediaStorage mediamodule.MediaStorageProvider
	if cfg.R2.AccountID != "" && cfg.R2.AccessKeyID != "" && cfg.R2.SecretKey != "" && cfg.R2.BucketName != "" {
		mediaStorage = webapi.NewR2Storage(
			cfg.R2.AccountID,
			cfg.R2.AccessKeyID,
			cfg.R2.SecretKey,
			cfg.R2.BucketName,
			cfg.R2.PublicURL,
		)
	} else {
		mediaStorage = webapi.NewLocalStorage(cfg.App.BaseURL)
	}

	catalogUseCase := catalogmodule.NewCatalogUseCase(catalogRepo, gripCatalogRepo)
	catalogBaseRepositories := persistent.NewCatalogRepositories(pg)
	catalogBaseUnitOfWork := persistent.NewUnitOfWork(pg)
	catalogBaseUseCase := catalogbase.New(catalogBaseRepositories, catalogBaseUnitOfWork)

	return useCases{
		user:        usermodule.NewUserUseCase(userRepo, jwtManager),
		catalog:     catalogUseCase,
		catalogBase: catalogBaseUseCase,
		auth:        usermodule.NewAuthUseCase(authRepo, jwtManager, 30*24*time.Hour, cfg.Admin.Users),
		checkout:    checkoutUC,
		orders:      ordermodule.NewOrdersUseCase(gripOrderRepo),
		profile:     usermodule.NewProfileUseCase(profileRepo),
		admin:       usermodule.NewAdminUseCase(adminRepo, adminNotifier, cfg.Admin.Users),
		maintenance: usermodule.NewMaintenanceUseCase(maintenanceRepo, 5*time.Minute),
		wishlist:    wishlistmodule.NewWishlistUseCase(wishlistRepo, gripOrderRepo),
		notify:      notificationmodule.NewNotificationCenterUseCase(notificationRepo),
		media: mediamodule.NewMediaUseCase(mediaRepo, mediaStorage, mediamodule.Config{
			MaxBytes: cfg.Ecommerce.MediaMaxBytes,
		}),
		homepage: contentmodule.NewHomepageUseCase(homepageRepo, supportRepo),
		cart:     cartmodule.NewCartUseCase(cartRepo, orderRepo, cartNotificationAdapter{notificationUseCase}),
		lead:     leadmodule.NewLeadUseCase(leadRepo),
		content:  contentmodule.NewContentUseCase(contentRepo),
		importer: importermodule.NewImporterUseCase(importRepo, cfg.Ecommerce.InitialImportMax),
	}
}

func initServers(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l logger.Interface) servers {
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	restapi.NewRouter(httpServer.App, cfg, uc.user, uc.catalog, uc.catalogBase, uc.auth, uc.checkout, uc.orders, uc.profile, uc.admin, uc.wishlist, uc.notify, uc.media, uc.homepage, uc.cart, uc.lead, uc.content, uc.importer, jwtManager, l)

	interval := cfg.Ecommerce.SchedulerInterval
	if interval <= 0 {
		interval = time.Minute
	}

	return servers{
		http:              httpServer,
		maintenance:       uc.maintenance,
		maintenanceTicker: time.NewTicker(interval),
		maintenanceDone:   make(chan struct{}),
	}
}

func (s *servers) startServers() {
	s.http.Start()
	s.startMaintenance()
}

func (s *servers) startMaintenance() {
	if s.maintenance == nil || s.maintenanceTicker == nil {
		return
	}

	go func() {
		for {
			select {
			case <-s.maintenanceDone:
				return
			case <-s.maintenanceTicker.C:
				ctx := context.Background()
				_ = s.maintenance.CancelExpiredPendingOrders(ctx)
				_ = s.maintenance.SyncProductAggregates(ctx)
			}
		}
	}()
}

func (s *servers) waitForShutdown(l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal: %s", sig.String())
	case err = <-s.http.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	s.shutdownServers(l)
}

func (s *servers) shutdownServers(l logger.Interface) {
	if s.maintenanceTicker != nil {
		s.maintenanceTicker.Stop()
	}
	if s.maintenanceDone != nil {
		close(s.maintenanceDone)
	}

	if err := s.http.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	jwtManager := jwt.New(cfg.JWT.Secret, cfg.JWT.TokenExpiry)

	uc := initUseCases(cfg, pg, jwtManager)
	s := initServers(cfg, uc, jwtManager, l)
	s.startServers()
	s.waitForShutdown(l)
}
