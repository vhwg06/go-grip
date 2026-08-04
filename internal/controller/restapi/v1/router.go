package v1

import (
	"time"

	"github.com/evrone/go-clean-template/internal/controller/restapi/middleware"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/internal/usecase/catalogbase"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// NewRoutes -.
func NewRoutes(apiV1Group fiber.Router, t usecase.Translation, u usecase.User, tk usecase.Task, catalog usecase.Catalog, auth usecase.Auth, checkout usecase.Checkout, orders usecase.Orders, profile usecase.Profile, admin usecase.Admin, wishlist usecase.Wishlist, notify usecase.NotificationCenter, media usecase.Media, homepage usecase.Homepage, cart usecase.Cart, lead usecase.Lead, content usecase.Content, importer usecase.Importer, jwtManager *jwt.Manager, adminUsers string, l logger.Interface) {
	r := &V1{
		t:          t,
		u:          u,
		tk:         tk,
		catalog:    catalog,
		authUC:     auth,
		checkout:   checkout,
		orders:     orders,
		profileUC:  profile,
		adminUC:    admin,
		wishlistUC: wishlist,
		notifyUC:   notify,
		media:      media,
		homepage:   homepage,
		cart:       cart,
		lead:       lead,
		content:    content,
		importer:   importer,
		jwtManager: jwtManager,
		adminUsers: adminUsers,
		l:          l,
		v:          validator.New(validator.WithRequiredStructEnabled()),
	}
	if provider, ok := catalog.(interface{ CatalogBaseService() *catalogbase.Service }); ok {
		r.catalogBase = provider.CatalogBaseService()
	}

	// Public routes
	authGroup := apiV1Group.Group("/auth", middleware.RateLimitByIP(5, time.Second))
	{
		authGroup.Post("/register", r.register)
		authGroup.Post("/login", r.login)
	}

	r.registerGripStoreRoutes(apiV1Group)
	r.registerEcommerceRoutes(apiV1Group)

	// Protected routes
	protected := apiV1Group.Group("", middleware.Auth(jwtManager))

	userGroup := protected.Group("/user")
	{
		userGroup.Get("/profile", r.profile)
	}

	usersGroup := protected.Group("/users")
	{
		usersGroup.Get("/", r.listUsers)
		usersGroup.Post("/", r.createAdminUser)
		usersGroup.Get("/:id", r.getUser)
		usersGroup.Patch("/:id", r.updateUserProfile)
		usersGroup.Post("/:id/lock", r.lockUser)
		usersGroup.Post("/:id/unlock", r.unlockUser)
	}

	taskGroup := protected.Group("/tasks")
	{
		taskGroup.Post("/", r.createTask)
		taskGroup.Get("/", r.listTasks)
		taskGroup.Get("/:id", r.getTask)
		taskGroup.Put("/:id", r.updateTask)
		taskGroup.Patch("/:id/status", r.transitionTask)
		taskGroup.Delete("/:id", r.deleteTask)
	}

	translationGroup := protected.Group("/translation")
	{
		translationGroup.Get("/history", r.history)
		translationGroup.Post("/do-translate", r.doTranslate)
	}
}
