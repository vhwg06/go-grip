package admin

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

// mapAdminError maps domain errors specific to Admin capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
//
// Both entity.Err* and usermodule.Err* sentinels are checked because
// ensureAdmin returns usermodule-package errors while other callers may
// propagate the shared entity sentinels. Checking only one set causes
// the other to fall through to 500.
// ErrUnauthorized (unauthenticated caller) → 401; ErrForbidden (authenticated
// but non-admin caller) → 403.
func mapAdminError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	// Unauthenticated: no valid actor identity → 401.
	if errors.Is(err, entity.ErrUnauthorized) || errors.Is(err, usermodule.ErrUnauthorized) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "UNAUTHORIZED",
			Message: "Authentication required",
		})
		return http.StatusUnauthorized, resp
	}

	// Authenticated but not an admin → 403.
	if errors.Is(err, entity.ErrForbidden) || errors.Is(err, usermodule.ErrForbidden) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "FORBIDDEN",
			Message: "Administrative access denied",
		})
		return http.StatusForbidden, resp
	}

	if errors.Is(err, entity.ErrNotFound) || errors.Is(err, usermodule.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "NOT_FOUND",
			Message: "Admin resource not found",
		})
		return http.StatusNotFound, resp
	}

	resp := openapi.ErrorResponse{}
	resp.Error.FromErrorPayload(openapi.ErrorPayload{
		Code:    "INTERNAL_ERROR",
		Message: "An internal server error occurred",
	})
	return http.StatusInternalServerError, resp
}
