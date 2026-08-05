package user

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// mapUserError maps domain errors specific to the User capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapUserError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, entity.ErrUserNotFound) || errors.Is(err, entity.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "USER_NOT_FOUND",
			Message: "User resource not found",
		})
		return http.StatusNotFound, resp
	}

	if errors.Is(err, entity.ErrUserAlreadyExists) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "USER_ALREADY_EXISTS",
			Message: "User with given username or email already exists",
		})
		return http.StatusConflict, resp
	}

	if errors.Is(err, entity.ErrForbidden) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "FORBIDDEN",
			Message: "Access forbidden",
		})
		return http.StatusForbidden, resp
	}

	if errors.Is(err, entity.ErrInvalidInput) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "INVALID_INPUT",
			Message: "Invalid request payload",
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
