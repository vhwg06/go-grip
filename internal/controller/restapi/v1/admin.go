package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/gofiber/fiber/v2"
)

type adminExtendedUseCase interface {
	ListProducts(ctx context.Context, actor entity.Actor, page entity.Pagination) ([]entity.Product, int, error)
	GetProduct(ctx context.Context, actor entity.Actor, productID string) (entity.Product, error)
	UpsertProduct(ctx context.Context, actor entity.Actor, product entity.Product) (entity.Product, error)
	DeleteProduct(ctx context.Context, actor entity.Actor, productID string) error
	ListCategories(ctx context.Context, actor entity.Actor) ([]entity.Category, error)
	UpsertCategory(ctx context.Context, actor entity.Actor, category entity.Category) (entity.Category, error)
	DeleteCategory(ctx context.Context, actor entity.Actor, categoryID string) error
	ListOrders(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Order, int, error)
	GetOrder(ctx context.Context, actor entity.Actor, orderID string) (entity.Order, error)
	UpdateOrderStatus(ctx context.Context, actor entity.Actor, orderID string, status entity.OrderStatus) error
	DeleteOrder(ctx context.Context, actor entity.Actor, orderID string) error
	ListRefunds(ctx context.Context, actor entity.Actor, status string) ([]entity.RefundRequest, error)
	GetRefund(ctx context.Context, actor entity.Actor, refundID int64) (entity.RefundRequest, error)
	GetOrderRefundStatus(ctx context.Context, actor entity.Actor, orderID string) (entity.RefundRequest, error)
	DecideRefund(ctx context.Context, actor entity.Actor, refundID int64, approve bool, note string) (entity.RefundRequest, error)
	ListReviews(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error)
	UpdateReviewStatus(ctx context.Context, actor entity.Actor, reviewID int64, status entity.ReviewStatus) (entity.Review, error)
	BulkPublishReviews(ctx context.Context, actor entity.Actor, reviewIDs []int64) (int, error)
	DeleteReview(ctx context.Context, actor entity.Actor, reviewID int64) error
	ListSettings(ctx context.Context, actor entity.Actor) ([]entity.Setting, error)
	SetSetting(ctx context.Context, actor entity.Actor, key, value string) error
	DeleteSetting(ctx context.Context, actor entity.Actor, key string) error
	SendBroadcast(ctx context.Context, actor entity.Actor, title, body string) error
	SendTargeted(ctx context.Context, actor entity.Actor, userID, title, body string) error
	ListAdminMessages(ctx context.Context, actor entity.Actor) ([]entity.AdminMessage, error)
	ListCards(ctx context.Context, actor entity.Actor) ([]entity.Card, error)
}

func parseBoolForm(value string) bool {
	normalized := strings.TrimSpace(strings.ToLower(value))
	return normalized == "1" || normalized == "true" || normalized == "on" || normalized == "yes"
}

func patchProductFromForm(product *entity.Product, ctx *fiber.Ctx) {
	if value := strings.TrimSpace(ctx.FormValue("id")); value != "" {
		product.ID = value
	}
	if value := strings.TrimSpace(ctx.FormValue("slug")); value != "" && product.ID == "" {
		product.ID = value
	}
	if value := strings.TrimSpace(ctx.FormValue("name")); value != "" {
		product.Title = value
	}
	if value := strings.TrimSpace(ctx.FormValue("description")); value != "" {
		product.Description = value
	}
	if value := strings.TrimSpace(ctx.FormValue("category")); value != "" {
		product.CategoryID = value
	}
	if value := strings.TrimSpace(ctx.FormValue("image")); value != "" {
		product.ImageURL = value
	}
	if value := strings.TrimSpace(ctx.FormValue("images")); value != "" {
		lines := strings.Split(value, "\n")
		images := make([]string, 0, len(lines))
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				images = append(images, line)
			}
		}
		product.Images = images
	}
	if value := strings.TrimSpace(ctx.FormValue("sku")); value != "" {
		product.SKU = value
	}
	if value := strings.TrimSpace(ctx.FormValue("brand")); value != "" {
		product.Brand = value
	}
	if value := strings.TrimSpace(ctx.FormValue("brandId")); value != "" && product.Brand == "" {
		product.Brand = value
	}
	if value := strings.TrimSpace(ctx.FormValue("price")); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			product.Price = int64(parsed)
		}
	}
	if value := strings.TrimSpace(ctx.FormValue("compareAtPrice")); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			price := int64(parsed)
			product.ComparePrice = &price
		}
	}
	if value := strings.TrimSpace(ctx.FormValue("purchaseLimit")); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			product.PurchaseLimit = parsed
		}
	}
	if value := strings.TrimSpace(ctx.FormValue("purchaseWarning")); value != "" {
		product.PurchaseWarning = value
	}
	if value := strings.TrimSpace(ctx.FormValue("visibilityLevel")); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			product.VisibilityLevel = parsed
		}
	}
	if value := strings.TrimSpace(ctx.FormValue("isHot")); value != "" {
		product.IsHot = parseBoolForm(value)
	}
	if value := strings.TrimSpace(ctx.FormValue("isActive")); value != "" {
		product.IsActive = parseBoolForm(value)
	}
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
	contentType := strings.ToLower(ctx.Get("Content-Type"))
	isMultipart := strings.HasPrefix(contentType, "multipart/form-data")
	if err := ctx.BodyParser(&product); err != nil && !isMultipart {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	patchProductFromForm(&product, ctx)

	if specsStr := ctx.FormValue("specs"); specsStr != "" {
		var specs []entity.ProductSpecItem
		if err := json.Unmarshal([]byte(specsStr), &specs); err == nil {
			product.Specs = specs
		}
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
func mergeProductPatch(existing *entity.Product, bodyMap map[string]any, ctx *fiber.Ctx, isMultipart bool) {
	if isMultipart {
		if value := strings.TrimSpace(ctx.FormValue("name")); value != "" {
			existing.Title = value
		}
		if value := strings.TrimSpace(ctx.FormValue("title")); value != "" {
			existing.Title = value
		}
		if value := strings.TrimSpace(ctx.FormValue("description")); value != "" {
			existing.Description = value
		}
		if value := strings.TrimSpace(ctx.FormValue("category")); value != "" {
			existing.CategoryID = value
		}
		if value := strings.TrimSpace(ctx.FormValue("categoryId")); value != "" {
			existing.CategoryID = value
		}
		if value := strings.TrimSpace(ctx.FormValue("category_id")); value != "" {
			existing.CategoryID = value
		}
		if value := strings.TrimSpace(ctx.FormValue("image")); value != "" {
			existing.ImageURL = value
		}
		if value := strings.TrimSpace(ctx.FormValue("image_url")); value != "" {
			existing.ImageURL = value
		}
		if value := strings.TrimSpace(ctx.FormValue("sku")); value != "" {
			existing.SKU = value
		}
		if value := strings.TrimSpace(ctx.FormValue("brand")); value != "" {
			existing.Brand = value
		}
		if value := strings.TrimSpace(ctx.FormValue("price")); value != "" {
			if parsed, err := strconv.ParseFloat(value, 64); err == nil {
				existing.Price = int64(parsed)
			}
		}
		if value := strings.TrimSpace(ctx.FormValue("compareAtPrice")); value != "" {
			if parsed, err := strconv.ParseFloat(value, 64); err == nil {
				val := int64(parsed)
				existing.ComparePrice = &val
			}
		}
		if value := strings.TrimSpace(ctx.FormValue("compare_price")); value != "" {
			if parsed, err := strconv.ParseFloat(value, 64); err == nil {
				val := int64(parsed)
				existing.ComparePrice = &val
			}
		}
		if value := strings.TrimSpace(ctx.FormValue("isHot")); value != "" {
			existing.IsHot = parseBoolForm(value)
		}
		if value := strings.TrimSpace(ctx.FormValue("is_hot")); value != "" {
			existing.IsHot = parseBoolForm(value)
		}
		if value := strings.TrimSpace(ctx.FormValue("isActive")); value != "" {
			existing.IsActive = parseBoolForm(value)
		}
		if value := strings.TrimSpace(ctx.FormValue("is_active")); value != "" {
			existing.IsActive = parseBoolForm(value)
		}
		if value := strings.TrimSpace(ctx.FormValue("purchaseLimit")); value != "" {
			if parsed, err := strconv.Atoi(value); err == nil {
				existing.PurchaseLimit = parsed
			}
		}
		if value := strings.TrimSpace(ctx.FormValue("purchase_limit")); value != "" {
			if parsed, err := strconv.Atoi(value); err == nil {
				existing.PurchaseLimit = parsed
			}
		}
		if value := strings.TrimSpace(ctx.FormValue("purchaseWarning")); value != "" {
			existing.PurchaseWarning = value
		}
		if value := strings.TrimSpace(ctx.FormValue("purchase_warning")); value != "" {
			existing.PurchaseWarning = value
		}
		if value := strings.TrimSpace(ctx.FormValue("visibilityLevel")); value != "" {
			if parsed, err := strconv.Atoi(value); err == nil {
				existing.VisibilityLevel = parsed
			}
		}
		if value := strings.TrimSpace(ctx.FormValue("visibility_level")); value != "" {
			if parsed, err := strconv.Atoi(value); err == nil {
				existing.VisibilityLevel = parsed
			}
		}
		return
	}

	if bodyMap == nil {
		return
	}

	if _, ok := bodyMap["name"]; ok {
		if val, ok := bodyMap["name"].(string); ok {
			existing.Title = val
		}
	}
	if _, ok := bodyMap["title"]; ok {
		if val, ok := bodyMap["title"].(string); ok {
			existing.Title = val
		}
	}
	if _, ok := bodyMap["description"]; ok {
		if val, ok := bodyMap["description"].(string); ok {
			existing.Description = val
		}
	}
	if _, ok := bodyMap["price"]; ok {
		if val, ok := bodyMap["price"].(float64); ok {
			existing.Price = int64(val)
		}
	}
	if _, ok := bodyMap["compareAtPrice"]; ok {
		if val, ok := bodyMap["compareAtPrice"].(float64); ok {
			comp := int64(val)
			existing.ComparePrice = &comp
		} else if bodyMap["compareAtPrice"] == nil {
			existing.ComparePrice = nil
		}
	}
	if _, ok := bodyMap["compare_price"]; ok {
		if val, ok := bodyMap["compare_price"].(float64); ok {
			comp := int64(val)
			existing.ComparePrice = &comp
		} else if bodyMap["compare_price"] == nil {
			existing.ComparePrice = nil
		}
	}
	if _, ok := bodyMap["categoryId"]; ok {
		if val, ok := bodyMap["categoryId"].(string); ok {
			existing.CategoryID = val
		} else if val, ok := bodyMap["categoryId"].(float64); ok {
			existing.CategoryID = strconv.Itoa(int(val))
		}
	}
	if _, ok := bodyMap["category_id"]; ok {
		if val, ok := bodyMap["category_id"].(string); ok {
			existing.CategoryID = val
		} else if val, ok := bodyMap["category_id"].(float64); ok {
			existing.CategoryID = strconv.Itoa(int(val))
		}
	}
	if _, ok := bodyMap["image"]; ok {
		if val, ok := bodyMap["image"].(string); ok {
			existing.ImageURL = val
		}
	}
	if _, ok := bodyMap["image_url"]; ok {
		if val, ok := bodyMap["image_url"].(string); ok {
			existing.ImageURL = val
		}
	}
	if _, ok := bodyMap["images"]; ok {
		if arr, ok := bodyMap["images"].([]any); ok {
			imgs := make([]string, 0, len(arr))
			for _, item := range arr {
				if s, ok := item.(string); ok {
					imgs = append(imgs, s)
				}
			}
			existing.Images = imgs
		}
	}
	if _, ok := bodyMap["sku"]; ok {
		if val, ok := bodyMap["sku"].(string); ok {
			existing.SKU = val
		}
	}
	if _, ok := bodyMap["brand"]; ok {
		if val, ok := bodyMap["brand"].(string); ok {
			existing.Brand = val
		}
	}
	if _, ok := bodyMap["isHot"]; ok {
		if val, ok := bodyMap["isHot"].(bool); ok {
			existing.IsHot = val
		}
	}
	if _, ok := bodyMap["is_hot"]; ok {
		if val, ok := bodyMap["is_hot"].(bool); ok {
			existing.IsHot = val
		}
	}
	if _, ok := bodyMap["isActive"]; ok {
		if val, ok := bodyMap["isActive"].(bool); ok {
			existing.IsActive = val
		}
	}
	if _, ok := bodyMap["is_active"]; ok {
		if val, ok := bodyMap["is_active"].(bool); ok {
			existing.IsActive = val
		}
	}
	if _, ok := bodyMap["purchaseLimit"]; ok {
		if val, ok := bodyMap["purchaseLimit"].(float64); ok {
			existing.PurchaseLimit = int(val)
		}
	}
	if _, ok := bodyMap["purchase_limit"]; ok {
		if val, ok := bodyMap["purchase_limit"].(float64); ok {
			existing.PurchaseLimit = int(val)
		}
	}
	if _, ok := bodyMap["purchaseWarning"]; ok {
		if val, ok := bodyMap["purchaseWarning"].(string); ok {
			existing.PurchaseWarning = val
		}
	}
	if _, ok := bodyMap["purchase_warning"]; ok {
		if val, ok := bodyMap["purchase_warning"].(string); ok {
			existing.PurchaseWarning = val
		}
	}
	if _, ok := bodyMap["visibilityLevel"]; ok {
		if val, ok := bodyMap["visibilityLevel"].(float64); ok {
			existing.VisibilityLevel = int(val)
		}
	}
	if _, ok := bodyMap["visibility_level"]; ok {
		if val, ok := bodyMap["visibility_level"].(float64); ok {
			existing.VisibilityLevel = int(val)
		}
	}
	if _, ok := bodyMap["specs"]; ok {
		if arr, ok := bodyMap["specs"].([]any); ok {
			jsBytes, _ := json.Marshal(arr)
			var specs []entity.ProductSpecItem
			if err := json.Unmarshal(jsBytes, &specs); err == nil {
				existing.Specs = specs
			}
		}
	}
}

// @Summary     Update admin product
// @Description Updates a product by ID
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

	productID := ctx.Params("id")
	existing, err := ext.GetProduct(ctx.UserContext(), r.gripActor(ctx), productID)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	contentType := strings.ToLower(ctx.Get("Content-Type"))
	isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

	var bodyMap map[string]any
	if !isMultipart && len(ctx.Body()) > 0 {
		_ = json.Unmarshal(ctx.Body(), &bodyMap)
	}

	mergeProductPatch(&existing, bodyMap, ctx, isMultipart)

	if specsStr := ctx.FormValue("specs"); specsStr != "" {
		var specs []entity.ProductSpecItem
		if err := json.Unmarshal([]byte(specsStr), &specs); err == nil {
			existing.Specs = specs
		}
	}

	updated, err := ext.UpsertProduct(ctx.UserContext(), r.gripActor(ctx), existing)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.JSON(apiSuccessEnvelope(updated))
}

// @Summary     Update admin product visibility status
// @Description Updates only the visibility state of a product by ID
// @ID          grip_admin_update_product_status
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param       id path string true "Product ID"
// @Success     200 {object} envelope
// @Failure     400 {object} envelope
// @Failure     403 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /admin/products/{id}/status [patch]
func (r *V1) gripAdminUpdateProductStatus(ctx *fiber.Ctx) error {
	if r.adminUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_usecase_not_configured"})
	}
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_products_not_available"})
	}

	var body struct {
		IsActive *bool `json:"isActive"`
	}
	if err := ctx.BodyParser(&body); err != nil || body.IsActive == nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	productID := ctx.Params("id")
	existing, err := ext.GetProduct(ctx.UserContext(), r.gripActor(ctx), productID)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	existing.IsActive = *body.IsActive
	updated, err := ext.UpsertProduct(ctx.UserContext(), r.gripActor(ctx), existing)
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

	pageNum := max(ctx.QueryInt("page", 1), 1)
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
	timeline := []fiber.Map{}
	isTerminal := order.Status == entity.OrderStatusCancelled ||
		order.Status == entity.OrderStatusFailed ||
		order.Status == entity.OrderStatusRefundPending ||
		order.Status == entity.OrderStatusRefunded

	if isTerminal {
		timeline = append(timeline, fiber.Map{
			"status":    status,
			"timestamp": order.UpdatedAt,
			"note":      "",
		})
		if order.DeliveredAt != nil {
			timeline = append(timeline, fiber.Map{
				"status":    "DELIVERED",
				"timestamp": *order.DeliveredAt,
				"note":      "",
			})
		}
		if order.PaidAt != nil {
			timeline = append(timeline, fiber.Map{
				"status":    "PAID",
				"timestamp": *order.PaidAt,
				"note":      "",
			})
		}
		timeline = append(timeline, fiber.Map{
			"status":    "PENDING",
			"timestamp": order.CreatedAt,
			"note":      "",
		})
	} else {
		timeline = append(timeline, fiber.Map{
			"status":    "PENDING",
			"timestamp": order.CreatedAt,
			"note":      "",
		})
		if order.PaidAt != nil {
			timeline = append(timeline, fiber.Map{
				"status":    "PAID",
				"timestamp": *order.PaidAt,
				"note":      "",
			})
		}
		if order.DeliveredAt != nil {
			timeline = append(timeline, fiber.Map{
				"status":    "DELIVERED",
				"timestamp": *order.DeliveredAt,
				"note":      "",
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"id":          order.ID,
		"orderNumber": order.ID,
		"status":      status,
		"createdAt":   order.CreatedAt,
		"paidAt":      order.PaidAt,
		"deliveredAt": order.DeliveredAt,
		"updatedAt":   order.UpdatedAt,
		"tradeNo":     order.TradeNo,
		"pointsUsed":  order.PointsUsed,
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
		"timeline":        timeline,
	})
}

func (r *V1) gripAdminUpdateOrder(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_orders_not_available"})
	}

	var body struct {
		Status entity.OrderStatus `json:"status"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	if err := ext.UpdateOrderStatus(ctx.UserContext(), r.gripActor(ctx), ctx.Params("id"), body.Status); err != nil {
		httpStatus, payload := mapDomainError(err)
		return ctx.Status(httpStatus).JSON(payload)
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) gripAdminDeleteOrder(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_orders_not_available"})
	}

	if err := ext.DeleteOrder(ctx.UserContext(), r.gripActor(ctx), ctx.Params("id")); err != nil {
		httpStatus, payload := mapDomainError(err)
		return ctx.Status(httpStatus).JSON(payload)
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) gripAdminListRefunds(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_refunds_not_available"})
	}

	refunds, err := ext.ListRefunds(ctx.UserContext(), r.gripActor(ctx), strings.TrimSpace(ctx.Query("status", "pending")))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(refunds))
}

func (r *V1) gripAdminApproveRefund(ctx *fiber.Ctx) error {
	return r.gripAdminDecideRefund(ctx, true)
}

func (r *V1) gripAdminRejectRefund(ctx *fiber.Ctx) error {
	return r.gripAdminDecideRefund(ctx, false)
}

func (r *V1) gripAdminDecideRefund(ctx *fiber.Ctx, approve bool) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_refunds_not_available"})
	}

	refundID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	var body struct {
		Note string `json:"note"`
	}
	if len(ctx.Body()) > 0 {
		if err := ctx.BodyParser(&body); err != nil {
			status, payload := mapDomainError(entity.ErrInvalidInput)
			return ctx.Status(status).JSON(payload)
		}
	}

	refund, err := ext.DecideRefund(ctx.UserContext(), r.gripActor(ctx), refundID, approve, body.Note)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(refund))
}

func (r *V1) gripAdminNotificationTest(ctx *fiber.Ctx) error {
	var body struct {
		Channel string `json:"channel"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if len(ctx.Body()) > 0 {
		if err := ctx.BodyParser(&body); err != nil {
			status, payload := mapDomainError(entity.ErrInvalidInput)
			return ctx.Status(status).JSON(payload)
		}
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{
		"status":  "queued",
		"channel": body.Channel,
	}))
}

func (r *V1) gripAdminImportData(ctx *fiber.Ctx) error {
	if r.importer == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "importer_not_configured"})
	}

	var body importBody
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	result, err := r.importer.Import(ctx.UserContext(), body.Items)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(result))
}

func (r *V1) gripAdminListSettings(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_settings_not_available"})
	}

	settings, err := ext.ListSettings(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(settings))
}

func (r *V1) gripAdminUpsertSetting(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_settings_not_available"})
	}

	var body struct {
		Value string `json:"value"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	if err := ext.SetSetting(ctx.UserContext(), r.gripActor(ctx), ctx.Params("key"), body.Value); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{
		"key":   ctx.Params("key"),
		"value": body.Value,
	}))
}

func (r *V1) gripAdminDeleteSetting(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_settings_not_available"})
	}

	if err := ext.DeleteSetting(ctx.UserContext(), r.gripActor(ctx), ctx.Params("key")); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.SendStatus(http.StatusNoContent)
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
	uctx := ctx.UserContext()
	if q := ctx.Query("q"); q != "" {
		uctx = context.WithValue(uctx, "query", q)
	}
	if role := ctx.Query("role"); role != "" {
		uctx = context.WithValue(uctx, "role", role)
	}
	users, total, err := r.adminUC.ListUsers(uctx, r.gripActor(ctx), page)
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

func (r *V1) gripAdminUsersUpdatePoints(ctx *fiber.Ctx) error {
	var body struct {
		Points int `json:"points"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	actor := r.gripActor(ctx)
	if err := r.adminUC.UpdateUserPoints(ctx.UserContext(), actor, ctx.Params("id"), body.Points); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.SendStatus(http.StatusOK)
}

func (r *V1) gripAdminUsersUpdateBlock(ctx *fiber.Ctx) error {
	var body struct {
		IsBlocked bool `json:"isBlocked"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	actor := r.gripActor(ctx)
	status := entity.UserStatusActive
	if body.IsBlocked {
		status = entity.UserStatusLocked
	}
	if err := r.adminUC.UpdateUserStatus(ctx.UserContext(), actor, ctx.Params("id"), status); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.SendStatus(http.StatusOK)
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

func (r *V1) gripAdminProductsNew(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	var categories []entity.Category
	if ok {
		if items, err := ext.ListCategories(ctx.UserContext(), r.gripActor(ctx)); err == nil {
			categories = items
		}
	}
	if categories == nil {
		categories = []entity.Category{}
	}

	return ctx.JSON(fiber.Map{
		"product":    nil,
		"categories": categories,
	})
}

func (r *V1) gripAdminProductForm(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_products_not_available"})
	}

	product, err := ext.GetProduct(ctx.UserContext(), r.gripActor(ctx), ctx.Params("id"))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	var categories []entity.Category
	if items, err := ext.ListCategories(ctx.UserContext(), r.gripActor(ctx)); err == nil {
		categories = items
	}
	if categories == nil {
		categories = []entity.Category{}
	}

	return ctx.JSON(fiber.Map{
		"product":    product,
		"categories": categories,
	})
}

func (r *V1) gripAdminGetCollect(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_settings_not_available"})
	}

	settings, err := ext.ListSettings(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	payee := ""
	payLink := ""
	for _, setting := range settings {
		if setting.Key == "payee" {
			payee = setting.Value
		}
		if setting.Key == "payLink" {
			payLink = setting.Value
		}
	}

	warnings := make([]string, 0)
	if payee == "" {
		warnings = append(warnings, "Payee is not configured")
	}
	if len(payLink) < 10 {
		warnings = append(warnings, "PayLink must be at least 10 characters")
	}
	ready := payee != "" && len(payLink) >= 10

	sources := []fiber.Map{
		{
			"id":      "vcb",
			"key":     "vcb",
			"label":   "Vietcombank QR",
			"enabled": ready,
			"status":  "active",
		},
		{
			"id":      "momo",
			"key":     "momo",
			"label":   "MoMo",
			"enabled": false,
			"status":  "inactive",
		},
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{
		"payee":    payee,
		"payLink":  payLink,
		"ready":    ready,
		"is_ready": ready,
		"warnings": warnings,
		"sources":  sources,
	}))
}

func (r *V1) gripAdminPutCollect(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_settings_not_available"})
	}

	var body struct {
		Payee   string `json:"payee"`
		PayLink string `json:"payLink"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	if len(body.PayLink) < 10 || body.Payee == "" {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	if err := ext.SetSetting(ctx.UserContext(), r.gripActor(ctx), "payee", body.Payee); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	if err := ext.SetSetting(ctx.UserContext(), r.gripActor(ctx), "payLink", body.PayLink); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{
		"payee":   body.Payee,
		"payLink": body.PayLink,
	}))
}

func (r *V1) gripAdminGetRefund(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_refunds_not_available"})
	}

	refundID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	refund, err := ext.GetRefund(ctx.UserContext(), r.gripActor(ctx), refundID)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	tradeNo := refund.TradeNo
	if tradeNo == "" {
		tradeNo = "AUTO-REFUND-TRADE-" + refund.OrderID
	}

	dataMap := fiber.Map{
		"id":             refund.ID,
		"order_id":       refund.OrderID,
		"user_id":        refund.UserID,
		"username":       refund.Username,
		"reason":         refund.Reason,
		"status":         refund.Status,
		"admin_username": refund.AdminUsername,
		"admin_note":     refund.AdminNote,
		"product_name":   refund.ProductName,
		"amount":         refund.Amount,
		"points_used":    refund.PointsUsed,
		"trade_no":       tradeNo,
		"order_status":   refund.OrderStatus,
		"created_at":     refund.CreatedAt,
		"updated_at":     refund.UpdatedAt,
	}
	if refund.ProcessedAt != nil {
		dataMap["processed_at"] = refund.ProcessedAt
	}

	res := fiber.Map{}
	for k, v := range dataMap {
		res[k] = v
	}
	res["data"] = dataMap

	return ctx.JSON(res)
}

func (r *V1) gripAdminGetOrderRefundStatus(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_refunds_not_available"})
	}

	orderID := ctx.Params("id")
	if strings.TrimSpace(orderID) == "" {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	refund, err := ext.GetOrderRefundStatus(ctx.UserContext(), r.gripActor(ctx), orderID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return ctx.JSON(fiber.Map{
				"success":          true,
				"hasPendingRefund": false,
				"status":           1,
				"msg":              "No pending refund request",
			})
		}
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(fiber.Map{
		"success":          true,
		"hasPendingRefund": true,
		"refundId":         refund.ID,
		"status":           1,
		"msg":              "Pending refund request exists",
	})
}

func (r *V1) gripAdminListCards(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_products_not_available"})
	}

	cards, err := ext.ListCards(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(cards))
}

func (r *V1) gripAdminGetNotifications(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_settings_not_available"})
	}

	settings, err := ext.ListSettings(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	values := map[string]string{}
	for _, setting := range settings {
		values[setting.Key] = setting.Value
	}

	telegramEnabled := parseBoolSetting(values["telegramEnabled"], false)
	barkEnabled := parseBoolSetting(values["barkEnabled"], false)
	resendEnabled := parseBoolSetting(values["resendEnabled"], false)

	res := fiber.Map{
		"telegramBotToken": values["telegramBotToken"],
		"telegramChatId":   values["telegramChatId"],
		"telegramLanguage": firstNonEmpty(values["telegramLanguage"], "vi"),
		"telegramEnabled":  telegramEnabled,
		"barkEnabled":      barkEnabled,
		"barkServerUrl":    firstNonEmpty(values["barkServerUrl"], "https://api.day.app"),
		"barkDeviceKey":    values["barkDeviceKey"],
		"resendApiKey":     values["resendApiKey"],
		"resendFromEmail":  values["resendFromEmail"],
		"resendFromName":   values["resendFromName"],
		"resendEnabled":    resendEnabled,
		"emailLanguage":    firstNonEmpty(values["emailLanguage"], "vi"),
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{
		"telegramEnabled": telegramEnabled,
		"barkEnabled":     barkEnabled,
		"resendEnabled":   resendEnabled,
		"settings":        res,
	}))
}

func (r *V1) gripAdminPostNotifications(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_settings_not_available"})
	}

	var body map[string]string
	_ = ctx.BodyParser(&body)

	getVal := func(key string) string {
		if val := ctx.FormValue(key); val != "" {
			return val
		}
		if body != nil {
			return body[key]
		}
		return ""
	}

	keys := []string{
		"telegramBotToken", "telegramChatId", "telegramLanguage", "telegramEnabled",
		"barkEnabled", "barkServerUrl", "barkDeviceKey",
		"resendApiKey", "resendFromEmail", "resendFromName", "resendEnabled", "emailLanguage",
	}

	for _, k := range keys {
		v := getVal(k)
		if v != "" {
			if err := ext.SetSetting(ctx.UserContext(), r.gripActor(ctx), k, v); err != nil {
				status, payload := mapDomainError(err)
				return ctx.Status(status).JSON(payload)
			}
		}
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{"success": true}))
}

func (r *V1) gripAdminListMessages(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_messages_not_available"})
	}

	msgs, err := ext.ListAdminMessages(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	rows := make([]fiber.Map, 0, len(msgs))
	for _, m := range msgs {
		status := "sent"
		rows = append(rows, fiber.Map{
			"id":          m.ID,
			"title":       m.Title,
			"subject":     m.Title,
			"body":        m.Body,
			"targetType":  m.TargetType,
			"targetValue": m.TargetValue,
			"sender":      m.Sender,
			"status":      status,
			"result":      status,
			"outcome":     status,
			"sentAt":      m.CreatedAt.Format(time.RFC3339),
			"sent_at":     m.CreatedAt.Format(time.RFC3339),
			"createdAt":   m.CreatedAt.Format(time.RFC3339),
			"created_at":  m.CreatedAt.Format(time.RFC3339),
		})
	}

	return ctx.JSON(apiSuccessEnvelope(rows))
}
