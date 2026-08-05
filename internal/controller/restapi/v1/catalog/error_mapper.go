package catalog

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
)

// mapCatalogError maps domain errors specific to Catalog capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapCatalogError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, catalogmodule.ErrNotFound) || errors.Is(err, entity.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "PRODUCT_NOT_FOUND",
			Message: "Product not found",
		})
		return http.StatusNotFound, resp
	}

	if errors.Is(err, entity.ErrDuplicateSKU) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "DUPLICATE_SKU",
			Message: "Product SKU already exists",
		})
		return http.StatusConflict, resp
	}

	if errors.Is(err, catalogmodule.ErrInvalidInput) || errors.Is(err, entity.ErrInvalidInput) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "INVALID_INPUT",
			Message: "Invalid product payload",
		})
		return http.StatusBadRequest, resp
	}

	if errors.Is(err, entity.ErrForbidden) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "FORBIDDEN",
			Message: "Insufficient permissions for catalog management",
		})
		return http.StatusForbidden, resp
	}

	resp := openapi.ErrorResponse{}
	resp.Error.FromErrorPayload(openapi.ErrorPayload{
		Code:    "INTERNAL_ERROR",
		Message: "An internal server error occurred",
	})
	return http.StatusInternalServerError, resp
}
