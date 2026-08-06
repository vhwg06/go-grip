package catalog

import (
	"context"
	"encoding/json"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/module/catalog/catalogbase"
)

// GetCatalogProductModel handles GET /catalog/product-models/{id}.
func (h *Handler) GetCatalogProductModel(ctx context.Context, request openapi.GetCatalogProductModelRequestObject) (openapi.GetCatalogProductModelResponseObject, error) {
	result, err := h.catalogBase.GetPublicModel(ctx, request.Id)
	if err != nil {
		status, _ := catalogbase.ErrorStatus(err)
		if status == 404 {
			return openapi.GetCatalogProductModel404JSONResponse{}, nil
		}
		return openapi.GetCatalogProductModel500JSONResponse{}, nil
	}
	return openapi.GetCatalogProductModel200JSONResponse(result), nil
}

// ListCatalogProductModels handles GET /catalog/product-models
func (h *Handler) ListCatalogProductModels(ctx context.Context, request openapi.ListCatalogProductModelsRequestObject) (openapi.ListCatalogProductModelsResponseObject, error) {
	limit := 20
	page := 1
	if request.Params.Page != nil && *request.Params.Page > 0 {
		page = *request.Params.Page
	}
	if request.Params.Limit != nil && *request.Params.Limit > 0 {
		limit = *request.Params.Limit
	}

	filter := catalogbase.PublicFilter{
		Page:  page,
		Limit: limit,
	}
	if request.Params.MinPrice != nil {
		min := int64(*request.Params.MinPrice)
		filter.MinPrice = &min
	}
	if request.Params.MaxPrice != nil {
		max := int64(*request.Params.MaxPrice)
		filter.MaxPrice = &max
	}
	if request.Params.CategoryId != nil {
		filter.CategoryID = *request.Params.CategoryId
	}
	if request.Params.MaterialId != nil {
		filter.MaterialID = *request.Params.MaterialId
	}
	if request.Params.FinishId != nil {
		filter.FinishID = *request.Params.FinishId
	}
	if request.Params.Search != nil {
		filter.Search = *request.Params.Search
	}
	if request.Params.Sort != nil {
		filter.Sort = *request.Params.Sort
	}

	resMap, err := h.catalogBase.ListPublicModels(ctx, filter)
	if err != nil {
		return openapi.ListCatalogProductModels500JSONResponse{}, nil
	}

	resp := toPublicProductModelListResponse(resMap)
	return openapi.ListCatalogProductModels200JSONResponse(resp), nil
}

// GetCatalogProductModelOptions handles GET /catalog/product-models/{id}/options
func (h *Handler) GetCatalogProductModelOptions(ctx context.Context, request openapi.GetCatalogProductModelOptionsRequestObject) (openapi.GetCatalogProductModelOptionsResponseObject, error) {
	selected := map[string]string{}
	if request.Params.Selected != nil && *request.Params.Selected != "" {
		value := *request.Params.Selected
		if err := json.Unmarshal([]byte(value), &selected); err != nil {
			selected["Size"] = value
		}
	}

	result, err := h.catalogBase.AvailableOptions(ctx, request.Id, selected)
	if err != nil {
		return openapi.GetCatalogProductModelOptions404JSONResponse{}, nil
	}
	return openapi.GetCatalogProductModelOptions200JSONResponse(result), nil
}

// ResolveCatalogProductModelVariant handles POST /catalog/product-models/{id}/variants:resolve
func (h *Handler) ResolveCatalogProductModelVariant(ctx context.Context, request openapi.ResolveCatalogProductModelVariantRequestObject) (openapi.ResolveCatalogProductModelVariantResponseObject, error) {
	selected := map[string]string{}
	if request.Body != nil {
		bodyMap := map[string]any(*request.Body)
		if selMap, ok := bodyMap["selectedOptions"].(map[string]any); ok {
			for k, v := range selMap {
				if str, ok := v.(string); ok {
					selected[k] = str
				}
			}
		}
	}

	result, err := h.catalogBase.ResolvePublicVariant(ctx, request.Id, selected)
	if err != nil {
		status, _ := catalogbase.ErrorStatus(err)
		if status == 404 {
			return openapi.ResolveCatalogProductModelVariant404JSONResponse{}, nil
		}
		return openapi.ResolveCatalogProductModelVariant400JSONResponse{}, nil
	}
	return openapi.ResolveCatalogProductModelVariant200JSONResponse(result), nil
}

// SearchCatalog handles GET /catalog/search
func (h *Handler) SearchCatalog(ctx context.Context, request openapi.SearchCatalogRequestObject) (openapi.SearchCatalogResponseObject, error) {
	q := ""
	if request.Params.Q != nil {
		q = *request.Params.Q
	}
	filter := catalogbase.PublicFilter{
		Search: q,
	}
	resMap, err := h.catalogBase.ListPublicModels(ctx, filter)
	if err != nil {
		return openapi.SearchCatalog500JSONResponse{}, nil
	}
	items := []map[string]any{}
	if rawItems, ok := resMap["items"].([]map[string]any); ok {
		items = rawItems
	}
	return openapi.SearchCatalog200JSONResponse(items), nil
}

// GetCatalogProductBuyMeta handles GET /catalog/products/{id}/buy-meta
func (h *Handler) GetCatalogProductBuyMeta(ctx context.Context, request openapi.GetCatalogProductBuyMetaRequestObject) (openapi.GetCatalogProductBuyMetaResponseObject, error) {
	model, err := h.catalogBase.GetPublicModel(ctx, request.Id)
	if err != nil {
		return openapi.GetCatalogProductBuyMeta404JSONResponse{}, nil
	}
	return openapi.GetCatalogProductBuyMeta200JSONResponse{
		"productModelId": request.Id,
		"productModel":   model,
	}, nil
}

// GetCatalogAnnouncement handles GET /catalog/announcement
func (h *Handler) GetCatalogAnnouncement(ctx context.Context, _ openapi.GetCatalogAnnouncementRequestObject) (openapi.GetCatalogAnnouncementResponseObject, error) {
	return openapi.GetCatalogAnnouncement200JSONResponse{
		Enabled: true,
		Message: "Welcome to Grip Store",
	}, nil
}
