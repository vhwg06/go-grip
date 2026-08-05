package orders

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Orders capability.
type Handler struct {
	ordersUC usecase.Orders
	logger   logger.Interface
}

// NewHandler constructs a new Orders vertical handler instance.
func NewHandler(ordersUC usecase.Orders, l logger.Interface) *Handler {
	return &Handler{
		ordersUC: ordersUC,
		logger:   l,
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

	page := entity.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	orders, total, err := h.ordersUC.List(ctx, actor, "", page)
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

	listDTO := toOrderListResponse(orders, total)
	return openapi.ListOrders200JSONResponse(listDTO), nil
}

// GetOrderByID handles GET /orders/{id}
func (h *Handler) GetOrderByID(ctx context.Context, request openapi.GetOrderByIDRequestObject) (openapi.GetOrderByIDResponseObject, error) {
	actor := getActor(ctx)
	order, err := h.ordersUC.Get(ctx, actor, request.Id)
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

	orderDTO := toOrderResponse(order)
	return openapi.GetOrderByID200JSONResponse(orderDTO), nil
}

// RequestOrderRefund handles POST /orders/{id}/refund
func (h *Handler) RequestOrderRefund(ctx context.Context, request openapi.RequestOrderRefundRequestObject) (openapi.RequestOrderRefundResponseObject, error) {
	if request.Body == nil {
		return openapi.RequestOrderRefund400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	rf, err := h.ordersUC.RequestRefund(ctx, actor, request.Id, request.Body.Reason)
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
			return openapi.RequestOrderRefund422JSONResponse{
				UnprocessableEntityResponseJSONResponse: openapi.UnprocessableEntityResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.RequestOrderRefund500JSONResponse{}, nil
		}
	}

	rfDTO := toRefundResponse(rf)
	return openapi.RequestOrderRefund200JSONResponse(rfDTO), nil
}
