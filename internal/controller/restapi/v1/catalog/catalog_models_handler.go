package catalog

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// ListCatalogProductModels handles GET /catalog/product-models
// Returns a list of product models for display in catalog browsing UIs.
func (h *Handler) ListCatalogProductModels(ctx context.Context, request openapi.ListCatalogProductModelsRequestObject) (openapi.ListCatalogProductModelsResponseObject, error) {
	limit := 20
	if request.Params.Limit != nil && *request.Params.Limit > 0 {
		limit = *request.Params.Limit
	}

	minPrice := 0
	maxPrice := 0
	if request.Params.MinPrice != nil {
		minPrice = *request.Params.MinPrice
	}
	if request.Params.MaxPrice != nil {
		maxPrice = *request.Params.MaxPrice
	}
	_ = minPrice
	_ = maxPrice

	filter := catalogmodule.ProductFilter{
		Pagination: pagination.Pagination{
			Limit:  limit,
			Offset: 0,
		},
	}

	products, total, err := h.catalogUC.ListProducts(ctx, filter)
	if err != nil {
		return openapi.ListCatalogProductModels500JSONResponse{}, nil
	}

	items := make([]map[string]interface{}, 0, len(products))
	for _, p := range products {
		items = append(items, map[string]interface{}{
			"id":          p.ID,
			"title":       p.Title,
			"price":       p.Price,
			"image_url":   p.ImageURL,
			"category_id": p.CategoryID,
		})
	}
	resp := openapi.PublicProductModelListResponse{Items: &items, Total: &total}
	return openapi.ListCatalogProductModels200JSONResponse(resp), nil
}
