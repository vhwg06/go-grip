package importer

import (
	usermodule "github.com/evrone/go-clean-template/internal/module/user"

	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	importermodule "github.com/evrone/go-clean-template/internal/module/importer"
)

// mapImporterError maps domain errors specific to Importer capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapImporterError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, usermodule.ErrForbidden) || errors.Is(err, usermodule.ErrUnauthorized) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "FORBIDDEN",
			Message: "Importer administrative access denied",
		})
		return http.StatusForbidden, resp
	}

	if errors.Is(err, importermodule.ErrInvalidInput) || errors.Is(err, usermodule.ErrInvalidInput) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "INVALID_INPUT",
			Message: "Invalid import payload structure",
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
