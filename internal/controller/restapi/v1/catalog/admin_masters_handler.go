package catalog

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
)

// AdminListCatalogMasters handles GET /admin/catalog/masters/{masterKind}
func (h *Handler) AdminListCatalogMasters(ctx context.Context, request openapi.AdminListCatalogMastersRequestObject) (openapi.AdminListCatalogMastersResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminListCatalogMasters401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminListCatalogMasters403JSONResponse{}, nil
	}

	_, err := h.catalogBase.ListMasters(ctx, request.MasterKind)
	if err != nil {
		return openapi.AdminListCatalogMasters500JSONResponse{}, nil
	}

	return openapi.AdminListCatalogMasters200Response{}, nil
}

// AdminCreateCatalogMaster handles POST /admin/catalog/masters/{masterKind}
func (h *Handler) AdminCreateCatalogMaster(ctx context.Context, request openapi.AdminCreateCatalogMasterRequestObject) (openapi.AdminCreateCatalogMasterResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminCreateCatalogMaster401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminCreateCatalogMaster403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	_, err := h.catalogBase.CreateMaster(ctx, request.MasterKind, input)
	if err != nil {
		return openapi.AdminCreateCatalogMaster400JSONResponse{}, nil
	}

	return openapi.AdminCreateCatalogMaster201Response{}, nil
}

// AdminUpdateCatalogMaster handles PATCH /admin/catalog/masters/{masterKind}/{masterId}
func (h *Handler) AdminUpdateCatalogMaster(ctx context.Context, request openapi.AdminUpdateCatalogMasterRequestObject) (openapi.AdminUpdateCatalogMasterResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUpdateCatalogMaster401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUpdateCatalogMaster403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	_, err := h.catalogBase.UpdateMaster(ctx, request.MasterKind, request.MasterId, input)
	if err != nil {
		return openapi.AdminUpdateCatalogMaster400JSONResponse{}, nil
	}

	return openapi.AdminUpdateCatalogMaster200Response{}, nil
}

// AdminDeactivateCatalogMaster handles POST /admin/catalog/masters/{masterKind}/{masterId}/deactivate
func (h *Handler) AdminDeactivateCatalogMaster(ctx context.Context, request openapi.AdminDeactivateCatalogMasterRequestObject) (openapi.AdminDeactivateCatalogMasterResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminDeactivateCatalogMaster401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminDeactivateCatalogMaster403JSONResponse{}, nil
	}

	_, err := h.catalogBase.DeactivateMaster(ctx, request.MasterKind, request.MasterId)
	if err != nil {
		return openapi.AdminDeactivateCatalogMaster404JSONResponse{}, nil
	}

	return openapi.AdminDeactivateCatalogMaster200Response{}, nil
}
