package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type listResponse struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

func (r *V1) registerEcommerceRoutes(apiV1Group fiber.Router, protected fiber.Router) {
	adminCatalog := protected.Group("/catalog")
	adminCatalog.Get("/products", r.listProducts)
	adminCatalog.Post("/products", r.createProduct)
	adminCatalog.Get("/products/:id", r.getProduct)
	adminCatalog.Patch("/products/:id", r.updateProduct)
	adminCatalog.Delete("/products/:id", r.deleteProduct)
	adminCatalog.Get("/categories", r.listCategories)
	adminCatalog.Post("/categories", r.createCategory)
	adminCatalog.Get("/tags", r.listTags)
	adminCatalog.Post("/tags", r.createTag)

	protected.Post("/media", r.createMedia)
	protected.Get("/media", r.listMedia)
	protected.Delete("/media/:id", r.deleteMedia)

	protected.Get("/homepage/blocks", r.listHomepageBlocks)
	protected.Post("/homepage/blocks", r.createHomepageBlock)
	protected.Patch("/homepage/blocks/:id", r.updateHomepageBlock)
	protected.Delete("/homepage/blocks/:id", r.deleteHomepageBlock)
	protected.Get("/support/channels", r.listSupportChannels)
	protected.Patch("/support/channels/:id", r.updateSupportChannel)

	protected.Get("/content/articles", r.listAdminArticles)
	protected.Post("/content/articles", r.createArticle)
	protected.Patch("/content/articles/:id", r.updateArticle)
	protected.Post("/content/articles/:id/schedule", r.updateArticle)
	protected.Post("/content/articles/:id/publish", r.updateArticle)
	protected.Get("/content/pages", r.getPage)
	protected.Post("/content/pages", r.createPage)
	protected.Patch("/content/pages/:slug", r.updatePage)
	protected.Post("/import/initial-content", r.importInitialContent)

	public := apiV1Group.Group("/public")
	public.Get("/search", r.listProducts)
	public.Get("/categories", r.listCategories)
	public.Get("/categories/:id/products", r.listProducts)
	public.Get("/products/:id", r.getProduct)
	public.Get("/homepage", r.listPublicHomepage)
	public.Get("/content/articles", r.listPublicArticles)
	public.Get("/content/articles/:id", r.getArticle)
	public.Get("/content/pages/:slug", r.getPage)
	public.Get("/footer", r.listPublicHomepage)
	public.Get("/support", r.listPublicSupport)

	apiV1Group.Post("/cart", r.createCart)
	apiV1Group.Get("/cart/:session_id", r.getCart)
	apiV1Group.Post("/cart/:session_id/items", r.addCartItem)
	apiV1Group.Patch("/cart/:session_id/items/:item_id", r.updateCartItem)
	apiV1Group.Delete("/cart/:session_id/items/:item_id", r.removeCartItem)
	apiV1Group.Post("/order-requests", r.submitOrder)
	apiV1Group.Post("/leads", r.submitLead)
	protected.Get("/leads/:id", r.getLead)
}

func (r *V1) listProducts(ctx *fiber.Ctx) error {
	items, total, err := r.catalog.ListProducts(ctx.UserContext(), entity.ProductFilter{Keyword: ctx.Query("q"), Brand: ctx.Query("brand"), Sort: ctx.Query("sort"), Pagination: queryPage(ctx)})
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(listResponse{Data: items, Meta: entity.Page{Limit: queryPage(ctx).Normalize().Limit, Offset: queryPage(ctx).Normalize().Offset, Total: total}})
}

func (r *V1) createProduct(ctx *fiber.Ctx) error {
	var product entity.Product
	if err := ctx.BodyParser(&product); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	product, err := r.catalog.CreateProduct(ctx.UserContext(), product)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusCreated).JSON(product)
}

func (r *V1) getProduct(ctx *fiber.Ctx) error {
	product, err := r.catalog.GetProduct(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, "product not found")
	}
	return ctx.JSON(product)
}

func (r *V1) updateProduct(ctx *fiber.Ctx) error {
	var product entity.Product
	if err := ctx.BodyParser(&product); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	product.ID = ctx.Params("id")
	product, err := r.catalog.UpdateProduct(ctx.UserContext(), product)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(product)
}

func (r *V1) deleteProduct(ctx *fiber.Ctx) error {
	if err := r.catalog.DeleteProduct(ctx.UserContext(), ctx.Params("id")); err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) createCategory(ctx *fiber.Ctx) error {
	var category entity.Category
	if err := ctx.BodyParser(&category); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	category, err := r.catalog.CreateCategory(ctx.UserContext(), category)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusCreated).JSON(category)
}

func (r *V1) listCategories(ctx *fiber.Ctx) error {
	items, err := r.catalog.ListCategories(ctx.UserContext())
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(listResponse{Data: items})
}

func (r *V1) createTag(ctx *fiber.Ctx) error {
	var tag entity.Tag
	if err := ctx.BodyParser(&tag); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	tag, err := r.catalog.CreateTag(ctx.UserContext(), tag)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusCreated).JSON(tag)
}

func (r *V1) listTags(ctx *fiber.Ctx) error {
	items, err := r.catalog.ListTags(ctx.UserContext())
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(listResponse{Data: items})
}

func queryPage(ctx *fiber.Ctx) entity.Pagination {
	return entity.Pagination{Limit: ctx.QueryInt("limit", 20), Offset: ctx.QueryInt("offset", 0)}
}
