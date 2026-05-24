package v1

import (
	"context"
	"net/http"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type adminExtendedUseCase interface {
	ListProducts(ctx context.Context, actor entity.Actor, page entity.Pagination) ([]entity.Product, int, error)
	UpsertProduct(ctx context.Context, actor entity.Actor, product entity.Product) (entity.Product, error)
	DeleteProduct(ctx context.Context, actor entity.Actor, productID string) error
	ListCategories(ctx context.Context, actor entity.Actor) ([]entity.Category, error)
	UpsertCategory(ctx context.Context, actor entity.Actor, category entity.Category) (entity.Category, error)
	DeleteCategory(ctx context.Context, actor entity.Actor, categoryID string) error
	ListOrders(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Order, int, error)
	GetOrder(ctx context.Context, actor entity.Actor, orderID string) (entity.Order, error)
	ImportCards(ctx context.Context, actor entity.Actor, productID string, keys []string) (int, error)
	SendBroadcast(ctx context.Context, actor entity.Actor, title, body string) error
	SendTargeted(ctx context.Context, actor entity.Actor, userID, title, body string) error
}

type gripImportCardsRequest struct {
	ProductID string   `json:"productId"`
	Keys      []string `json:"keys"`
}

// @Summary     List admin products
// @Description Lists products for admin management
// @ID          grip_admin_list_products
// @Tags        admin
// @Produce     json
// @Success     200 {object} gripListResponse
// @Failure     403 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /admin/products [get]
func (r *V1) gripAdminListProducts(ctx *fiber.Ctx) error {
	if r.adminUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_usecase_not_configured"})
	}

	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_products_not_available"})
	}

	page := gripPage(ctx)
	items, total, err := ext.ListProducts(ctx.UserContext(), r.gripActor(ctx), page)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	normalized := page.Normalize()
	return ctx.JSON(gripListResponse{Data: items, Meta: entity.Page{Limit: normalized.Limit, Offset: normalized.Offset, Total: total}})
}

// @Summary     Create admin product
// @Description Creates or inserts a product from admin panel
// @ID          grip_admin_create_product
// @Tags        admin
// @Accept      json
// @Produce     json
// @Success     201 {object} envelope
// @Failure     400 {object} envelope
// @Failure     403 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /admin/products [post]
func (r *V1) gripAdminCreateProduct(ctx *fiber.Ctx) error {
	if r.adminUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_usecase_not_configured"})
	}
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_products_not_available"})
	}

	var product entity.Product
	if err := ctx.BodyParser(&product); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	created, err := ext.UpsertProduct(ctx.UserContext(), r.gripActor(ctx), product)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.Status(http.StatusCreated).JSON(apiSuccessEnvelope(created))
}

// @Summary     Update admin product
// @Description Updates existing product fields
// @ID          grip_admin_update_product
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param       id path string true "Product ID"
// @Success     200 {object} envelope
// @Failure     400 {object} envelope
// @Failure     403 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /admin/products/{id} [patch]
func (r *V1) gripAdminUpdateProduct(ctx *fiber.Ctx) error {
	if r.adminUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_usecase_not_configured"})
	}
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_products_not_available"})
	}

	var product entity.Product
	if err := ctx.BodyParser(&product); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	product.ID = ctx.Params("id")

	updated, err := ext.UpsertProduct(ctx.UserContext(), r.gripActor(ctx), product)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.JSON(apiSuccessEnvelope(updated))
}

// @Summary     Delete admin product
// @Description Deletes a product by ID
// @ID          grip_admin_delete_product
// @Tags        admin
// @Produce     json
// @Param       id path string true "Product ID"
// @Success     204 {string} string ""
// @Failure     403 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /admin/products/{id} [delete]
func (r *V1) gripAdminDeleteProduct(ctx *fiber.Ctx) error {
	if r.adminUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_usecase_not_configured"})
	}
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_products_not_available"})
	}
	if err := ext.DeleteProduct(ctx.UserContext(), r.gripActor(ctx), ctx.Params("id")); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) gripAdminListCategories(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_categories_not_available"})
	}
	items, err := ext.ListCategories(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.JSON(apiSuccessEnvelope(items))
}

func (r *V1) gripAdminCreateCategory(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_categories_not_available"})
	}

	var category entity.Category
	if err := ctx.BodyParser(&category); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	created, err := ext.UpsertCategory(ctx.UserContext(), r.gripActor(ctx), category)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.Status(http.StatusCreated).JSON(apiSuccessEnvelope(created))
}

func (r *V1) gripAdminUpdateCategory(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_categories_not_available"})
	}

	var category entity.Category
	if err := ctx.BodyParser(&category); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	category.ID = ctx.Params("id")
	updated, err := ext.UpsertCategory(ctx.UserContext(), r.gripActor(ctx), category)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.JSON(apiSuccessEnvelope(updated))
}

func (r *V1) gripAdminDeleteCategory(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_categories_not_available"})
	}
	if err := ext.DeleteCategory(ctx.UserContext(), r.gripActor(ctx), ctx.Params("id")); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) gripAdminListOrders(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_orders_not_available"})
	}

	pageNum := ctx.QueryInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize := ctx.QueryInt("pageSize", 50)
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}
	query := strings.TrimSpace(ctx.Query("q"))
	status := strings.TrimSpace(ctx.Query("status", "all"))

	page := entity.Pagination{
		Limit:  pageSize,
		Offset: (pageNum - 1) * pageSize,
	}

	items, total, err := ext.ListOrders(ctx.UserContext(), r.gripActor(ctx), page, query, status)
	if err != nil {
		httpStatus, payload := mapDomainError(err)
		return ctx.Status(httpStatus).JSON(payload)
	}

	orders := make([]fiber.Map, 0, len(items))
	for _, order := range items {
		orders = append(orders, fiber.Map{
			"orderId":     order.ID,
			"userId":      order.UserID,
			"username":    order.Username,
			"email":       order.Email,
			"productName": order.ProductName,
			"amount":      order.Amount,
			"status":      string(order.Status),
			"cardKey":     order.CardKey,
			"tradeNo":     order.TradeNo,
			"createdAt":   order.CreatedAt,
		})
	}

	return ctx.JSON(fiber.Map{
		"orders":   orders,
		"total":    total,
		"page":     pageNum,
		"pageSize": pageSize,
		"query":    query,
		"status":   status,
	})
}

func (r *V1) gripAdminGetOrder(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_orders_not_available"})
	}

	order, err := ext.GetOrder(ctx.UserContext(), r.gripActor(ctx), ctx.Params("id"))
	if err != nil {
		httpStatus, payload := mapDomainError(err)
		return ctx.Status(httpStatus).JSON(payload)
	}

	status := strings.ToUpper(string(order.Status))
	return ctx.JSON(fiber.Map{
		"id":          order.ID,
		"orderNumber": order.ID,
		"status":      status,
		"createdAt":   order.CreatedAt,
		"items": []fiber.Map{
			{
				"productName": order.ProductName,
				"sku":         "",
				"price":       order.Amount,
				"quantity":    max(order.Quantity, 1),
			},
		},
		"totalAmount":     order.Amount,
		"customerName":    order.Username,
		"customerPhone":   "",
		"customerEmail":   order.Email,
		"shippingAddress": "",
		"paymentMethod":   "",
		"timeline": []fiber.Map{
			{
				"status":    status,
				"timestamp": order.UpdatedAt,
				"note":      "",
			},
		},
	})
}

// @Summary     Import product cards
// @Description Bulk imports card keys for a product
// @ID          grip_admin_import_cards
// @Tags        admin
// @Accept      json
// @Produce     json
// @Success     200 {object} envelope
// @Failure     400 {object} envelope
// @Failure     403 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /admin/cards/import [post]
func (r *V1) gripAdminCardsImport(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_cards_not_available"})
	}

	var body gripImportCardsRequest
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	count, err := ext.ImportCards(ctx.UserContext(), r.gripActor(ctx), body.ProductID, body.Keys)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.JSON(apiSuccessEnvelope(fiber.Map{"imported": count}))
}

// @Summary     List users for admin
// @Description Lists users with pagination
// @ID          grip_admin_list_users
// @Tags        admin
// @Produce     json
// @Success     200 {object} gripListResponse
// @Failure     403 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /admin/users [get]
func (r *V1) gripAdminUsersList(ctx *fiber.Ctx) error {
	page := gripPage(ctx)
	users, total, err := r.adminUC.ListUsers(ctx.UserContext(), r.gripActor(ctx), page)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	normalized := page.Normalize()
	return ctx.JSON(gripListResponse{Data: users, Meta: entity.Page{Limit: normalized.Limit, Offset: normalized.Offset, Total: total}})
}

// @Summary     Update admin user state
// @Description Updates user status or points from admin panel
// @ID          grip_admin_update_user
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param       id path string true "User ID"
// @Success     204 {string} string ""
// @Failure     400 {object} envelope
// @Failure     403 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /admin/users/{id} [patch]
func (r *V1) gripAdminUsersUpdate(ctx *fiber.Ctx) error {
	var body struct {
		Status *entity.UserStatus `json:"status"`
		Points *int               `json:"points"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	actor := r.gripActor(ctx)
	if body.Status != nil {
		if err := r.adminUC.UpdateUserStatus(ctx.UserContext(), actor, ctx.Params("id"), *body.Status); err != nil {
			status, payload := mapDomainError(err)
			return ctx.Status(status).JSON(payload)
		}
	}
	if body.Points != nil {
		if err := r.adminUC.UpdateUserPoints(ctx.UserContext(), actor, ctx.Params("id"), *body.Points); err != nil {
			status, payload := mapDomainError(err)
			return ctx.Status(status).JSON(payload)
		}
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Repair aggregates
// @Description Recalculates product stock aggregates
// @ID          grip_admin_repair_aggregates
// @Tags        admin
// @Produce     json
// @Success     200 {object} envelope
// @Failure     403 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /admin/data/repair-aggregates [post]
func (r *V1) gripAdminRepairAggregates(ctx *fiber.Ctx) error {
	if err := r.adminUC.RepairAggregates(ctx.UserContext(), r.gripActor(ctx)); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.JSON(apiSuccessEnvelope(fiber.Map{"status": "ok"}))
}

func (r *V1) gripAdminBroadcast(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_messages_not_available"})
	}

	var body struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if err := ext.SendBroadcast(ctx.UserContext(), r.gripActor(ctx), body.Title, body.Body); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) gripAdminTargeted(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_messages_not_available"})
	}

	var body struct {
		UserID string `json:"userId"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if err := ext.SendTargeted(ctx.UserContext(), r.gripActor(ctx), body.UserID, body.Title, body.Body); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) gripAdminNoop(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusNotImplemented).JSON(envelope{Error: "not_implemented"})
}
