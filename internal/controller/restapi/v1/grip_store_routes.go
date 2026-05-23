package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (r *V1) registerGripStoreRoutes(apiV1Group fiber.Router, protected fiber.Router) {
	authGroup := apiV1Group.Group("/auth")
	authGroup.Get("/oauth/linuxdo", r.notImplemented)
	authGroup.Get("/oauth/github", r.notImplemented)
	authGroup.Get("/callback/linuxdo", r.notImplemented)
	authGroup.Get("/callback/github", r.notImplemented)
	authGroup.Post("/refresh", r.notImplemented)
	authGroup.Post("/logout", r.notImplemented)
	authGroup.Get("/me", r.notImplemented)

	catalogGroup := apiV1Group.Group("/catalog")
	catalogGroup.Get("/products/:id/buy-meta", r.notImplemented)
	catalogGroup.Get("/search", r.notImplemented)
	catalogGroup.Get("/settings", r.notImplemented)
	catalogGroup.Get("/announcement", r.notImplemented)

	checkoutGroup := apiV1Group.Group("/checkout")
	checkoutGroup.Get("/preview", r.notImplemented)
	checkoutGroup.Post("/orders", r.notImplemented)
	checkoutGroup.Post("/payment-orders", r.notImplemented)
	checkoutGroup.Get("/orders/:id/payment-params", r.notImplemented)
	checkoutGroup.Get("/orders/:id/status", r.notImplemented)
	checkoutGroup.Post("/orders/:id/cancel", r.notImplemented)
	checkoutGroup.Post("/notify", r.notImplemented)
	checkoutGroup.Get("/callback/:id", r.notImplemented)

	ordersGroup := protected.Group("/orders")
	ordersGroup.Get("/", r.notImplemented)
	ordersGroup.Get("/:id", r.notImplemented)
	ordersGroup.Post("/:id/refund-request", r.notImplemented)

	profileGroup := protected.Group("/profile")
	profileGroup.Get("/", r.notImplemented)
	profileGroup.Patch("/", r.notImplemented)
	profileGroup.Post("/check-in", r.notImplemented)

	wishlistGroup := apiV1Group.Group("/wishlist")
	wishlistGroup.Get("/", r.notImplemented)
	wishlistGroup.Post("/", r.notImplemented)
	wishlistGroup.Patch("/:id", r.notImplemented)
	wishlistGroup.Delete("/:id", r.notImplemented)
	wishlistGroup.Post("/:id/vote", r.notImplemented)

	reviewsGroup := protected.Group("/reviews")
	reviewsGroup.Post("/", r.notImplemented)

	notificationGroup := protected.Group("/notifications")
	notificationGroup.Get("/", r.notImplemented)
	notificationGroup.Get("/unread-count", r.notImplemented)
	notificationGroup.Post("/:id/read", r.notImplemented)
	notificationGroup.Post("/read-all", r.notImplemented)
	notificationGroup.Delete("/", r.notImplemented)

	adminGroup := protected.Group("/admin")
	adminGroup.Get("/products", r.notImplemented)
	adminGroup.Post("/products", r.notImplemented)
	adminGroup.Patch("/products/:id", r.notImplemented)
	adminGroup.Delete("/products/:id", r.notImplemented)
	adminGroup.Get("/categories", r.notImplemented)
	adminGroup.Post("/categories", r.notImplemented)
	adminGroup.Patch("/categories/:id", r.notImplemented)
	adminGroup.Delete("/categories/:id", r.notImplemented)
	adminGroup.Get("/cards", r.notImplemented)
	adminGroup.Post("/cards", r.notImplemented)
	adminGroup.Delete("/cards/:id", r.notImplemented)
	adminGroup.Post("/cards/import", r.notImplemented)
	adminGroup.Post("/cards/replenish", r.notImplemented)
	adminGroup.Get("/orders", r.notImplemented)
	adminGroup.Patch("/orders/:id", r.notImplemented)
	adminGroup.Delete("/orders/:id", r.notImplemented)
	adminGroup.Get("/refunds", r.notImplemented)
	adminGroup.Post("/refunds/:id/approve", r.notImplemented)
	adminGroup.Post("/refunds/:id/reject", r.notImplemented)
	adminGroup.Get("/users", r.notImplemented)
	adminGroup.Patch("/users/:id", r.notImplemented)
	adminGroup.Get("/settings", r.notImplemented)
	adminGroup.Put("/settings/:key", r.notImplemented)
	adminGroup.Delete("/settings/:key", r.notImplemented)
	adminGroup.Post("/messages/broadcast", r.notImplemented)
	adminGroup.Post("/messages/targeted", r.notImplemented)
	adminGroup.Post("/notifications/test", r.notImplemented)
	adminGroup.Post("/data/import", r.notImplemented)
	adminGroup.Post("/data/repair-aggregates", r.notImplemented)
}

func (r *V1) notImplemented(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusNotImplemented).JSON(fiber.Map{
		"error": "not_implemented",
	})
}
