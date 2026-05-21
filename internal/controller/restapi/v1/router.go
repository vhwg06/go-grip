package v1

import (
	"github.com/evrone/go-clean-template/internal/controller/restapi/middleware"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// NewRoutes -.
func NewRoutes(apiV1Group fiber.Router, t usecase.Translation, u usecase.User, tk usecase.Task, catalog usecase.Catalog, media usecase.Media, homepage usecase.Homepage, cart usecase.Cart, lead usecase.Lead, content usecase.Content, importer usecase.Importer, jwtManager *jwt.Manager, l logger.Interface) {
	r := &V1{t: t, u: u, tk: tk, catalog: catalog, media: media, homepage: homepage, cart: cart, lead: lead, content: content, importer: importer, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	// Public routes
	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.Post("/register", r.register)
		authGroup.Post("/login", r.login)
	}

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

	r.registerEcommerceRoutes(apiV1Group, protected)
}
