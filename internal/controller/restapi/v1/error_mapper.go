package v1

import (
	"errors"
	"net/http"

	cartmodule "github.com/evrone/go-clean-template/internal/module/cart"
	usermodule "github.com/evrone/go-clean-template/internal/module/user" 
)

// SharedErrorResponse represents the standard HTTP error DTO response.
type SharedErrorResponse struct {
	Error string `json:"error"`
}

// MapSharedError maps common transport and foundation errors (400, 401, 403, 404, 429, 500)
// to HTTP status code and SharedErrorResponse payload.
// Does NOT import capability-specific domain errors.
func MapSharedError(err error) (int, SharedErrorResponse) {
	if err == nil {
		return http.StatusOK, SharedErrorResponse{}
	}

	if errors.Is(err, usermodule.ErrInvalidInput) {
		return http.StatusBadRequest, SharedErrorResponse{Error: "invalid request body"}
	}

	if errors.Is(err, usermodule.ErrUnauthorized) || errors.Is(err, usermodule.ErrInvalidCredentials) {
		return http.StatusUnauthorized, SharedErrorResponse{Error: "unauthorized"}
	}

	if errors.Is(err, usermodule.ErrForbidden) || errors.Is(err, cartmodule.ErrCartBlocked) {
		return http.StatusForbidden, SharedErrorResponse{Error: "forbidden"}
	}

	if errors.Is(err, usermodule.ErrNotFound) {
		return http.StatusNotFound, SharedErrorResponse{Error: "not found"}
	}

	if errors.Is(err, usermodule.ErrForbidden) {
		return http.StatusTooManyRequests, SharedErrorResponse{Error: "rate limit exceeded"}
	}

	return http.StatusInternalServerError, SharedErrorResponse{Error: "internal server error"}
}
