package catalog

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
)

// AdminAddCatalogVariantDimension handles POST /admin/catalog/product-models/{modelId}/variant-dimensions
func (h *Handler) AdminAddCatalogVariantDimension(ctx context.Context, request openapi.AdminAddCatalogVariantDimensionRequestObject) (openapi.AdminAddCatalogVariantDimensionResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminAddCatalogVariantDimension401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminAddCatalogVariantDimension403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	res, err := h.catalogBase.CreateDimension(ctx, request.ModelId, input)
	if err != nil {
		return openapi.AdminAddCatalogVariantDimension400JSONResponse{}, nil
	}

	return openapi.AdminAddCatalogVariantDimension201JSONResponse(res), nil
}

// AdminUpdateCatalogVariantDimension handles PATCH /admin/catalog/product-models/{modelId}/variant-dimensions/{dimensionId}
func (h *Handler) AdminUpdateCatalogVariantDimension(ctx context.Context, request openapi.AdminUpdateCatalogVariantDimensionRequestObject) (openapi.AdminUpdateCatalogVariantDimensionResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUpdateCatalogVariantDimension401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUpdateCatalogVariantDimension403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	res, err := h.catalogBase.UpdateDimension(ctx, request.ModelId, request.DimensionId, input)
	if err != nil {
		return openapi.AdminUpdateCatalogVariantDimension400JSONResponse{}, nil
	}

	return openapi.AdminUpdateCatalogVariantDimension200JSONResponse(res), nil
}

// AdminAddCatalogVariantDimensionValue handles POST /admin/catalog/product-models/{modelId}/variant-dimensions/{dimensionId}/values
func (h *Handler) AdminAddCatalogVariantDimensionValue(ctx context.Context, request openapi.AdminAddCatalogVariantDimensionValueRequestObject) (openapi.AdminAddCatalogVariantDimensionValueResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminAddCatalogVariantDimensionValue401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminAddCatalogVariantDimensionValue403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	res, err := h.catalogBase.AddDimensionValue(ctx, request.ModelId, request.DimensionId, input)
	if err != nil {
		return openapi.AdminAddCatalogVariantDimensionValue400JSONResponse{}, nil
	}

	return openapi.AdminAddCatalogVariantDimensionValue201JSONResponse(res), nil
}

// AdminDeactivateCatalogVariantDimensionValue handles POST /admin/catalog/product-models/{modelId}/variant-dimensions/{dimensionId}/values/{valueId}/deactivate
func (h *Handler) AdminDeactivateCatalogVariantDimensionValue(ctx context.Context, request openapi.AdminDeactivateCatalogVariantDimensionValueRequestObject) (openapi.AdminDeactivateCatalogVariantDimensionValueResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminDeactivateCatalogVariantDimensionValue401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminDeactivateCatalogVariantDimensionValue403JSONResponse{}, nil
	}

	res, err := h.catalogBase.DeactivateDimensionValue(ctx, request.ModelId, request.DimensionId, request.ValueId)
	if err != nil {
		return openapi.AdminDeactivateCatalogVariantDimensionValue404JSONResponse{}, nil
	}

	return openapi.AdminDeactivateCatalogVariantDimensionValue200JSONResponse(res), nil
}
