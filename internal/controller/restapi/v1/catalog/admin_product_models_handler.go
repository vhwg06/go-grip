package catalog

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
)

// AdminListCatalogProductModels handles GET /admin/catalog/product-models
func (h *Handler) AdminListCatalogProductModels(ctx context.Context, _ openapi.AdminListCatalogProductModelsRequestObject) (openapi.AdminListCatalogProductModelsResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminListCatalogProductModels401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminListCatalogProductModels403JSONResponse{}, nil
	}

	_, err := h.catalogBase.ListModels(ctx)
	if err != nil {
		return openapi.AdminListCatalogProductModels500JSONResponse{}, nil
	}

	return openapi.AdminListCatalogProductModels200Response{}, nil
}

// AdminCreateCatalogProductModel handles POST /admin/catalog/product-models
func (h *Handler) AdminCreateCatalogProductModel(ctx context.Context, request openapi.AdminCreateCatalogProductModelRequestObject) (openapi.AdminCreateCatalogProductModelResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminCreateCatalogProductModel401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminCreateCatalogProductModel403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	_, err := h.catalogBase.CreateModel(ctx, input)
	if err != nil {
		return openapi.AdminCreateCatalogProductModel400JSONResponse{}, nil
	}

	return openapi.AdminCreateCatalogProductModel201Response{}, nil
}

// AdminGetCatalogProductModel handles GET /admin/catalog/product-models/{modelId}
func (h *Handler) AdminGetCatalogProductModel(ctx context.Context, request openapi.AdminGetCatalogProductModelRequestObject) (openapi.AdminGetCatalogProductModelResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminGetCatalogProductModel401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminGetCatalogProductModel403JSONResponse{}, nil
	}

	_, err := h.catalogBase.GetModel(ctx, request.ModelId)
	if err != nil {
		return openapi.AdminGetCatalogProductModel404JSONResponse{}, nil
	}

	return openapi.AdminGetCatalogProductModel200Response{}, nil
}

// AdminUpdateCatalogProductModel handles PATCH /admin/catalog/product-models/{modelId}
func (h *Handler) AdminUpdateCatalogProductModel(ctx context.Context, request openapi.AdminUpdateCatalogProductModelRequestObject) (openapi.AdminUpdateCatalogProductModelResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUpdateCatalogProductModel401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUpdateCatalogProductModel403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	_, err := h.catalogBase.UpdateModel(ctx, request.ModelId, input)
	if err != nil {
		return openapi.AdminUpdateCatalogProductModel400JSONResponse{}, nil
	}

	return openapi.AdminUpdateCatalogProductModel200Response{}, nil
}

// AdminDeleteCatalogProductModel handles DELETE /admin/catalog/product-models/{modelId}
func (h *Handler) AdminDeleteCatalogProductModel(ctx context.Context, request openapi.AdminDeleteCatalogProductModelRequestObject) (openapi.AdminDeleteCatalogProductModelResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminDeleteCatalogProductModel401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminDeleteCatalogProductModel403JSONResponse{}, nil
	}

	_, err := h.catalogBase.DeleteModel(ctx, request.ModelId)
	if err != nil {
		return openapi.AdminDeleteCatalogProductModel404JSONResponse{}, nil
	}

	return openapi.AdminDeleteCatalogProductModel204Response{}, nil
}

// AdminPublishCatalogProductModel handles POST /admin/catalog/product-models/{modelId}/publish
func (h *Handler) AdminPublishCatalogProductModel(ctx context.Context, request openapi.AdminPublishCatalogProductModelRequestObject) (openapi.AdminPublishCatalogProductModelResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminPublishCatalogProductModel401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminPublishCatalogProductModel403JSONResponse{}, nil
	}

	_, err := h.catalogBase.PublishModel(ctx, request.ModelId)
	if err != nil {
		return openapi.AdminPublishCatalogProductModel400JSONResponse{}, nil
	}

	return openapi.AdminPublishCatalogProductModel200Response{}, nil
}

// AdminUnpublishCatalogProductModel handles POST /admin/catalog/product-models/{modelId}/unpublish
func (h *Handler) AdminUnpublishCatalogProductModel(ctx context.Context, request openapi.AdminUnpublishCatalogProductModelRequestObject) (openapi.AdminUnpublishCatalogProductModelResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUnpublishCatalogProductModel401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUnpublishCatalogProductModel403JSONResponse{}, nil
	}

	_, err := h.catalogBase.UnpublishModel(ctx, request.ModelId)
	if err != nil {
		return openapi.AdminUnpublishCatalogProductModel400JSONResponse{}, nil
	}

	return openapi.AdminUnpublishCatalogProductModel200Response{}, nil
}

// AdminDiscontinueCatalogProductModel handles POST /admin/catalog/product-models/{modelId}/discontinue
func (h *Handler) AdminDiscontinueCatalogProductModel(ctx context.Context, request openapi.AdminDiscontinueCatalogProductModelRequestObject) (openapi.AdminDiscontinueCatalogProductModelResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminDiscontinueCatalogProductModel401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminDiscontinueCatalogProductModel403JSONResponse{}, nil
	}

	_, err := h.catalogBase.DiscontinueModel(ctx, request.ModelId)
	if err != nil {
		return openapi.AdminDiscontinueCatalogProductModel400JSONResponse{}, nil
	}

	return openapi.AdminDiscontinueCatalogProductModel200Response{}, nil
}

// AdminUpdateCatalogProductModelMedia handles PUT /admin/catalog/product-models/{modelId}/media
func (h *Handler) AdminUpdateCatalogProductModelMedia(ctx context.Context, request openapi.AdminUpdateCatalogProductModelMediaRequestObject) (openapi.AdminUpdateCatalogProductModelMediaResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUpdateCatalogProductModelMedia401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUpdateCatalogProductModelMedia403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	_, err := h.catalogBase.ReplaceMedia(ctx, request.ModelId, input)
	if err != nil {
		return openapi.AdminUpdateCatalogProductModelMedia400JSONResponse{}, nil
	}

	return openapi.AdminUpdateCatalogProductModelMedia200Response{}, nil
}

// AdminListCatalogModelVariants handles GET /admin/catalog/product-models/{modelId}/variants
func (h *Handler) AdminListCatalogModelVariants(ctx context.Context, request openapi.AdminListCatalogModelVariantsRequestObject) (openapi.AdminListCatalogModelVariantsResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminListCatalogModelVariants401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminListCatalogModelVariants403JSONResponse{}, nil
	}

	_, err := h.catalogBase.ListVariants(ctx, request.ModelId)
	if err != nil {
		return openapi.AdminListCatalogModelVariants500JSONResponse{}, nil
	}

	return openapi.AdminListCatalogModelVariants200Response{}, nil
}
