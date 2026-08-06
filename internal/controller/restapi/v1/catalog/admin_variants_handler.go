package catalog

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
)

// AdminCreateCatalogModelVariant handles POST /admin/catalog/product-models/{modelId}/variants
func (h *Handler) AdminCreateCatalogModelVariant(ctx context.Context, request openapi.AdminCreateCatalogModelVariantRequestObject) (openapi.AdminCreateCatalogModelVariantResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminCreateCatalogModelVariant401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminCreateCatalogModelVariant403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	res, err := h.catalogBase.CreateVariant(ctx, request.ModelId, input)
	if err != nil {
		if h.logger != nil {
			h.logger.Error("catalog base create Variant failed", "error", err)
		}
		return openapi.AdminCreateCatalogModelVariant400JSONResponse{}, nil
	}

	return openapi.AdminCreateCatalogModelVariant201JSONResponse(res), nil
}

// AdminGetCatalogVariant handles GET /admin/catalog/variants/{variantId}
func (h *Handler) AdminGetCatalogVariant(ctx context.Context, request openapi.AdminGetCatalogVariantRequestObject) (openapi.AdminGetCatalogVariantResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminGetCatalogVariant401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminGetCatalogVariant403JSONResponse{}, nil
	}

	res, err := h.catalogBase.GetVariant(ctx, request.VariantId)
	if err != nil {
		return openapi.AdminGetCatalogVariant404JSONResponse{}, nil
	}

	return openapi.AdminGetCatalogVariant200JSONResponse(res), nil
}

// AdminUpdateCatalogVariant handles PATCH /admin/catalog/variants/{variantId}
func (h *Handler) AdminUpdateCatalogVariant(ctx context.Context, request openapi.AdminUpdateCatalogVariantRequestObject) (openapi.AdminUpdateCatalogVariantResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUpdateCatalogVariant401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUpdateCatalogVariant403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	res, err := h.catalogBase.UpdateVariant(ctx, request.VariantId, input)
	if err != nil {
		return openapi.AdminUpdateCatalogVariant400JSONResponse{}, nil
	}

	return openapi.AdminUpdateCatalogVariant200JSONResponse(res), nil
}

// AdminActivateCatalogVariant handles POST /admin/catalog/variants/{variantId}/activate
func (h *Handler) AdminActivateCatalogVariant(ctx context.Context, request openapi.AdminActivateCatalogVariantRequestObject) (openapi.AdminActivateCatalogVariantResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminActivateCatalogVariant401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminActivateCatalogVariant403JSONResponse{}, nil
	}

	res, err := h.catalogBase.ActivateVariant(ctx, request.VariantId)
	if err != nil {
		return openapi.AdminActivateCatalogVariant404JSONResponse{}, nil
	}

	return openapi.AdminActivateCatalogVariant200JSONResponse(res), nil
}

// AdminInactivateCatalogVariant handles POST /admin/catalog/variants/{variantId}/inactivate
func (h *Handler) AdminInactivateCatalogVariant(ctx context.Context, request openapi.AdminInactivateCatalogVariantRequestObject) (openapi.AdminInactivateCatalogVariantResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminInactivateCatalogVariant401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminInactivateCatalogVariant403JSONResponse{}, nil
	}

	res, err := h.catalogBase.InactivateVariant(ctx, request.VariantId)
	if err != nil {
		return openapi.AdminInactivateCatalogVariant404JSONResponse{}, nil
	}

	return openapi.AdminInactivateCatalogVariant200JSONResponse(res), nil
}

// AdminBulkUpdateVariantPrices handles POST /admin/catalog/variants/prices:bulk
func (h *Handler) AdminBulkUpdateVariantPrices(ctx context.Context, request openapi.AdminBulkUpdateVariantPricesRequestObject) (openapi.AdminBulkUpdateVariantPricesResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminBulkUpdateVariantPrices401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminBulkUpdateVariantPrices403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	res, err := h.catalogBase.BulkSetPrice(ctx, input)
	if err != nil {
		return openapi.AdminBulkUpdateVariantPrices400JSONResponse{}, nil
	}

	return openapi.AdminBulkUpdateVariantPrices200JSONResponse(map[string]interface{}{"result": res}), nil
}

// GetCatalogSettings handles GET /catalog/settings
func (h *Handler) GetCatalogSettings(ctx context.Context, _ openapi.GetCatalogSettingsRequestObject) (openapi.GetCatalogSettingsResponseObject, error) {
	settings, err := h.catalogUC.ListPublicSettings(ctx)
	if err != nil {
		return openapi.GetCatalogSettings500JSONResponse{}, nil
	}
	resp := map[string]any{
		"site_name":       "Grip Store",
		"currency":        "VND",
		"shopName":        "Grip Store",
		"shopDescription": "",
		"themeColor":      "",
	}
	for _, setting := range settings {
		switch setting.Key {
		case "brand.shopName":
			resp["shopName"] = setting.Value
		case "brand.shopDescription":
			resp["shopDescription"] = setting.Value
		case "brand.themeColor":
			resp["themeColor"] = setting.Value
		}
	}
	return openapi.GetCatalogSettings200JSONResponse(resp), nil
}
