package lead

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	leadmodule "github.com/evrone/go-clean-template/internal/module/lead"
)

// mapLeadError maps domain errors specific to Lead capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapLeadError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, leadmodule.ErrNotFound) || errors.Is(err, entity.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "LEAD_NOT_FOUND",
			Message: "Lead submission not found",
		})
		return http.StatusNotFound, resp
	}

	if errors.Is(err, entity.ErrInvalidInput) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "INVALID_INPUT",
			Message: "Invalid lead submission data",
		})
		return http.StatusBadRequest, resp
	}

	if errors.Is(err, entity.ErrForbidden) || errors.Is(err, entity.ErrUnauthorized) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "FORBIDDEN",
			Message: "Access to lead resource denied",
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
