package task

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// mapTaskError maps domain errors specific to the Task capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapTaskError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, entity.ErrTaskNotFound) || errors.Is(err, entity.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "TASK_NOT_FOUND",
			Message: "Task resource not found",
		})
		return http.StatusNotFound, resp
	}

	if errors.Is(err, entity.ErrTaskForbidden) || errors.Is(err, entity.ErrForbidden) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "TASK_FORBIDDEN",
			Message: "Task does not belong to user",
		})
		return http.StatusForbidden, resp
	}

	if errors.Is(err, entity.ErrInvalidTransition) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "INVALID_TRANSITION",
			Message: "Invalid status transition for task",
		})
		return http.StatusUnprocessableEntity, resp
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
