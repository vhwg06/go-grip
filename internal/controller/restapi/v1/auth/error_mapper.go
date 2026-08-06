package auth

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

// mapAuthError maps domain errors specific to the Auth capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
//
// Both entity.Err* and usermodule.Err* sentinels are checked because
// the auth usecase returns usermodule-package errors while callers may
// also propagate the shared entity errors. Checking only one set causes
// the other to fall through to 500.
func mapAuthError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, usermodule.ErrInvalidInput) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "USER_ALREADY_EXISTS",
			Message: "User with given username or email already exists",
		})
		return http.StatusConflict, resp
	}

	// ErrInvalidCredentials is returned by authUseCase (usermodule package) on
	// bad email/password; usermodule.ErrInvalidCredentials covers the shared sentinel.
	if errors.Is(err, usermodule.ErrInvalidCredentials) ||
		errors.Is(err, usermodule.ErrInvalidCredentials) ||
		errors.Is(err, usermodule.ErrNotFound) ||
		errors.Is(err, usermodule.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
		})
		return http.StatusUnauthorized, resp
	}

	if errors.Is(err, usermodule.ErrUnauthorized) || errors.Is(err, usermodule.ErrUnauthorized) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "UNAUTHORIZED",
			Message: "Authentication token missing or invalid",
		})
		return http.StatusUnauthorized, resp
	}

	if errors.Is(err, usermodule.ErrForbidden) || errors.Is(err, usermodule.ErrForbidden) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "FORBIDDEN",
			Message: "Access forbidden",
		})
		return http.StatusForbidden, resp
	}

	if errors.Is(err, usermodule.ErrInvalidInput) || errors.Is(err, usermodule.ErrInvalidInput) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "INVALID_INPUT",
			Message: "Invalid request payload",
		})
		return http.StatusBadRequest, resp
	}

	// Fallback to internal error
	resp := openapi.ErrorResponse{}
	resp.Error.FromErrorPayload(openapi.ErrorPayload{
		Code:    "INTERNAL_ERROR",
		Message: "An internal server error occurred",
	})
	return http.StatusInternalServerError, resp
}
