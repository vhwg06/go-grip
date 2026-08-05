package wishlist

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
)

// mapWishlistError maps domain errors specific to Wishlist capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapWishlistError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, wishlistmodule.ErrNotFound) || errors.Is(err, entity.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "WISHLIST_ITEM_NOT_FOUND",
			Message: "Wishlist item not found",
		})
		return http.StatusNotFound, resp
	}

	if errors.Is(err, wishlistmodule.ErrUnauthorized) || errors.Is(err, entity.ErrUnauthorized) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "UNAUTHORIZED",
			Message: "Authentication required",
		})
		return http.StatusUnauthorized, resp
	}

	resp := openapi.ErrorResponse{}
	resp.Error.FromErrorPayload(openapi.ErrorPayload{
		Code:    "INTERNAL_ERROR",
		Message: "An internal server error occurred",
	})
	return http.StatusInternalServerError, resp
}
