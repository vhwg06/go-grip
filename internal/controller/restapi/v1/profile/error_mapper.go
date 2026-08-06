package profile

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

// mapProfileError maps domain errors specific to Profile capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
//
// Both entity.Err* and usermodule.Err* sentinels are checked because
// the profile usecase returns usermodule-package errors while shared
// infrastructure may propagate entity sentinels.
func mapProfileError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, usermodule.ErrNotFound) || errors.Is(err, usermodule.ErrNotFound) ||
		errors.Is(err, usermodule.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "USER_NOT_FOUND",
			Message: "User profile not found",
		})
		return http.StatusNotFound, resp
	}

	if errors.Is(err, usermodule.ErrUnauthorized) || errors.Is(err, usermodule.ErrUnauthorized) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "UNAUTHORIZED",
			Message: "Authentication required",
		})
		return http.StatusUnauthorized, resp
	}

	if errors.Is(err, usermodule.ErrInvalidInput) || errors.Is(err, usermodule.ErrInvalidInput) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "INVALID_INPUT",
			Message: "Invalid profile update payload",
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
