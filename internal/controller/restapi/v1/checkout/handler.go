package checkout

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Checkout capability.
type Handler struct {
	checkoutUC ordermodule.CheckoutUseCase
	logger     logger.Interface
}

// NewHandler constructs a new Checkout vertical handler instance.
func NewHandler(checkoutUC ordermodule.CheckoutUseCase, l logger.Interface) *Handler {
	return &Handler{
		checkoutUC: checkoutUC,
		logger:     l,
	}
}

func getActor(ctx context.Context) usermodule.Actor {
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(usermodule.Actor); ok {
			return a
		}
	}
	return usermodule.Actor{}
}

// PreviewCheckout handles POST /checkout/preview
func (h *Handler) PreviewCheckout(ctx context.Context, request openapi.PreviewCheckoutRequestObject) (openapi.PreviewCheckoutResponseObject, error) {
	if request.Body == nil {
		return openapi.PreviewCheckout400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	productID := request.Body.ProductId
	quantity := 1
	if request.Body.Quantity > 0 {
		quantity = request.Body.Quantity
	}

	breakdown, err := h.checkoutUC.Preview(ctx, actor, productID, quantity)
	if err != nil {
		status, errResp := mapCheckoutError(err)
		switch status {
		case 400:
			return openapi.PreviewCheckout400JSONResponse{}, nil
		default:
			return openapi.PreviewCheckout500JSONResponse{
				InternalErrorResponseJSONResponse: openapi.InternalErrorResponseJSONResponse(errResp),
			}, nil
		}
	}

	res := toCheckoutPreviewResponse(breakdown)
	return openapi.PreviewCheckout200JSONResponse(res), nil
}

// CreateCheckoutOrder handles POST /checkout/orders
func (h *Handler) CreateCheckoutOrder(ctx context.Context, request openapi.CreateCheckoutOrderRequestObject) (openapi.CreateCheckoutOrderResponseObject, error) {
	if request.Body == nil {
		return openapi.CreateCheckoutOrder400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	productID := request.Body.ProductId
	quantity := 1
	if request.Body.Quantity > 0 {
		quantity = request.Body.Quantity
	}
	email := request.Body.Email

	orderEntity, err := h.checkoutUC.CreateOrder(ctx, actor, productID, quantity, email)
	if err != nil {
		status, errResp := mapCheckoutError(err)
		switch status {
		case 400:
			return openapi.CreateCheckoutOrder400JSONResponse{}, nil
		case 401:
			return openapi.CreateCheckoutOrder401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.CreateCheckoutOrder500JSONResponse{}, nil
		}
	}

	res := toCheckoutOrderResponse(orderEntity)
	return openapi.CreateCheckoutOrder201JSONResponse(res), nil
}

// GetPaymentParams handles GET /checkout/orders/{orderId}/payment-params
func (h *Handler) GetPaymentParams(ctx context.Context, request openapi.GetPaymentParamsRequestObject) (openapi.GetPaymentParamsResponseObject, error) {
	actor := getActor(ctx)
	params, err := h.checkoutUC.PaymentParams(ctx, actor, request.OrderId)
	if err != nil {
		status, errResp := mapCheckoutError(err)
		switch status {
		case 401:
			return openapi.GetPaymentParams401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.GetPaymentParams404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetPaymentParams500JSONResponse{}, nil
		}
	}

	res := toPaymentParamsResponse(params)
	return openapi.GetPaymentParams200JSONResponse(res), nil
}

// PaymentNotify handles POST /checkout/notify
func (h *Handler) PaymentNotify(ctx context.Context, request openapi.PaymentNotifyRequestObject) (openapi.PaymentNotifyResponseObject, error) {
	_ = request
	payload := map[string]string{}

	if err := h.checkoutUC.PaymentNotify(ctx, payload); err != nil {
		status, _ := mapCheckoutError(err)
		switch status {
		case 400:
			return openapi.PaymentNotify400JSONResponse{}, nil
		default:
			return openapi.PaymentNotify500JSONResponse{}, nil
		}
	}

	return openapi.PaymentNotify200Response{}, nil
}

// GetPaymentStatus handles GET /checkout/orders/{orderId}/status
func (h *Handler) GetPaymentStatus(ctx context.Context, request openapi.GetPaymentStatusRequestObject) (openapi.GetPaymentStatusResponseObject, error) {
	orderEntity, err := h.checkoutUC.PaymentStatus(ctx, request.OrderId)
	if err != nil {
		status, errResp := mapCheckoutError(err)
		switch status {
		case 404:
			return openapi.GetPaymentStatus404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetPaymentStatus500JSONResponse{}, nil
		}
	}

	res := toCheckoutOrderResponse(orderEntity)
	return openapi.GetPaymentStatus200JSONResponse(res), nil
}
