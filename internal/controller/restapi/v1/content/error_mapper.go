package content

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
)

// mapContentError maps domain errors specific to Content & Homepage capabilities
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapContentError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, contentmodule.ErrNotFound) || errors.Is(err, entity.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "PAGE_NOT_FOUND",
			Message: "Static page or content resource not found",
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
