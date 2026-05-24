package restapi

import (
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/evrone/go-clean-template/config"
	_ "github.com/evrone/go-clean-template/docs" // Swagger docs.
	"github.com/evrone/go-clean-template/internal/controller/restapi/middleware"
	v1 "github.com/evrone/go-clean-template/internal/controller/restapi/v1"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// NewRouter -.
// Swagger spec:
//
//	@title       Grip Store Backend REST API
//	@description REST-only backend for catalog, checkout, orders, profile, admin, and notifications
//	@version     1.0
//	@host        localhost:8080
//	@BasePath    /v1
//	@securityDefinitions.apikey BearerAuth
//	@in header
//	@name Authorization
func NewRouter(app *fiber.App, cfg *config.Config, t usecase.Translation, u usecase.User, tk usecase.Task, catalog usecase.Catalog, auth usecase.Auth, checkout usecase.Checkout, orders usecase.Orders, profile usecase.Profile, admin usecase.Admin, wishlist usecase.Wishlist, notify usecase.NotificationCenter, media usecase.Media, homepage usecase.Homepage, cart usecase.Cart, lead usecase.Lead, content usecase.Content, importer usecase.Importer, jwtManager *jwt.Manager, l logger.Interface) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))
	app.Use(middleware.RejectBlockedMutations())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.HTTP.CORSAllowedOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Playwright-Test",
		AllowCredentials: cfg.HTTP.CORSAllowCredentials,
	}))


	// Prometheus metrics
	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("my-service-name")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// K8s probe
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewRoutes(apiV1Group, t, u, tk, catalog, auth, checkout, orders, profile, admin, wishlist, notify, media, homepage, cart, lead, content, importer, jwtManager, cfg.Admin.Users, l)
	}
}
