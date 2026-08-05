package v1

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
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

	if errors.Is(err, entity.ErrInvalidInput) {
		return http.StatusBadRequest, SharedErrorResponse{Error: "invalid request body"}
	}

	if errors.Is(err, entity.ErrUnauthorized) || errors.Is(err, entity.ErrInvalidCredentials) {
		return http.StatusUnauthorized, SharedErrorResponse{Error: "unauthorized"}
	}

	if errors.Is(err, entity.ErrForbidden) || errors.Is(err, entity.ErrCartBlocked) {
		return http.StatusForbidden, SharedErrorResponse{Error: "forbidden"}
	}

	if errors.Is(err, entity.ErrNotFound) {
		return http.StatusNotFound, SharedErrorResponse{Error: "not found"}
	}

	if errors.Is(err, entity.ErrRateLimited) {
		return http.StatusTooManyRequests, SharedErrorResponse{Error: "rate limit exceeded"}
	}

	return http.StatusInternalServerError, SharedErrorResponse{Error: "internal server error"}
}
