package v1

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

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
	article, err := r.content.GetArticle(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, "article not found")
	}
	if err := mergeArticlePatch(&article, ctx); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	article, err = r.content.UpdateArticle(ctx.UserContext(), article)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(article)
}

func mergeArticlePatch(article *entity.ContentArticle, ctx *fiber.Ctx) error {
	if article == nil {
		return entity.ErrInvalidInput
	}
	patch := map[string]json.RawMessage{}
	if body := strings.TrimSpace(string(ctx.Body())); body != "" {
		if err := json.Unmarshal([]byte(body), &patch); err != nil {
			return err
		}
	}
	decode := func(target any, names ...string) error {
		for _, name := range names {
			if raw, ok := patch[name]; ok {
				return json.Unmarshal(raw, target)
			}
		}
		return nil
	}
	if err := decode(&article.Title, "title"); err != nil {
		return err
	}
	if err := decode(&article.Slug, "slug"); err != nil {
		return err
	}
	if err := decode(&article.Body, "body"); err != nil {
		return err
	}
	if err := decode(&article.Status, "status"); err != nil {
		return err
	}
	if err := decode(&article.ScheduledAt, "scheduled_at", "scheduledAt"); err != nil {
		return err
	}
	if err := decode(&article.PublishedAt, "published_at", "publishedAt"); err != nil {
		return err
	}
	if err := decode(&article.AuthorID, "author_id", "authorId"); err != nil {
		return err
	}
	if err := decode(&article.ImageURL, "image_url", "imageUrl"); err != nil {
		return err
	}
	if err := decode(&article.Tags, "tags"); err != nil {
		return err
	}
	if err := decode(&article.Topic, "topic"); err != nil {
		return err
	}
	if err := decode(&article.Priority, "priority"); err != nil {
		return err
	}
	if strings.HasSuffix(ctx.Path(), "/publish") {
		article.Status = entity.ContentStatusPublished
	}
	if strings.HasSuffix(ctx.Path(), "/schedule") {
		article.Status = entity.ContentStatusScheduled
	}
	if article.Status == entity.ContentStatusPublished && article.PublishedAt == nil {
		publishedAt := time.Now().UTC()
		article.PublishedAt = &publishedAt
	}
	return nil
}

func (r *V1) previewArticle(ctx *fiber.Ctx) error {
	article, err := r.content.GetArticle(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, "article not found")
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
	filter := entity.ArticleFilter{
		PublicOnly: publicOnly,
		Topic:      ctx.Query("topic"),
		Tag:        ctx.Query("tag"),
		Pagination: queryPage(ctx),
	}
	items, total, err := r.content.ListArticles(ctx.UserContext(), filter)
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

func (r *V1) deleteArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := r.content.DeleteArticle(ctx.UserContext(), id); err != nil {
		return errorResponse(ctx, http.StatusNotFound, "article not found")
	}
	return ctx.SendStatus(http.StatusNoContent)
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
