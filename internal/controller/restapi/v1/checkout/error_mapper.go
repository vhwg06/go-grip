package checkout

import (
	usermodule "github.com/evrone/go-clean-template/internal/module/user"

	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
)

// mapCheckoutError maps domain errors specific to Checkout capability
// to standard HTTP status codes and openapi.ErrorResponse DTOs.
func mapCheckoutError(err error) (int, openapi.ErrorResponse) {
	if err == nil {
		return http.StatusOK, openapi.ErrorResponse{}
	}

	if errors.Is(err, ordermodule.ErrNotFound) || errors.Is(err, usermodule.ErrNotFound) || errors.Is(err, ordermodule.ErrNotFound) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "ORDER_NOT_FOUND",
			Message: "Order or resource not found",
		})
		return http.StatusNotFound, resp
	}

	if errors.Is(err, ordermodule.ErrOutOfStock) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "OUT_OF_STOCK",
			Message: "The requested product quantity is not available",
		})
		return http.StatusBadRequest, resp
	}

	if errors.Is(err, ordermodule.ErrPaymentInvalidSign) || errors.Is(err, ordermodule.ErrInvalidInput) || errors.Is(err, ordermodule.ErrPaymentInvalidSign) {
		resp := openapi.ErrorResponse{}
		resp.Error.FromErrorPayload(openapi.ErrorPayload{
			Code:    "PAYMENT_FAILED",
			Message: "Payment processing failed or signature invalid",
		})
		return http.StatusBadRequest, resp
	}

	if errors.Is(err, ordermodule.ErrInvalidInput) || errors.Is(err, usermodule.ErrInvalidInput) {
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
