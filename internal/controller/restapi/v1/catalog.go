package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type gripCatalogUseCase interface {
	ListVisibleProducts(ctx context.Context, actor entity.Actor, filter entity.ProductFilter) ([]entity.Product, int, error)
	GetVisibleProduct(ctx context.Context, actor entity.Actor, productID string) (entity.Product, error)
	ListPublicSettings(ctx context.Context) ([]entity.Setting, error)
	GetPublicSetting(ctx context.Context, key string) (entity.Setting, error)
}

type gripListResponse struct {
	Data any         `json:"data"`
	Meta entity.Page `json:"meta"`
}

type gripBuyMeta struct {
	ProductID       string          `json:"product_id"`
	ProductID2      string          `json:"productId"`
	Rating          float64         `json:"rating"`
	ReviewCount     int             `json:"reviewCount"`
	ReviewCount2    int             `json:"review_count"`
	CanReview       bool            `json:"canReview"`
	CanReview2      bool            `json:"can_review"`
	Available       bool            `json:"available"`
	Stock           int             `json:"stock"`
	Reviews         []entity.Review `json:"reviews"`
	AverageRating   float64         `json:"averageRating"`
	ReviewOrderID   *string         `json:"reviewOrderId"`
	EmailConfigured bool            `json:"emailConfigured"`
}

func (r *V1) gripCatalogUC() (gripCatalogUseCase, bool) {
	uc, ok := r.catalog.(gripCatalogUseCase)
	return uc, ok
}

func (r *V1) gripParseProductFilter(ctx *fiber.Ctx) (entity.ProductFilter, error) {
	minPriceRaw := ctx.Query("minPrice")
	if minPriceRaw == "" {
		minPriceRaw = ctx.Query("min_price")
	}

	maxPriceRaw := ctx.Query("maxPrice")
	if maxPriceRaw == "" {
		maxPriceRaw = ctx.Query("max_price")
	}

	filter := entity.ProductFilter{
		Keyword:    ctx.Query("q"),
		CategoryID: ctx.Query("category"),
		Brand:      ctx.Query("brand"),
		Sort:       ctx.Query("sort"),
		Pagination: gripPage(ctx),
	}

	if minPriceRaw != "" {
		minPrice, err := strconv.ParseInt(minPriceRaw, 10, 64)
		if err != nil {
			return entity.ProductFilter{}, entity.ErrInvalidInput
		}
		filter.MinPrice = &minPrice
	}
	if maxPriceRaw != "" {
		maxPrice, err := strconv.ParseInt(maxPriceRaw, 10, 64)
		if err != nil {
			return entity.ProductFilter{}, entity.ErrInvalidInput
		}
		filter.MaxPrice = &maxPrice
	}

	return filter, nil
}

// @Summary     List visible products
// @Description Lists active products visible to anonymous or signed-in users
// @ID          grip_catalog_list_products
// @Tags        catalog
// @Produce     json
// @Param       q query string false "Search keyword"
// @Param       brand query string false "Brand filter"
// @Param       min_price query int false "Minimum price"
// @Param       max_price query int false "Maximum price"
// @Param       limit query int false "Page size"
// @Param       offset query int false "Page offset"
// @Success     200 {object} gripListResponse
// @Failure     400 {object} envelope
// @Failure     500 {object} envelope
// @Router      /catalog/products [get]
func (r *V1) gripListProducts(ctx *fiber.Ctx) error {
	uc, ok := r.gripCatalogUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "catalog_usecase_not_configured"})
	}

	filter, err := r.gripParseProductFilter(ctx)
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	actor := r.gripActor(ctx)
	items, total, err := uc.ListVisibleProducts(ctx.UserContext(), actor, filter)
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	page := filter.Pagination.Normalize()
	pageNum := (page.Offset / page.Limit) + 1
	return ctx.JSON(fiber.Map{
		"items": items,
		"page":  pageNum,
		"limit": page.Limit,
		"total": total,
	})
}

// @Summary     Get visible product detail
// @Description Returns visible product detail and purchase metadata
// @ID          grip_catalog_get_product
// @Tags        catalog
// @Produce     json
// @Param       id path string true "Product ID"
// @Success     200 {object} entity.Product
// @Failure     404 {object} envelope
// @Failure     500 {object} envelope
// @Router      /catalog/products/{id} [get]
func (r *V1) gripGetProduct(ctx *fiber.Ctx) error {
	uc, ok := r.gripCatalogUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "catalog_usecase_not_configured"})
	}

	product, err := uc.GetVisibleProduct(ctx.UserContext(), r.gripActor(ctx), ctx.Params("id"))
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	return ctx.JSON(apiSuccessEnvelope(product))
}

// @Summary     Get buy metadata
// @Description Returns lightweight review metadata for purchase views
// @ID          grip_catalog_get_buy_meta
// @Tags        catalog
// @Produce     json
// @Param       id path string true "Product ID"
// @Success     200 {object} gripBuyMeta
// @Failure     404 {object} envelope
// @Failure     500 {object} envelope
// @Router      /catalog/products/{id}/buy-meta [get]
func (r *V1) gripGetBuyMeta(ctx *fiber.Ctx) error {
	uc, ok := r.gripCatalogUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "catalog_usecase_not_configured"})
	}

	product, err := uc.GetVisibleProduct(ctx.UserContext(), r.gripActor(ctx), ctx.Params("id"))
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	var reviews []entity.Review
	if r.wishlistUC != nil {
		if revs, err := r.wishlistUC.ListReviews(ctx.UserContext(), product.ID); err == nil {
			reviews = revs
		}
	}
	if reviews == nil {
		reviews = []entity.Review{}
	}

	available := product.IsActive && product.StockCount > 0

	return ctx.JSON(apiSuccessEnvelope(gripBuyMeta{
		ProductID:       product.ID,
		ProductID2:      product.ID,
		Rating:          product.Rating,
		ReviewCount:     product.ReviewCount,
		ReviewCount2:    product.ReviewCount,
		CanReview:       false,
		CanReview2:      false,
		Available:       available,
		Stock:           product.StockCount,
		Reviews:         reviews,
		AverageRating:   product.Rating,
		ReviewOrderID:   nil,
		EmailConfigured: false,
	}))
}

// @Summary     Search visible products
// @Description Searches active products visible to the actor
// @ID          grip_catalog_search_products
// @Tags        catalog
// @Produce     json
// @Param       q query string false "Search keyword"
// @Param       brand query string false "Brand filter"
// @Param       min_price query int false "Minimum price"
// @Param       max_price query int false "Maximum price"
// @Param       limit query int false "Page size"
// @Param       offset query int false "Page offset"
// @Success     200 {object} gripListResponse
// @Failure     400 {object} envelope
// @Failure     500 {object} envelope
// @Router      /catalog/search [get]
func (r *V1) gripSearchProducts(ctx *fiber.Ctx) error {
	uc, ok := r.gripCatalogUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "catalog_usecase_not_configured"})
	}

	filter, err := r.gripParseProductFilter(ctx)
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	actor := r.gripActor(ctx)
	items, _, err := uc.ListVisibleProducts(ctx.UserContext(), actor, filter)
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	return ctx.JSON(items)
}

// @Summary     List categories
// @Description Lists public categories for product navigation
// @ID          grip_catalog_list_categories
// @Tags        catalog
// @Produce     json
// @Success     200 {object} gripListResponse
// @Failure     500 {object} envelope
// @Router      /catalog/categories [get]
func (r *V1) gripListCategories(ctx *fiber.Ctx) error {
	items, err := r.catalog.ListCategories(ctx.UserContext())
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}
	response := make([]fiber.Map, 0, len(items))
	for _, item := range items {
		response = append(response, fiber.Map{
			"id":        item.ID,
			"name":      item.Name,
			"slug":      item.ID,
			"parent_id": item.ParentID,
			"position":  item.Position,
			"is_active": item.IsActive,
		})
	}
	return ctx.JSON(gripListResponse{Data: response})
}

// @Summary     List public settings
// @Description Lists public store settings used by the storefront
// @ID          grip_catalog_list_settings
// @Tags        catalog
// @Produce     json
// @Success     200 {object} gripListResponse
// @Failure     500 {object} envelope
// @Router      /catalog/settings [get]
func (r *V1) gripListSettings(ctx *fiber.Ctx) error {
	uc, ok := r.gripCatalogUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "catalog_usecase_not_configured"})
	}

	settings, err := uc.ListPublicSettings(ctx.UserContext())
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	return ctx.JSON(apiSuccessEnvelope(buildCatalogSettingsProjection(settings)))
}

// @Summary     Get active announcement
// @Description Returns the current announcement setting if configured
// @ID          grip_catalog_get_announcement
// @Tags        catalog
// @Produce     json
// @Success     200 {object} envelope
// @Failure     500 {object} envelope
// @Router      /catalog/announcement [get]
func (r *V1) gripGetAnnouncement(ctx *fiber.Ctx) error {
	uc, ok := r.gripCatalogUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "catalog_usecase_not_configured"})
	}

	setting, err := uc.GetPublicSetting(ctx.UserContext(), "announcement")
	if err != nil {
		setting, err = uc.GetPublicSetting(ctx.UserContext(), "test.announcement")
		if err != nil {
			return ctx.JSON(apiSuccessEnvelope(nil))
		}
	}

	res := fiber.Map{
		"id":         setting.Key,
		"content":    setting.Value,
		"active":     setting.Value != "",
		"updated_at": setting.UpdatedAt.UnixMilli(),
	}

	return ctx.JSON(apiSuccessEnvelope(res))
}
