package checkout

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
)

// CancelCheckoutOrder handles POST /checkout/orders/{orderId}/cancel
func (h *Handler) CancelCheckoutOrder(ctx context.Context, request openapi.CancelCheckoutOrderRequestObject) (openapi.CancelCheckoutOrderResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.CancelCheckoutOrder401JSONResponse{}, nil
	}

	if err := h.checkoutUC.Cancel(ctx, actor, request.OrderId); err != nil {
		status, errResp := mapCheckoutError(err)
		switch status {
		case 400:
			return openapi.CancelCheckoutOrder400JSONResponse{BadRequestResponseJSONResponse: openapi.BadRequestResponseJSONResponse(errResp)}, nil
		case 404:
			return openapi.CancelCheckoutOrder404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		default:
			return openapi.CancelCheckoutOrder400JSONResponse{BadRequestResponseJSONResponse: openapi.BadRequestResponseJSONResponse(errResp)}, nil
		}
	}

	return openapi.CancelCheckoutOrder200Response{}, nil
}

// GetCheckoutOrderStatus handles GET /checkout/orders/{orderId}/status
func (h *Handler) GetCheckoutOrderStatus(ctx context.Context, request openapi.GetCheckoutOrderStatusRequestObject) (openapi.GetCheckoutOrderStatusResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.GetCheckoutOrderStatus401JSONResponse{}, nil
	}

	order, err := h.checkoutUC.PaymentStatus(ctx, request.OrderId)
	if err != nil {
		status, errResp := mapCheckoutError(err)
		switch status {
		case 404:
			return openapi.GetCheckoutOrderStatus404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		default:
			return openapi.GetCheckoutOrderStatus500JSONResponse{}, nil
		}
	}

	resp := toCheckoutOrderResponse(order)
	return openapi.GetCheckoutOrderStatus200JSONResponse(resp), nil
}
