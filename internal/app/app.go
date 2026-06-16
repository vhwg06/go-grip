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
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/evrone/go-clean-template/internal/repo/webapi"
	"github.com/evrone/go-clean-template/internal/usecase"
	adminuc "github.com/evrone/go-clean-template/internal/usecase/admin"
	"github.com/evrone/go-clean-template/internal/usecase/auth"
	"github.com/evrone/go-clean-template/internal/usecase/cart"
	"github.com/evrone/go-clean-template/internal/usecase/catalog"
	"github.com/evrone/go-clean-template/internal/usecase/checkout"
	"github.com/evrone/go-clean-template/internal/usecase/content"
	"github.com/evrone/go-clean-template/internal/usecase/importer"
	"github.com/evrone/go-clean-template/internal/usecase/lead"
	"github.com/evrone/go-clean-template/internal/usecase/media"
	"github.com/evrone/go-clean-template/internal/usecase/notification"
	"github.com/evrone/go-clean-template/internal/usecase/orders"
	"github.com/evrone/go-clean-template/internal/usecase/profile"
	"github.com/evrone/go-clean-template/internal/usecase/task"
	"github.com/evrone/go-clean-template/internal/usecase/translation"
	"github.com/evrone/go-clean-template/internal/usecase/user"
	"github.com/evrone/go-clean-template/internal/usecase/wishlist"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type useCases struct {
	translation *translation.UseCase
	user        *user.UseCase
	task        *task.UseCase
	catalog     *catalog.UseCase
	auth        *auth.UseCase
	checkout    *checkout.UseCase
	orders      *orders.UseCase
	profile     *profile.UseCase
	admin       *adminuc.UseCase
	maintenance *adminuc.MaintenanceUseCase
	wishlist    *wishlist.UseCase
	notify      *notification.CenterUseCase
	media       *media.UseCase
	homepage    *content.HomepageUseCase
	cart        *cart.UseCase
	lead        *lead.UseCase
	content     *content.UseCase
	importer    *importer.UseCase
}

type servers struct {
	http              *httpserver.Server
	maintenance       *adminuc.MaintenanceUseCase
	maintenanceTicker *time.Ticker
	maintenanceDone   chan struct{}
}

func initUseCases(cfg *config.Config, pg *postgres.Postgres, jwtManager *jwt.Manager) useCases {
	userRepo := persistent.NewUserRepo(pg)
	taskRepo := persistent.NewTaskRepo(pg)
	translationRepo := persistent.NewTranslationRepo(pg)
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
	notificationUseCase := notification.New(cfg.Notification.Enabled)
	epayVerifier := webapi.NewEpayVerifier(cfg.Payment.SecretKey)
	adminNotifier := webapi.NewNoopAdminNotifier()

	checkoutUC := checkout.New(gripCheckoutRepo, gripOrderRepo)
	checkoutUC.SetPaymentVerifier(epayVerifier)

	var mediaStorage usecase.MediaStorage
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

	return useCases{
		user:        user.New(userRepo, jwtManager),
		task:        task.New(taskRepo),
		translation: translation.New(translationRepo, webapi.New()),
		catalog:     catalog.NewWithGrip(catalogRepo, gripCatalogRepo),
		auth:        auth.New(authRepo, jwtManager, 30*24*time.Hour, cfg.Admin.Users),
		checkout:    checkoutUC,
		orders:      orders.New(gripOrderRepo),
		profile:     profile.New(profileRepo, 10),
		admin:       adminuc.New(adminRepo, adminNotifier, cfg.Admin.Users),
		maintenance: adminuc.NewMaintenance(maintenanceRepo, 5*time.Minute),
		wishlist:    wishlist.New(wishlistRepo, gripOrderRepo),
		notify:      notification.NewCenter(notificationRepo),
		media:       media.New(mediaRepo, mediaStorage, media.Config{
			MaxBytes:    cfg.Ecommerce.MediaMaxBytes,
		}),
		homepage:    content.NewHomepage(homepageRepo, supportRepo),
		cart:        cart.New(cartRepo, orderRepo, notificationUseCase),
		lead:        lead.New(leadRepo),
		content:     content.New(contentRepo),
		importer:    importer.New(importRepo, cfg.Ecommerce.InitialImportMax),
	}
}

func initServers(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l logger.Interface) servers {
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	restapi.NewRouter(httpServer.App, cfg, uc.translation, uc.user, uc.task, uc.catalog, uc.auth, uc.checkout, uc.orders, uc.profile, uc.admin, uc.wishlist, uc.notify, uc.media, uc.homepage, uc.cart, uc.lead, uc.content, uc.importer, jwtManager, l)

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
				_ = s.maintenance.CleanupExpiredCards(ctx)
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
