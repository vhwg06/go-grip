package orders

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Orders capability.
type Handler struct {
	ordersUC ordermodule.OrdersUseCase
	logger   logger.Interface
}

// NewHandler constructs a new Orders vertical handler instance.
func NewHandler(ordersUC ordermodule.OrdersUseCase, l logger.Interface) *Handler {
	return &Handler{
		ordersUC: ordersUC,
		logger:   l,
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

// ListOrders handles GET /orders
func (h *Handler) ListOrders(ctx context.Context, request openapi.ListOrdersRequestObject) (openapi.ListOrdersResponseObject, error) {
	actor := getActor(ctx)
	limit := 10
	if request.Params.Limit != nil && *request.Params.Limit > 0 {
		limit = *request.Params.Limit
	}
	offset := 0
	if request.Params.Offset != nil && *request.Params.Offset >= 0 {
		offset = *request.Params.Offset
	}

	page := pagination.Pagination{
		Limit:  limit,
		Offset: offset,
	}
	email := ""

	orderList, total, err := h.ordersUC.List(ctx, ordermodule.Actor{UserID: actor.UserID, Username: actor.Username, IsAdmin: actor.IsAdmin}, email, page)
	if err != nil {
		status, errResp := mapOrdersError(err)
		switch status {
		case 401:
			return openapi.ListOrders401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.ListOrders500JSONResponse{}, nil
		}
	}

	res := toOrderListResponse(orderList, total)
	return openapi.ListOrders200JSONResponse(res), nil
}

// GetOrderByID handles GET /orders/{id}
func (h *Handler) GetOrderByID(ctx context.Context, request openapi.GetOrderByIDRequestObject) (openapi.GetOrderByIDResponseObject, error) {
	actor := getActor(ctx)
	orderEntity, err := h.ordersUC.Get(ctx, ordermodule.Actor{UserID: actor.UserID, Username: actor.Username, IsAdmin: actor.IsAdmin}, request.Id)
	if err != nil {
		status, errResp := mapOrdersError(err)
		switch status {
		case 401:
			return openapi.GetOrderByID401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.GetOrderByID404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetOrderByID500JSONResponse{}, nil
		}
	}

	res := toOrderResponse(orderEntity)
	return openapi.GetOrderByID200JSONResponse(res), nil
}

// RequestOrderRefund handles POST /orders/{id}/refund
func (h *Handler) RequestOrderRefund(ctx context.Context, request openapi.RequestOrderRefundRequestObject) (openapi.RequestOrderRefundResponseObject, error) {
	actor := getActor(ctx)
	reason := ""
	if request.Body != nil {
		reason = request.Body.Reason
	}

	refundEntity, err := h.ordersUC.RequestRefund(ctx, ordermodule.Actor{UserID: actor.UserID, Username: actor.Username, IsAdmin: actor.IsAdmin}, request.Id, reason)
	if err != nil {
		status, errResp := mapOrdersError(err)
		switch status {
		case 400:
			return openapi.RequestOrderRefund400JSONResponse{}, nil
		case 401:
			return openapi.RequestOrderRefund401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.RequestOrderRefund404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		case 422:
			return openapi.RequestOrderRefund422JSONResponse{}, nil
		default:
			return openapi.RequestOrderRefund500JSONResponse{}, nil
		}
	}

	res := toRefundResponse(refundEntity)
	return openapi.RequestOrderRefund200JSONResponse(res), nil
}
