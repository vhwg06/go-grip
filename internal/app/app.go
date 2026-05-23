// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/internal/controller/restapi"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/evrone/go-clean-template/internal/repo/webapi"
	"github.com/evrone/go-clean-template/internal/usecase/cart"
	"github.com/evrone/go-clean-template/internal/usecase/catalog"
	"github.com/evrone/go-clean-template/internal/usecase/content"
	"github.com/evrone/go-clean-template/internal/usecase/importer"
	"github.com/evrone/go-clean-template/internal/usecase/lead"
	"github.com/evrone/go-clean-template/internal/usecase/media"
	"github.com/evrone/go-clean-template/internal/usecase/notification"
	"github.com/evrone/go-clean-template/internal/usecase/task"
	"github.com/evrone/go-clean-template/internal/usecase/translation"
	"github.com/evrone/go-clean-template/internal/usecase/user"
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
	media       *media.UseCase
	homepage    *content.HomepageUseCase
	cart        *cart.UseCase
	lead        *lead.UseCase
	content     *content.UseCase
	importer    *importer.UseCase
}

type servers struct {
	http *httpserver.Server
}

func initUseCases(cfg *config.Config, pg *postgres.Postgres, jwtManager *jwt.Manager) useCases {
	userRepo := persistent.NewUserRepo(pg)
	taskRepo := persistent.NewTaskRepo(pg)
	translationRepo := persistent.NewTranslationRepo(pg)
	catalogRepo := persistent.NewCatalogRepo(pg)
	mediaRepo := persistent.NewMediaRepo(pg)
	homepageRepo := persistent.NewHomepageRepo(pg)
	supportRepo := persistent.NewSupportChannelRepo(pg)
	cartRepo := persistent.NewCartRepo(pg)
	orderRepo := persistent.NewOrderRequestRepo(pg)
	leadRepo := persistent.NewLeadRepo(pg)
	contentRepo := persistent.NewContentRepo(pg)
	importRepo := persistent.NewImportRepo(pg, catalogRepo, contentRepo)
	notificationUseCase := notification.New(cfg.Notification.Enabled)

	return useCases{
		user:        user.New(userRepo, jwtManager),
		task:        task.New(taskRepo),
		translation: translation.New(translationRepo, webapi.New()),
		catalog:     catalog.New(catalogRepo),
		media:       media.New(mediaRepo, cfg.Ecommerce.MediaMaxBytes),
		homepage:    content.NewHomepage(homepageRepo, supportRepo),
		cart:        cart.New(cartRepo, orderRepo, notificationUseCase),
		lead:        lead.New(leadRepo),
		content:     content.New(contentRepo),
		importer:    importer.New(importRepo, cfg.Ecommerce.InitialImportMax),
	}
}

func initServers(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l logger.Interface) servers {
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	restapi.NewRouter(httpServer.App, cfg, uc.translation, uc.user, uc.task, uc.catalog, uc.media, uc.homepage, uc.cart, uc.lead, uc.content, uc.importer, jwtManager, l)

	return servers{http: httpServer}
}

func (s *servers) startServers() {
	s.http.Start()
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
