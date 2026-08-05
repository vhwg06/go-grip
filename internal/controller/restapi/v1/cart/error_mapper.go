package cart

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// mapCartError maps domain errors specific to Cart capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapCartError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, entity.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "ITEM_NOT_FOUND",
			Message: "Cart item or product not found",
		})
		return http.StatusNotFound, resp
	}

	if errors.Is(err, entity.ErrCartBlocked) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "CART_BLOCKED",
			Message: "Cart is currently blocked for checkout",
		})
		return http.StatusBadRequest, resp
	}

	if errors.Is(err, entity.ErrOutOfStock) || errors.Is(err, entity.ErrProductUnavailable) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "OUT_OF_STOCK",
			Message: "Product out of stock or unavailable",
		})
		return http.StatusBadRequest, resp
	}

	if errors.Is(err, entity.ErrInvalidInput) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "INVALID_INPUT",
			Message: "Invalid cart request payload",
		})
		return http.StatusBadRequest, resp
	}

	resp := openapi.ErrorResponse{}
	resp.Error.FromErrorPayload(openapi.ErrorPayload{
		Code:    "INTERNAL_ERROR",
		Message: "An internal server error occurred",
	})
	return http.StatusInternalServerError, resp
}
