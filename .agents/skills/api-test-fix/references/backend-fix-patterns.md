# Backend Fix Patterns in Go-Grip

Use these established patterns when modifying delivery handlers in `internal/controller/restapi/v1/`:

---

## Pattern 1: Complete DTO Projection Mapping
When returning response objects, ensure every field from the domain entity is populated into the OpenAPI DTO struct:

```go
func toArticleResponse(a contentmodule.Article) openapi.ArticleResponse {
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
```

---

## Pattern 2: Partial PATCH Update Merge
For `PATCH` endpoints, load stored state first, overlay non-nil request body properties, then persist:

```go
func (h *Handler) UpdateContentArticle(ctx context.Context, request openapi.UpdateContentArticleRequestObject) (openapi.UpdateContentArticleResponseObject, error) {
	// 1. Fetch existing entity
	article, err := h.contentUC.GetArticle(ctx, request.Id)
	if err != nil {
		return openapi.UpdateContentArticle500JSONResponse{}, nil
	}

	// 2. Overlay non-nil request fields
	if request.Body != nil {
		if request.Body.Title != nil {
			article.Title = *request.Body.Title
		}
		if request.Body.Status != nil {
			article.Status = contentmodule.ContentStatus(*request.Body.Status)
		}
		if request.Body.Priority != nil {
			article.Priority = *request.Body.Priority
		}
		if request.Body.Tags != nil {
			article.Tags = append([]string(nil), (*request.Body.Tags)...)
		}
	}

	// 3. Persist updated model
	updated, err := h.contentUC.UpdateArticle(ctx, article)
	if err != nil {
		return openapi.UpdateContentArticle500JSONResponse{}, nil
	}

	return openapi.UpdateContentArticle200JSONResponse(toArticleResponse(updated)), nil
}
```

---

## Pattern 3: Query Parameter Filtering
Forward query parameters from `request.Params` into the usecase filter:

```go
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
	for _, a := range articles {
		items = append(items, toArticleResponse(a))
	}
	return openapi.ListPublicContentArticles200JSONResponse{Items: &items, Total: &total}, nil
}
```

---

## Pattern 4: Custom Response Formatting & Sorting
When the endpoint contract defines wrapper objects or custom sort orders:

```go
func settingsProjection(settings []catalogmodule.Setting) openapi.AdminStoreSettingsResponse {
	result := make(openapi.AdminStoreSettingsResponse, len(settings)+3)
	brand := map[string]any{}
	contact := map[string]any{}
	for _, setting := range settings {
		result[setting.Key] = setting.Value
		if strings.HasPrefix(setting.Key, "brand.") {
			brand[strings.TrimPrefix(setting.Key, "brand.")] = setting.Value
		}
	}
	result["config"] = map[string]any{"brand": brand, "contact": contact}
	result["stats"] = map[string]any{"settingsCount": len(settings)}
	result["visitorCount"] = 0
	return result
}
```

---

## Pattern 5: Validation & Error Mapper
Map invalid input / domain validation errors to HTTP 400 in `error_mapper.go`:

```go
if errors.Is(err, usermodule.ErrInvalidInput) || errors.Is(err, ordermodule.ErrInvalidInput) {
	resp := openapi.ErrorResponse{}
	resp.Error.FromErrorPayload(openapi.ErrorPayload{
		Code:    "INVALID_INPUT",
		Message: "Invalid administrative request",
	})
	return http.StatusBadRequest, resp
}
```
