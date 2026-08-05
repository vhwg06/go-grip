package checkout

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Checkout capability.
type Handler struct {
	checkoutUC usecase.Checkout
	logger     logger.Interface
}

// NewHandler constructs a new Checkout vertical handler instance.
func NewHandler(checkoutUC usecase.Checkout, l logger.Interface) *Handler {
	return &Handler{
		checkoutUC: checkoutUC,
		logger:     l,
	}
}

func getActor(ctx context.Context) entity.Actor {
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(entity.Actor); ok {
			return a
		}
	}
	return entity.Actor{}
}

// PreviewCheckout handles POST /checkout/preview
func (h *Handler) PreviewCheckout(ctx context.Context, request openapi.PreviewCheckoutRequestObject) (openapi.PreviewCheckoutResponseObject, error) {
	if request.Body == nil {
		return openapi.PreviewCheckout400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	breakdown, err := h.checkoutUC.Preview(ctx, actor, request.Body.ProductId, request.Body.Quantity)
	if err != nil {
		status, errResp := mapCheckoutError(err)
		switch status {
		case 400:
			return openapi.PreviewCheckout400JSONResponse{}, nil
		case 401:
			return openapi.PreviewCheckout401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.PreviewCheckout500JSONResponse{}, nil
		}
	}

	previewDTO := toCheckoutPreviewResponse(breakdown)
	return openapi.PreviewCheckout200JSONResponse(previewDTO), nil
}

// CreateCheckoutOrder handles POST /checkout/orders
func (h *Handler) CreateCheckoutOrder(ctx context.Context, request openapi.CreateCheckoutOrderRequestObject) (openapi.CreateCheckoutOrderResponseObject, error) {
	if request.Body == nil {
		return openapi.CreateCheckoutOrder400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	order, err := h.checkoutUC.CreateOrder(ctx, actor, request.Body.ProductId, request.Body.Quantity, request.Body.Email)
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

	orderDTO := toCheckoutOrderResponse(order)
	return openapi.CreateCheckoutOrder201JSONResponse(orderDTO), nil
}

// GetPaymentParams handles POST /checkout/orders/{orderId}/payment-params
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

	paramsDTO := toPaymentParamsResponse(params)
	return openapi.GetPaymentParams200JSONResponse(paramsDTO), nil
}

// PaymentNotify handles POST /checkout/payment-notify
func (h *Handler) PaymentNotify(ctx context.Context, request openapi.PaymentNotifyRequestObject) (openapi.PaymentNotifyResponseObject, error) {
	payload := make(map[string]string)
	if request.Body != nil {
		for k, v := range *request.Body {
			payload[k] = v
		}
	}

	err := h.checkoutUC.PaymentNotify(ctx, payload)
	if err != nil {
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

// GetPaymentStatus handles GET /checkout/orders/{orderId}/payment-status
func (h *Handler) GetPaymentStatus(ctx context.Context, request openapi.GetPaymentStatusRequestObject) (openapi.GetPaymentStatusResponseObject, error) {
	order, err := h.checkoutUC.PaymentStatus(ctx, request.OrderId)
	if err != nil {
		status, errResp := mapCheckoutError(err)
		switch status {
		case 401:
			return openapi.GetPaymentStatus401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.GetPaymentStatus404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetPaymentStatus500JSONResponse{}, nil
		}
	}

	orderDTO := toCheckoutOrderResponse(order)
	return openapi.GetPaymentStatus200JSONResponse(orderDTO), nil
}
