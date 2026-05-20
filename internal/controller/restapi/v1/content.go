package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func (r *V1) createArticle(ctx *fiber.Ctx) error {
	var article entity.ContentArticle
	if err := ctx.BodyParser(&article); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	article, err := r.content.CreateArticle(ctx.UserContext(), article)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusCreated).JSON(article)
}

func (r *V1) updateArticle(ctx *fiber.Ctx) error {
	var article entity.ContentArticle
	if err := ctx.BodyParser(&article); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	article.ID = ctx.Params("id")
	article, err := r.content.UpdateArticle(ctx.UserContext(), article)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(article)
}

func (r *V1) listAdminArticles(ctx *fiber.Ctx) error {
	return r.listArticles(ctx, false)
}

func (r *V1) listPublicArticles(ctx *fiber.Ctx) error {
	return r.listArticles(ctx, true)
}

func (r *V1) listArticles(ctx *fiber.Ctx, publicOnly bool) error {
	items, total, err := r.content.ListArticles(ctx.UserContext(), publicOnly, queryPage(ctx))
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(listResponse{Data: items, Meta: entity.Page{Limit: queryPage(ctx).Normalize().Limit, Offset: queryPage(ctx).Normalize().Offset, Total: total}})
}

func (r *V1) getArticle(ctx *fiber.Ctx) error {
	article, err := r.content.GetArticle(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, "article not found")
	}
	return ctx.JSON(article)
}

func (r *V1) createPage(ctx *fiber.Ctx) error {
	var page entity.StaticPage
	if err := ctx.BodyParser(&page); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	page, err := r.content.CreatePage(ctx.UserContext(), page)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusCreated).JSON(page)
}

func (r *V1) updatePage(ctx *fiber.Ctx) error {
	var page entity.StaticPage
	if err := ctx.BodyParser(&page); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	page.Slug = ctx.Params("slug")
	page, err := r.content.UpdatePage(ctx.UserContext(), page)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(page)
}

func (r *V1) getPage(ctx *fiber.Ctx) error {
	page, err := r.content.GetPage(ctx.UserContext(), ctx.Params("slug"))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, "page not found")
	}
	return ctx.JSON(page)
}
