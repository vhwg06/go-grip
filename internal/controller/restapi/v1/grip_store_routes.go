package v1

import (
	"net/http"
	"time"

	"github.com/evrone/go-clean-template/internal/controller/restapi/middleware"
	"github.com/gofiber/fiber/v2"
)

func (r *V1) registerGripStoreRoutes(apiV1Group fiber.Router, protected fiber.Router) {
	authGroup := apiV1Group.Group("/auth")
	authGroup.Get("/oauth/linuxdo", r.gripBeginLinuxDO)
	authGroup.Get("/oauth/github", r.gripBeginGitHub)
	authGroup.Get("/callback/linuxdo", r.gripCompleteLinuxDO)
	authGroup.Get("/callback/github", r.gripCompleteGitHub)
	authGroup.Post("/refresh", r.gripRefresh)
	authGroup.Post("/logout", middleware.Auth(r.jwtManager), r.gripLogout)
	authGroup.Get("/me", middleware.Auth(r.jwtManager), r.gripMe)

	catalogGroup := apiV1Group.Group("/catalog")
	catalogGroup.Get("/products", r.gripListProducts)
	catalogGroup.Get("/products/:id", r.gripGetProduct)
	catalogGroup.Get("/products/:id/buy-meta", r.gripGetBuyMeta)
	catalogGroup.Get("/search", r.gripSearchProducts)
	catalogGroup.Get("/categories", r.gripListCategories)
	catalogGroup.Get("/settings", r.gripListSettings)
	catalogGroup.Get("/announcement", r.gripGetAnnouncement)

	checkoutGroup := apiV1Group.Group("/checkout", middleware.RateLimitByIP(30, time.Minute))
	checkoutGroup.Get("/preview", r.gripCheckoutPreview)
	checkoutGroup.Post("/orders", r.gripCreateOrder)
	checkoutGroup.Post("/payment-orders", r.gripCreatePaymentOrder)
	checkoutGroup.Get("/orders/:id/payment-params", r.gripPaymentParams)
	checkoutGroup.Get("/orders/:id/status", r.gripOrderStatus)
	checkoutGroup.Post("/orders/:id/cancel", middleware.Auth(r.jwtManager), r.gripCancelOrder)
	checkoutGroup.Post("/notify", r.gripPaymentNotify)
	checkoutGroup.Get("/callback/:id", r.gripPaymentCallback)

	ordersGroup := apiV1Group.Group("/orders", middleware.Auth(r.jwtManager))
	ordersGroup.Get("/", r.gripListOrders)
	ordersGroup.Get("/:id", r.gripGetOrder)
	ordersGroup.Post("/:id/refund-request", r.gripRequestRefund)

	profileGroup := protected.Group("/profile")
	profileGroup.Get("/", r.gripProfileGet)
	profileGroup.Patch("/", r.gripProfileUpdate)
	profileGroup.Post("/check-in", r.gripProfileCheckin)

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
