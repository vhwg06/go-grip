package restapi

import (
	"net/http"
	"os"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/config"
	_ "github.com/evrone/go-clean-template/docs" // Swagger docs.
	"github.com/evrone/go-clean-template/internal/controller/restapi/middleware"
	v1 "github.com/evrone/go-clean-template/internal/controller/restapi/v1"
	cartmodule "github.com/evrone/go-clean-template/internal/module/cart"
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	"github.com/evrone/go-clean-template/internal/module/catalog/catalogbase"
	importermodule "github.com/evrone/go-clean-template/internal/module/importer"
	leadmodule "github.com/evrone/go-clean-template/internal/module/lead"
	mediamodule "github.com/evrone/go-clean-template/internal/module/media"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

const selfContainedDocsHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Go-Grip Backend API Documentation</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" charset="UTF-8"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/openapi.yaml",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIBundle.SwaggerUIStandalonePreset
        ]
      });
    };
  </script>
</body>
</html>`

// NewRouter registers strict OpenAPI 3.0 routes and global middlewares.
func NewRouter(app *fiber.App, cfg *config.Config, u usermodule.UserUseCase, catalog catalogmodule.CatalogUseCase, catalogBase catalogbase.UseCase, auth usermodule.AuthUseCase, checkout ordermodule.CheckoutUseCase, orders ordermodule.OrdersUseCase, profile usermodule.ProfileUseCase, admin usermodule.AdminUseCase, wishlist wishlistmodule.WishlistUseCase, notify notificationmodule.NotificationCenterUseCase, media mediamodule.MediaUseCase, homepage contentmodule.HomepageUseCase, cart cartmodule.CartUseCase, lead leadmodule.LeadUseCase, content contentmodule.ContentUseCase, importer importermodule.ImporterUseCase, jwtManager *jwt.Manager, l logger.Interface) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))
	app.Use(middleware.RejectBlockedMutations())
	if jwtManager != nil {
		app.Use(middleware.Auth(jwtManager))
	}
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

	// Serve raw OpenAPI spec
	app.Get("/openapi.yaml", func(c *fiber.Ctx) error {
		content, err := os.ReadFile("docs/api/openapi.yaml")
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Spec file not found")
		}
		c.Set(fiber.HeaderContentType, "text/yaml; charset=utf-8")
		return c.Send(content)
	})

	// Swagger & Self-Contained UI Docs
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
		app.Get("/docs", func(c *fiber.Ctx) error {
			c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
			return c.SendString(selfContainedDocsHTML)
		})
	}

	// K8s probe
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// Strict OpenAPI Composition Server
	server := v1.NewServer(
		cfg, u, catalog, catalogBase, auth, checkout, orders, profile, admin, wishlist, notify, media, homepage, cart, lead, content, importer, jwtManager, l,
	)

	strictHandler := openapi.NewStrictHandler(server, nil)
	openapi.RegisterHandlersWithOptions(app, strictHandler, openapi.FiberServerOptions{
		BaseURL: "/v1",
	})
}
