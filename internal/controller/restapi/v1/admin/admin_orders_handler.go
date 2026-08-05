package admin

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// AdminListOrders handles GET /admin/orders
func (h *Handler) AdminListOrders(ctx context.Context, request openapi.AdminListOrdersRequestObject) (openapi.AdminListOrdersResponseObject, error) {
	actor := getActor(ctx)

	page := 1
	pageSize := 20
	query := ""
	status := ""

	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}
	if request.Params.Q != nil {
		query = *request.Params.Q
	}
	if request.Params.Status != nil {
		status = *request.Params.Status
	}

	offset := (page - 1) * pageSize
	pag := pagination.New(pageSize, offset)

	orders, total, err := h.adminUC.ListOrders(ctx, actor, pag, query, status)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminListOrders401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminListOrders403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminListOrders500JSONResponse{}, nil
		}
	}

	items := make([]openapi.AdminOrderDetailResponse, 0, len(orders))
	for _, o := range orders {
		items = append(items, toAdminOrderDetail(o))
	}

	totalInt := total
	pageInt := page
	sizeInt := pageSize
	resp := openapi.AdminOrderListResponse{
		Items:    &items,
		Total:    &totalInt,
		Page:     &pageInt,
		PageSize: &sizeInt,
	}
	return openapi.AdminListOrders200JSONResponse(resp), nil
}

// AdminGetOrder handles GET /admin/orders/{orderId}
func (h *Handler) AdminGetOrder(ctx context.Context, request openapi.AdminGetOrderRequestObject) (openapi.AdminGetOrderResponseObject, error) {
	actor := getActor(ctx)

	order, err := h.adminUC.GetOrder(ctx, actor, request.OrderId)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminGetOrder401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminGetOrder403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		case 404:
			return openapi.AdminGetOrder404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminGetOrder500JSONResponse{}, nil
		}
	}

	detail := toAdminOrderDetail(order)
	return openapi.AdminGetOrder200JSONResponse(detail), nil
}

// AdminUpdateOrder handles PATCH /admin/orders/{orderId}
func (h *Handler) AdminUpdateOrder(ctx context.Context, request openapi.AdminUpdateOrderRequestObject) (openapi.AdminUpdateOrderResponseObject, error) {
	actor := getActor(ctx)

	var orderStatus entity.OrderStatus
	if request.Body != nil && request.Body.Status != nil {
		orderStatus = entity.OrderStatus(*request.Body.Status)
	}

	if err := h.adminUC.UpdateOrderStatus(ctx, actor, request.OrderId, orderStatus); err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminUpdateOrder401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminUpdateOrder403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		case 404:
			return openapi.AdminUpdateOrder404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminUpdateOrder500JSONResponse{}, nil
		}
	}

	// Return the updated order detail.
	order, err := h.adminUC.GetOrder(ctx, actor, request.OrderId)
	if err != nil {
		return openapi.AdminUpdateOrder500JSONResponse{}, nil
	}
	detail := toAdminOrderDetail(order)
	return openapi.AdminUpdateOrder200JSONResponse(detail), nil
}

// AdminGetCollect handles GET /admin/collect — returns aggregated backoffice counts.
func (h *Handler) AdminGetCollect(ctx context.Context, _ openapi.AdminGetCollectRequestObject) (openapi.AdminGetCollectResponseObject, error) {
	actor := getActor(ctx)

	// Confirm admin access by attempting listing; propagate auth errors.
	_, userCount, err := h.adminUC.ListUsers(ctx, actor, pagination.New(1, 0))
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminGetCollect401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminGetCollect403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminGetCollect500JSONResponse{}, nil
		}
	}

	_, orderCount, _ := h.adminUC.ListOrders(ctx, actor, pagination.New(1, 0), "", "")

	// AdminCollectResponse is map[string]interface{} — cast directly.
	resp := openapi.AdminCollectResponse{
		"users":  userCount,
		"orders": orderCount,
	}
	return openapi.AdminGetCollect200JSONResponse(resp), nil
}


// toAdminOrderDetail maps entity.Order to openapi.AdminOrderDetailResponse.
func toAdminOrderDetail(o entity.Order) openapi.AdminOrderDetailResponse {
	orderID := o.ID
	userID := o.UserID
	username := o.Username
	email := o.Email
	productName := o.ProductName
	amount := int64(o.Amount)
	qty := o.Quantity
	status := string(o.Status)

	return openapi.AdminOrderDetailResponse{
		OrderId:     &orderID,
		UserId:      &userID,
		Username:    &username,
		Email:       &email,
		ProductName: &productName,
		Amount:      &amount,
		Quantity:    &qty,
		Status:      &status,
		CreatedAt:   &o.CreatedAt,
		UpdatedAt:   &o.UpdatedAt,
	}
}
