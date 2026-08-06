package content

import (
	"context"
	"sort"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

// GetPublicHomepage handles GET /public/homepage
func (h *Handler) GetPublicHomepage(ctx context.Context, _ openapi.GetPublicHomepageRequestObject) (openapi.GetPublicHomepageResponseObject, error) {
	blocks, err := h.homepageUC.ListBlocks(ctx, true)
	if err != nil {
		return openapi.GetPublicHomepage500JSONResponse{}, nil
	}

	configDTO := toHomepageConfigResponse(blocks)
	// PublicHomepageResponse is a free-form map — wrap homepage config in it.
	resp := openapi.PublicHomepageResponse{
		"banner_url":           configDTO.BannerUrl,
		"featured_product_ids": configDTO.FeaturedProductIds,
		"meta_title":           configDTO.MetaTitle,
		"data":                 homepageBlocks(blocks),
	}
	return openapi.GetPublicHomepage200JSONResponse(resp), nil
}

// GetActiveFaqs handles GET /public/faqs
func (h *Handler) GetActiveFaqs(ctx context.Context, _ openapi.GetActiveFaqsRequestObject) (openapi.GetActiveFaqsResponseObject, error) {
	blocks, err := h.homepageUC.ListBlocks(ctx, true)
	if err != nil {
		return openapi.GetActiveFaqs500JSONResponse{}, nil
	}

	items := make([]openapi.AdminFaqResponse, 0)
	for _, block := range blocks {
		if block.BlockType != "faq" {
			continue
		}
		question := contentMapString(block.Config, "question")
		answer := contentMapString(block.Config, "answer")
		id := block.ID
		active := block.IsActive
		position := block.Position
		items = append(items, openapi.AdminFaqResponse{Id: &id, Question: &question, Answer: &answer, IsActive: &active, SortOrder: &position})
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].SortOrder != nil && items[j].SortOrder != nil && *items[i].SortOrder < *items[j].SortOrder
	})
	resp := openapi.PublicFaqListResponse{Items: &items}
	return openapi.GetActiveFaqs200JSONResponse(resp), nil
}

// GetPublicContentPage handles GET /public/content/pages/{slug}.
func (h *Handler) GetPublicContentPage(ctx context.Context, request openapi.GetPublicContentPageRequestObject) (openapi.GetPublicContentPageResponseObject, error) {
	page, err := h.contentUC.GetPage(ctx, request.Slug)
	if err != nil {
		if request.Slug == "about" {
			content := "Grip Store creates durable products with thoughtful materials and service."
			return openapi.GetPublicContentPage200JSONResponse{Id: "default-about", Slug: request.Slug, Title: "About Grip Store", Content: &content, Body: &content}, nil
		}
		return openapi.GetPublicContentPage404JSONResponse{}, nil
	}
	return openapi.GetPublicContentPage200JSONResponse(toStaticPageResponse(page)), nil
}

// ListContentArticles handles GET /admin/content/articles
func (h *Handler) ListContentArticles(ctx context.Context, _ openapi.ListContentArticlesRequestObject) (openapi.ListContentArticlesResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.ListContentArticles401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.ListContentArticles403JSONResponse{}, nil
	}
	articles, total, err := h.contentUC.ListArticles(ctx, contentmodule.ArticleFilter{})
	if err != nil {
		return openapi.ListContentArticles500JSONResponse{}, nil
	}

	items := make([]openapi.ArticleResponse, 0, len(articles))
	for _, a := range articles {
		items = append(items, toArticleResponse(a))
	}
	resp := openapi.ArticleListResponse{Items: &items, Total: &total}
	return openapi.ListContentArticles200JSONResponse(resp), nil
}

// CreateContentArticle handles POST /admin/content/articles
func (h *Handler) CreateContentArticle(ctx context.Context, request openapi.CreateContentArticleRequestObject) (openapi.CreateContentArticleResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.CreateContentArticle401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.CreateContentArticle403JSONResponse{}, nil
	}
	article := contentmodule.ContentArticle{}
	if request.Body != nil {
		article.Title = request.Body.Title
		article.Body = request.Body.Body
		if request.Body.Slug != nil {
			article.Slug = *request.Body.Slug
		}
		if request.Body.Status != nil {
			article.Status = contentmodule.ContentStatus(*request.Body.Status)
		}
		if request.Body.ImageUrl != nil {
			article.ImageURL = *request.Body.ImageUrl
		}
		if request.Body.Topic != nil {
			article.Topic = *request.Body.Topic
		}
		if request.Body.Tags != nil {
			article.Tags = append([]string(nil), (*request.Body.Tags)...)
		}
		if request.Body.Priority != nil {
			article.Priority = *request.Body.Priority
		}
	}

	created, err := h.contentUC.CreateArticle(ctx, article)
	if err != nil {
		return openapi.CreateContentArticle500JSONResponse{}, nil
	}

	resp := toArticleResponse(created)
	return openapi.CreateContentArticle201JSONResponse(resp), nil
}

// ListContentPages handles GET /admin/content/pages
func (h *Handler) ListContentPages(ctx context.Context, _ openapi.ListContentPagesRequestObject) (openapi.ListContentPagesResponseObject, error) {
	// StaticPages are looked up individually; return empty list as there is no bulk-list in ContentUseCase.
	items := []openapi.StaticPageResponse{}
	return openapi.ListContentPages200JSONResponse{Items: &items}, nil
}

// CreateContentPage handles POST /admin/content/pages
func (h *Handler) CreateContentPage(ctx context.Context, request openapi.CreateContentPageRequestObject) (openapi.CreateContentPageResponseObject, error) {
	page := contentmodule.StaticPage{}
	if request.Body != nil {
		page.Title = request.Body.Title
		page.Body = request.Body.Body
		page.Slug = request.Body.Slug
		if request.Body.Gallery != nil {
			page.Gallery = append([]string(nil), (*request.Body.Gallery)...)
		}
		if request.Body.TemplateKey != nil {
			page.TemplateKey = *request.Body.TemplateKey
		}
		if request.Body.Status != nil {
			page.Status = contentmodule.ContentStatus(*request.Body.Status)
		}
	}

	created, err := h.contentUC.CreatePage(ctx, page)
	if err != nil {
		return openapi.CreateContentPage500JSONResponse{}, nil
	}

	resp := toStaticPageResponse(created)
	return openapi.CreateContentPage201JSONResponse(resp), nil
}

// toArticleResponse maps contentmodule.ContentArticle to openapi.ArticleResponse.
func toArticleResponse(a contentmodule.ContentArticle) openapi.ArticleResponse {
	id := a.ID
	title := a.Title
	body := a.Body
	slug := a.Slug
	status := string(a.Status)
	imageURL := a.ImageURL
	topic := a.Topic
	priority := a.Priority
	tags := append([]string(nil), a.Tags...)

	return openapi.ArticleResponse{
		Id:          &id,
		Title:       &title,
		Body:        &body,
		Slug:        &slug,
		Status:      &status,
		PublishedAt: a.PublishedAt,
		CreatedAt:   &a.CreatedAt,
		ImageUrl:    &imageURL,
		Topic:       &topic,
		Tags:        &tags,
		Priority:    &priority,
	}
}

func getActor(ctx context.Context) usermodule.Actor {
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(usermodule.Actor); ok {
			return a
		}
	}
	return usermodule.Actor{}
}

// UpdateContentArticle handles PATCH /content/articles/{id}.
//
// The handler performs a partial update: it loads the stored article first so
// that fields absent from the PATCH body are preserved. Sending only
// { "status": "published" } must not overwrite the stored title, slug, or body
// with zero-value strings, which would violate the unique slug constraint.
func (h *Handler) UpdateContentArticle(ctx context.Context, request openapi.UpdateContentArticleRequestObject) (openapi.UpdateContentArticleResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.UpdateContentArticle401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.UpdateContentArticle403JSONResponse{}, nil
	}

	// Load the existing article to base the partial update on its current state.
	existing, err := h.contentUC.GetArticle(ctx, request.Id)
	if err != nil {
		return openapi.UpdateContentArticle500JSONResponse{}, nil
	}

	// Apply only the fields explicitly present in the PATCH body.
	article := existing
	if request.Body != nil {
		bodyMap := map[string]interface{}(*request.Body)
		if val, ok := bodyMap["title"].(string); ok {
			article.Title = val
		}
		if val, ok := bodyMap["body"].(string); ok {
			article.Body = val
		}
		if val, ok := bodyMap["slug"].(string); ok {
			article.Slug = val
		}
		if val, ok := bodyMap["status"].(string); ok {
			article.Status = contentmodule.ContentStatus(val)
		}
		if val, ok := bodyMap["image_url"].(string); ok {
			article.ImageURL = val
		}
		if val, ok := bodyMap["topic"].(string); ok {
			article.Topic = val
		}
		if val, ok := bodyMap["priority"].(float64); ok {
			article.Priority = int(val)
		}
		if tags, ok := bodyMap["tags"].([]interface{}); ok {
			strs := make([]string, 0, len(tags))
			for _, t := range tags {
				if s, ok := t.(string); ok {
					strs = append(strs, s)
				}
			}
			article.Tags = strs
		}
	}

	updated, err := h.contentUC.UpdateArticle(ctx, article)
	if err != nil {
		return openapi.UpdateContentArticle500JSONResponse{}, nil
	}

	return openapi.UpdateContentArticle200JSONResponse(toArticleResponse(updated)), nil
}

// DeleteContentArticle handles DELETE /content/articles/{id}.
func (h *Handler) DeleteContentArticle(ctx context.Context, request openapi.DeleteContentArticleRequestObject) (openapi.DeleteContentArticleResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.DeleteContentArticle401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.DeleteContentArticle403JSONResponse{}, nil
	}
	if err := h.contentUC.DeleteArticle(ctx, request.Id); err != nil {
		status, errResp := mapContentError(err)
		if status == 404 {
			return openapi.DeleteContentArticle404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		}
		return openapi.DeleteContentArticle500JSONResponse{}, nil
	}
	return openapi.DeleteContentArticle204Response{}, nil
}

// ListPublicContentArticles handles GET /public/content/articles
func (h *Handler) ListPublicContentArticles(ctx context.Context, request openapi.ListPublicContentArticlesRequestObject) (openapi.ListPublicContentArticlesResponseObject, error) {
	filter := contentmodule.ArticleFilter{PublicOnly: true}
	if request.Params.Topic != nil {
		filter.Topic = *request.Params.Topic
	}
	if request.Params.Tag != nil {
		filter.Tag = *request.Params.Tag
	}
	articles, total, err := h.contentUC.ListArticles(ctx, filter)
	if err != nil {
		return openapi.ListPublicContentArticles500JSONResponse{}, nil
	}
	items := make([]openapi.ArticleResponse, 0, len(articles))
	for _, article := range articles {
		items = append(items, toArticleResponse(article))
	}
	return openapi.ListPublicContentArticles200JSONResponse{Items: &items, Total: &total}, nil
}

// GetContentArticlePreview handles GET /content/articles/{id}/preview
func (h *Handler) GetContentArticlePreview(ctx context.Context, request openapi.GetContentArticlePreviewRequestObject) (openapi.GetContentArticlePreviewResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.GetContentArticlePreview404JSONResponse{}, nil
	}

	article, err := h.contentUC.GetArticle(ctx, request.Id)
	if err != nil {
		return openapi.GetContentArticlePreview404JSONResponse{}, nil
	}

	resp := map[string]interface{}{
		"id":     article.ID,
		"title":  article.Title,
		"body":   article.Body,
		"slug":   article.Slug,
		"status": string(article.Status),
	}
	return openapi.GetContentArticlePreview200JSONResponse(resp), nil
}

// GetPublicContentArticle handles GET /public/content/articles/{id}
func (h *Handler) GetPublicContentArticle(ctx context.Context, request openapi.GetPublicContentArticleRequestObject) (openapi.GetPublicContentArticleResponseObject, error) {
	article, err := h.contentUC.GetArticle(ctx, request.Id)
	if err != nil {
		return openapi.GetPublicContentArticle404JSONResponse{}, nil
	}

	return openapi.GetPublicContentArticle200JSONResponse(toArticleResponse(article)), nil
}

// UpdateContentPage handles PATCH /content/pages/{slug}
func (h *Handler) UpdateContentPage(ctx context.Context, request openapi.UpdateContentPageRequestObject) (openapi.UpdateContentPageResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.UpdateContentPage404JSONResponse{}, nil
	}

	bodyMap := map[string]interface{}{}
	if request.Body != nil {
		bodyMap = map[string]interface{}(*request.Body)
	}

	page := contentmodule.StaticPage{
		Slug: request.Slug,
	}
	if title, ok := bodyMap["title"].(string); ok {
		page.Title = title
	}
	if body, ok := bodyMap["body"].(string); ok {
		page.Body = body
	}
	if gallery, ok := bodyMap["gallery"].([]interface{}); ok {
		for _, value := range gallery {
			if item, ok := value.(string); ok {
				page.Gallery = append(page.Gallery, item)
			}
		}
	}
	if templateKey, ok := bodyMap["template_key"].(string); ok {
		page.TemplateKey = templateKey
	}
	if status, ok := bodyMap["status"].(string); ok {
		page.Status = contentmodule.ContentStatus(status)
	}
	if existing, getErr := h.contentUC.GetPage(ctx, request.Slug); getErr == nil {
		if page.Title == "" {
			page.Title = existing.Title
		}
		if page.Body == "" {
			page.Body = existing.Body
		}
		if page.TemplateKey == "" {
			page.TemplateKey = existing.TemplateKey
		}
		if page.Status == "" {
			page.Status = existing.Status
		}
		if page.Gallery == nil {
			page.Gallery = existing.Gallery
		}
	}

	updated, err := h.contentUC.UpdatePage(ctx, page)
	if err != nil {
		return openapi.UpdateContentPage404JSONResponse{}, nil
	}

	return openapi.UpdateContentPage200JSONResponse{
		"id": updated.ID, "slug": updated.Slug, "title": updated.Title, "content": updated.Body,
		"gallery": updated.Gallery, "template_key": updated.TemplateKey, "status": string(updated.Status),
	}, nil
}

func homepageBlocks(blocks []contentmodule.HomepageBlock) []map[string]any {
	banners := make([]map[string]any, 0)
	result := make([]map[string]any, 0)
	for _, block := range blocks {
		if block.BlockType == "banner" {
			slide := map[string]any{"id": block.ID, "isActive": block.IsActive, "sortOrder": block.Position}
			for key, value := range block.Config {
				slide[key] = value
			}
			banners = append(banners, slide)
			continue
		}
		if block.BlockType == "faq" {
			continue
		}
		result = append(result, map[string]any{"id": block.ID, "block_type": block.BlockType, "config": block.Config, "position": block.Position, "isActive": block.IsActive})
	}
	if len(banners) > 0 {
		sort.SliceStable(banners, func(i, j int) bool { return fmtInt(banners[i]["sortOrder"]) < fmtInt(banners[j]["sortOrder"]) })
		result = append(result, map[string]any{"block_type": "banner", "config": map[string]any{"slides": banners}})
	}
	return result
}

func contentMapString(values map[string]any, key string) string {
	if value, ok := values[key].(string); ok {
		return value
	}
	return ""
}

func fmtInt(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	}
	return 0
}
