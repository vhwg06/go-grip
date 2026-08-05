package checkout

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
)

// mapCheckoutError maps domain errors specific to Checkout capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapCheckoutError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, ordermodule.ErrNotFound) || errors.Is(err, entity.ErrNotFound) || errors.Is(err, entity.ErrOrderNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "ORDER_NOT_FOUND",
			Message: "Order or resource not found",
		})
		return http.StatusNotFound, resp
	}

	if errors.Is(err, ordermodule.ErrPaymentInvalidSign) || errors.Is(err, entity.ErrPaymentFailed) || errors.Is(err, entity.ErrPaymentInvalidSign) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "PAYMENT_FAILED",
			Message: "Payment processing failed or signature invalid",
		})
		return http.StatusBadRequest, resp
	}

	if errors.Is(err, ordermodule.ErrInvalidInput) || errors.Is(err, entity.ErrInvalidInput) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "INVALID_INPUT",
			Message: "Invalid checkout request payload",
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
