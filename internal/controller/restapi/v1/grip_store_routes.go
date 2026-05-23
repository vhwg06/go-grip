package v1

import (
	"net/http"
	"time"

	"github.com/evrone/go-clean-template/internal/controller/restapi/middleware"
	"github.com/gofiber/fiber/v2"
)

func (r *V1) registerGripStoreRoutes(apiV1Group fiber.Router, protected fiber.Router) {
	authGroup := apiV1Group.Group("/auth")
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
	wishlistGroup.Get("/", r.gripWishlistList)
	wishlistGroup.Post("/", middleware.Auth(r.jwtManager), r.gripWishlistCreate)
	wishlistGroup.Patch("/:id", middleware.Auth(r.jwtManager), r.gripWishlistUpdate)
	wishlistGroup.Delete("/:id", middleware.Auth(r.jwtManager), r.gripWishlistDelete)
	wishlistGroup.Post("/:id/vote", middleware.Auth(r.jwtManager), r.gripWishlistVote)

	reviewsGroup := protected.Group("/reviews")
	reviewsGroup.Post("/", r.gripReviewCreate)

	notificationGroup := protected.Group("/notifications")
	notificationGroup.Get("/", r.gripNotificationsList)
	notificationGroup.Get("/unread-count", r.gripNotificationsUnread)
	notificationGroup.Post("/:id/read", r.gripNotificationsMarkRead)
	notificationGroup.Post("/read-all", r.gripNotificationsReadAll)
	notificationGroup.Delete("/", r.gripNotificationsClear)

	adminGroup := protected.Group("/admin", middleware.RequireAdminUsernames(r.adminUsers))
	adminGroup.Get("/products", r.gripAdminListProducts)
	adminGroup.Post("/products", r.gripAdminCreateProduct)
	adminGroup.Patch("/products/:id", r.gripAdminUpdateProduct)
	adminGroup.Delete("/products/:id", r.gripAdminDeleteProduct)
	adminGroup.Get("/categories", r.gripAdminListCategories)
	adminGroup.Post("/categories", r.gripAdminCreateCategory)
	adminGroup.Patch("/categories/:id", r.gripAdminUpdateCategory)
	adminGroup.Delete("/categories/:id", r.gripAdminDeleteCategory)
	adminGroup.Get("/cards", r.gripAdminNoop)
	adminGroup.Post("/cards", r.gripAdminNoop)
	adminGroup.Delete("/cards/:id", r.gripAdminNoop)
	adminGroup.Post("/cards/import", r.gripAdminCardsImport)
	adminGroup.Post("/cards/replenish", r.gripAdminNoop)
	adminGroup.Get("/orders", r.gripAdminNoop)
	adminGroup.Patch("/orders/:id", r.gripAdminNoop)
	adminGroup.Delete("/orders/:id", r.gripAdminNoop)
	adminGroup.Get("/refunds", r.gripAdminNoop)
	adminGroup.Post("/refunds/:id/approve", r.gripAdminNoop)
	adminGroup.Post("/refunds/:id/reject", r.gripAdminNoop)
	adminGroup.Get("/users", r.gripAdminUsersList)
	adminGroup.Patch("/users/:id", r.gripAdminUsersUpdate)
	adminGroup.Get("/settings", r.gripAdminNoop)
	adminGroup.Put("/settings/:key", r.gripAdminNoop)
	adminGroup.Delete("/settings/:key", r.gripAdminNoop)
	adminGroup.Post("/messages/broadcast", r.gripAdminBroadcast)
	adminGroup.Post("/messages/targeted", r.gripAdminTargeted)
	adminGroup.Post("/notifications/test", r.gripAdminNoop)
	adminGroup.Post("/data/import", r.gripAdminNoop)
	adminGroup.Post("/data/repair-aggregates", r.gripAdminRepairAggregates)
}

func (r *V1) notImplemented(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusNotImplemented).JSON(fiber.Map{
		"error": "not_implemented",
	})
}
