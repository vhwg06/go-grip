package catalog

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/module/catalog/catalogbase"
)

// ListCatalogProductModels handles GET /catalog/product-models
func (h *Handler) ListCatalogProductModels(ctx context.Context, request openapi.ListCatalogProductModelsRequestObject) (openapi.ListCatalogProductModelsResponseObject, error) {
	limit := 20
	if request.Params.Limit != nil && *request.Params.Limit > 0 {
		limit = *request.Params.Limit
	}

	filter := catalogbase.PublicFilter{
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

	_, err := h.catalogBase.ListPublicModels(ctx, filter)
	if err != nil {
		return openapi.ListCatalogProductModels500JSONResponse{}, nil
	}

	items := make([]map[string]interface{}, 0)
	total := 0
	resp := openapi.PublicProductModelListResponse{
		Items: &items,
		Total: &total,
	}
	return openapi.ListCatalogProductModels200JSONResponse(resp), nil
}

// GetCatalogProductModelOptions handles GET /catalog/product-models/{id}/options
func (h *Handler) GetCatalogProductModelOptions(ctx context.Context, request openapi.GetCatalogProductModelOptionsRequestObject) (openapi.GetCatalogProductModelOptionsResponseObject, error) {
	selected := map[string]string{}
	if request.Params.Selected != nil && *request.Params.Selected != "" {
		selected["Size"] = *request.Params.Selected
	}

	_, err := h.catalogBase.AvailableOptions(ctx, request.Id, selected)
	if err != nil {
		return openapi.GetCatalogProductModelOptions404JSONResponse{}, nil
	}
	return openapi.GetCatalogProductModelOptions200Response{}, nil
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

	_, err := h.catalogBase.ResolvePublicVariant(ctx, request.Id, selected)
	if err != nil {
		return openapi.ResolveCatalogProductModelVariant400JSONResponse{}, nil
	}
	return openapi.ResolveCatalogProductModelVariant200Response{}, nil
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
	_, err := h.catalogBase.ListPublicModels(ctx, filter)
	if err != nil {
		return openapi.SearchCatalog500JSONResponse{}, nil
	}
	return openapi.SearchCatalog200Response{}, nil
}

// GetCatalogProductBuyMeta handles GET /catalog/products/{id}/buy-meta
func (h *Handler) GetCatalogProductBuyMeta(ctx context.Context, request openapi.GetCatalogProductBuyMetaRequestObject) (openapi.GetCatalogProductBuyMetaResponseObject, error) {
	_, err := h.catalogBase.GetPublicModel(ctx, request.Id)
	if err != nil {
		return openapi.GetCatalogProductBuyMeta404JSONResponse{}, nil
	}
	return openapi.GetCatalogProductBuyMeta200Response{}, nil
}

// GetCatalogAnnouncement handles GET /catalog/announcement
func (h *Handler) GetCatalogAnnouncement(ctx context.Context, _ openapi.GetCatalogAnnouncementRequestObject) (openapi.GetCatalogAnnouncementResponseObject, error) {
	return openapi.GetCatalogAnnouncement200Response{}, nil
}
