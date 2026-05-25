package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/controller/restapi/middleware"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type listResponse struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

func (r *V1) registerEcommerceRoutes(apiV1Group fiber.Router) {
	adminCatalog := apiV1Group.Group("/catalog", middleware.Auth(r.jwtManager), middleware.RequireAdminUsernames(r.adminUsers))
	adminCatalog.Post("/products", r.createProduct)
	adminCatalog.Patch("/products/:id", r.updateProduct)
	adminCatalog.Delete("/products/:id", r.deleteProduct)
	adminCatalog.Post("/categories", r.createCategory)
	adminCatalog.Post("/tags", r.createTag)

	apiV1Group.Post("/media", middleware.Auth(r.jwtManager), r.createMedia)
	apiV1Group.Get("/media", middleware.Auth(r.jwtManager), r.listMedia)
	apiV1Group.Delete("/media/:id", middleware.Auth(r.jwtManager), r.deleteMedia)
	apiV1Group.Put("/media/simulate-upload/:filename", r.simulateUpload)

	apiV1Group.Get("/homepage/blocks", middleware.Auth(r.jwtManager), r.listHomepageBlocks)
	apiV1Group.Post("/homepage/blocks", middleware.Auth(r.jwtManager), r.createHomepageBlock)
	apiV1Group.Patch("/homepage/blocks/:id", middleware.Auth(r.jwtManager), r.updateHomepageBlock)
	apiV1Group.Delete("/homepage/blocks/:id", middleware.Auth(r.jwtManager), r.deleteHomepageBlock)
	apiV1Group.Get("/support/channels", middleware.Auth(r.jwtManager), r.listSupportChannels)
	apiV1Group.Patch("/support/channels/:id", middleware.Auth(r.jwtManager), r.updateSupportChannel)

	apiV1Group.Get("/content/articles", middleware.Auth(r.jwtManager), r.listAdminArticles)
	apiV1Group.Post("/content/articles", middleware.Auth(r.jwtManager), r.createArticle)
	apiV1Group.Patch("/content/articles/:id", middleware.Auth(r.jwtManager), r.updateArticle)
	apiV1Group.Post("/content/articles/:id/schedule", middleware.Auth(r.jwtManager), r.updateArticle)
	apiV1Group.Post("/content/articles/:id/publish", middleware.Auth(r.jwtManager), r.updateArticle)
	apiV1Group.Get("/content/pages", middleware.Auth(r.jwtManager), r.getPage)
	apiV1Group.Post("/content/pages", middleware.Auth(r.jwtManager), r.createPage)
	apiV1Group.Patch("/content/pages/:slug", middleware.Auth(r.jwtManager), r.updatePage)
	apiV1Group.Post("/import/initial-content", middleware.Auth(r.jwtManager), r.importInitialContent)

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

	cartGroup := apiV1Group.Group("/cart", middleware.Auth(r.jwtManager))
	cartGroup.Post("/", r.createCart)
	cartGroup.Get("/:session_id", r.getCart)
	cartGroup.Post("/:session_id/items", r.addCartItem)
	cartGroup.Patch("/:session_id/items/:item_id", r.updateCartItem)
	cartGroup.Delete("/:session_id/items/:item_id", r.removeCartItem)
	apiV1Group.Post("/order-requests", middleware.Auth(r.jwtManager), r.submitOrder)
	apiV1Group.Post("/leads", r.submitLead)
	apiV1Group.Get("/leads/:id", middleware.Auth(r.jwtManager), r.getLead)
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
