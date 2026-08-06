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

// GetCheckoutPreview handles GET /checkout/preview
func (h *Handler) GetCheckoutPreview(ctx context.Context, request openapi.GetCheckoutPreviewRequestObject) (openapi.GetCheckoutPreviewResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.GetCheckoutPreview401JSONResponse{}, nil
	}

	productID := ""
	if request.Params.ProductId != nil {
		productID = *request.Params.ProductId
	}
	quantity := 1
	if request.Params.Quantity != nil && *request.Params.Quantity > 0 {
		quantity = *request.Params.Quantity
	}

	breakdown, err := h.checkoutUC.Preview(ctx, actor, productID, quantity)
	if err != nil {
		status, errResp := mapCheckoutError(err)
		switch status {
		case 400:
			return openapi.GetCheckoutPreview400JSONResponse{}, nil
		default:
			return openapi.GetCheckoutPreview500JSONResponse{
				InternalErrorResponseJSONResponse: openapi.InternalErrorResponseJSONResponse(errResp),
			}, nil
		}
	}

	previewDTO := toCheckoutPreviewResponse(breakdown)
	return openapi.GetCheckoutPreview200JSONResponse(previewDTO), nil
}

// PreviewCheckout handles POST /checkout/preview
func (h *Handler) PreviewCheckout(ctx context.Context, request openapi.PreviewCheckoutRequestObject) (openapi.PreviewCheckoutResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.PreviewCheckout401JSONResponse{}, nil
	}
	if request.Body == nil {
		return openapi.PreviewCheckout400JSONResponse{}, nil
	}

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
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.CreateCheckoutOrder401JSONResponse{}, nil
	}
	if request.Body == nil {
		return openapi.CreateCheckoutOrder400JSONResponse{}, nil
	}

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
		case 400, 404:
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
	if actor.UserID == "" {
		return openapi.GetPaymentParams401JSONResponse{}, nil
	}

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

// PostPaymentParams handles POST /checkout/orders/{orderId}/payment-params
func (h *Handler) PostPaymentParams(ctx context.Context, request openapi.PostPaymentParamsRequestObject) (openapi.PostPaymentParamsResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.PostPaymentParams401JSONResponse{}, nil
	}

	params, err := h.checkoutUC.PaymentParams(ctx, actor, request.OrderId)
	if err != nil {
		status, errResp := mapCheckoutError(err)
		switch status {
		case 401:
			return openapi.PostPaymentParams401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.PostPaymentParams404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.PostPaymentParams500JSONResponse{}, nil
		}
	}

	res := toPaymentParamsResponse(params)
	return openapi.PostPaymentParams200JSONResponse(res), nil
}

// CreatePaymentOrder handles POST /checkout/payment-orders
func (h *Handler) CreatePaymentOrder(ctx context.Context, request openapi.CreatePaymentOrderRequestObject) (openapi.CreatePaymentOrderResponseObject, error) {
	_ = request
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.CreatePaymentOrder401JSONResponse{}, nil
	}
	return openapi.CreatePaymentOrder500JSONResponse{}, nil
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
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.GetPaymentStatus401JSONResponse{}, nil
	}

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
