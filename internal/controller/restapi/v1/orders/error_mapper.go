package orders

import (
	usermodule "github.com/evrone/go-clean-template/internal/module/user"

	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
)

// mapOrdersError maps domain errors specific to Orders capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapOrdersError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, ordermodule.ErrNotFound) || errors.Is(err, ordermodule.ErrNotFound) || errors.Is(err, usermodule.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "ORDER_NOT_FOUND",
			Message: "Order or resource not found",
		})
		return http.StatusNotFound, resp
	}

	if errors.Is(err, ordermodule.ErrRefundNotAllowed) || errors.Is(err, ordermodule.ErrRefundNotAllowed) || errors.Is(err, ordermodule.ErrInvalidInput) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "REFUND_NOT_ALLOWED",
			Message: "Order refund is not permitted in current status",
		})
		return http.StatusUnprocessableEntity, resp
	}

	if errors.Is(err, ordermodule.ErrInvalidInput) || errors.Is(err, usermodule.ErrInvalidInput) {
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
