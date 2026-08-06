package catalog

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

func getActor(ctx context.Context) usermodule.Actor {
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(usermodule.Actor); ok {
			return a
		}
	}
	return usermodule.Actor{}
}

// AdminListCatalogCategories handles GET /admin/catalog/categories
func (h *Handler) AdminListCatalogCategories(ctx context.Context, _ openapi.AdminListCatalogCategoriesRequestObject) (openapi.AdminListCatalogCategoriesResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminListCatalogCategories401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminListCatalogCategories403JSONResponse{}, nil
	}

	_, err := h.catalogBase.ListCategories(ctx)
	if err != nil {
		return openapi.AdminListCatalogCategories500JSONResponse{}, nil
	}

	return openapi.AdminListCatalogCategories200Response{}, nil
}

// AdminCreateCatalogCategory handles POST /admin/catalog/categories
func (h *Handler) AdminCreateCatalogCategory(ctx context.Context, request openapi.AdminCreateCatalogCategoryRequestObject) (openapi.AdminCreateCatalogCategoryResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminCreateCatalogCategory401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminCreateCatalogCategory403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	_, err := h.catalogBase.CreateCategory(ctx, input)
	if err != nil {
		return openapi.AdminCreateCatalogCategory400JSONResponse{}, nil
	}

	return openapi.AdminCreateCatalogCategory201Response{}, nil
}

// AdminUpdateCatalogCategory handles PATCH /admin/catalog/categories/{categoryId}
func (h *Handler) AdminUpdateCatalogCategory(ctx context.Context, request openapi.AdminUpdateCatalogCategoryRequestObject) (openapi.AdminUpdateCatalogCategoryResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUpdateCatalogCategory401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUpdateCatalogCategory403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	_, err := h.catalogBase.UpdateCategory(ctx, request.CategoryId, input)
	if err != nil {
		return openapi.AdminUpdateCatalogCategory400JSONResponse{}, nil
	}

	return openapi.AdminUpdateCatalogCategory200Response{}, nil
}

// AdminDeleteCatalogCategory handles DELETE /admin/catalog/categories/{categoryId}
func (h *Handler) AdminDeleteCatalogCategory(ctx context.Context, request openapi.AdminDeleteCatalogCategoryRequestObject) (openapi.AdminDeleteCatalogCategoryResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminDeleteCatalogCategory401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminDeleteCatalogCategory403JSONResponse{}, nil
	}

	_, err := h.catalogBase.DeleteCategory(ctx, request.CategoryId)
	if err != nil {
		return openapi.AdminDeleteCatalogCategory404JSONResponse{}, nil
	}

	return openapi.AdminDeleteCatalogCategory204Response{}, nil
}

// AdminDeactivateCatalogCategory handles POST /admin/catalog/categories/{categoryId}/deactivate
func (h *Handler) AdminDeactivateCatalogCategory(ctx context.Context, request openapi.AdminDeactivateCatalogCategoryRequestObject) (openapi.AdminDeactivateCatalogCategoryResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminDeactivateCatalogCategory401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminDeactivateCatalogCategory403JSONResponse{}, nil
	}

	_, err := h.catalogBase.DeactivateCategory(ctx, request.CategoryId)
	if err != nil {
		return openapi.AdminDeactivateCatalogCategory404JSONResponse{}, nil
	}

	return openapi.AdminDeactivateCatalogCategory200Response{}, nil
}
