package notification

import (
	usermodule "github.com/evrone/go-clean-template/internal/module/user"

	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"
)

// mapNotificationError maps domain errors specific to Notification capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapNotificationError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, notificationmodule.ErrUnauthorized) || errors.Is(err, usermodule.ErrUnauthorized) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "UNAUTHORIZED",
			Message: "Authentication token required",
		})
		return http.StatusUnauthorized, resp
	}

	resp := openapi.ErrorResponse{}
	resp.Error.FromErrorPayload(openapi.ErrorPayload{
		Code:    "INTERNAL_ERROR",
		Message: "An internal server error occurred",
	})
	return http.StatusInternalServerError, resp
}
