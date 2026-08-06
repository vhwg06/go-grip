package content

import (
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
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
	}
	return openapi.GetPublicHomepage200JSONResponse(resp), nil
}

// GetActiveFaqs handles GET /public/faqs
func (h *Handler) GetActiveFaqs(ctx context.Context, _ openapi.GetActiveFaqsRequestObject) (openapi.GetActiveFaqsResponseObject, error) {
	// Re-use article listing as FAQ source using the published filter.
	articles, _, err := h.contentUC.ListArticles(ctx, contentmodule.ArticleFilter{PublicOnly: true})
	if err != nil {
		return openapi.GetActiveFaqs500JSONResponse{}, nil
	}

	items := make([]openapi.AdminFaqResponse, 0, len(articles))
	for _, a := range articles {
		title := a.Title
		body := a.Body
		items = append(items, openapi.AdminFaqResponse{Question: &title, Answer: &body})
	}
	resp := openapi.PublicFaqListResponse{Items: &items}
	return openapi.GetActiveFaqs200JSONResponse(resp), nil
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

	return openapi.ArticleResponse{
		Id:          &id,
		Title:       &title,
		Body:        &body,
		Slug:        &slug,
		Status:      &status,
		PublishedAt: a.PublishedAt,
		CreatedAt:   &a.CreatedAt,
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


// UpdateContentArticle handles PATCH /content/articles/{id}
func (h *Handler) UpdateContentArticle(ctx context.Context, request openapi.UpdateContentArticleRequestObject) (openapi.UpdateContentArticleResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.UpdateContentArticle401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.UpdateContentArticle403JSONResponse{}, nil
	}

	article := contentmodule.ContentArticle{ID: request.Id}
	if request.Body != nil {
		bodyMap := map[string]interface{}(*request.Body)
		if val, ok := bodyMap["title"].(string); ok { article.Title = val }
		if val, ok := bodyMap["body"].(string); ok { article.Body = val }
		if val, ok := bodyMap["slug"].(string); ok { article.Slug = val }
		if val, ok := bodyMap["status"].(string); ok { article.Status = contentmodule.ContentStatus(val) }
	}

	_, err := h.contentUC.UpdateArticle(ctx, article)
	if err != nil {
		return openapi.UpdateContentArticle500JSONResponse{}, nil
	}

	return openapi.UpdateContentArticle200Response{}, nil
}

// ListPublicContentArticles handles GET /public/content/articles
func (h *Handler) ListPublicContentArticles(ctx context.Context, _ openapi.ListPublicContentArticlesRequestObject) (openapi.ListPublicContentArticlesResponseObject, error) {
	_, _, err := h.contentUC.ListArticles(ctx, contentmodule.ArticleFilter{PublicOnly: true})
	if err != nil {
		return openapi.ListPublicContentArticles500JSONResponse{}, nil
	}

	return openapi.ListPublicContentArticles200Response{}, nil
}
