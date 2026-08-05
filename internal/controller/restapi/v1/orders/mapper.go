package orders

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
)

// toOrderResponse maps ordermodule.Order to openapi.OrderResponse DTO.
func toOrderResponse(o ordermodule.Order) openapi.OrderResponse {
	totalInt := int(o.Amount)
	statusStr := string(o.Status)
	email := o.Email

	return openapi.OrderResponse{
		Id:          o.ID,
		OrderSn:     o.ID,
		Status:      statusStr,
		TotalAmount: totalInt,
		Email:       &email,
		CreatedAt:   &o.CreatedAt,
	}
}

// toOrderListResponse maps []ordermodule.Order and total count to openapi.OrderListResponse DTO.
func toOrderListResponse(orders []ordermodule.Order, total int) openapi.OrderListResponse {
	items := make([]openapi.OrderResponse, len(orders))
	for i, o := range orders {
		items[i] = toOrderResponse(o)
	}

	return openapi.OrderListResponse{
		Items: items,
		Total: total,
	}
}

// toRefundResponse maps ordermodule.RefundRequest to openapi.RefundResponse DTO.
func toRefundResponse(r ordermodule.RefundRequest) openapi.RefundResponse {
	idInt := int(r.ID)
	statusStr := string(r.Status)
	return openapi.RefundResponse{
		Id:        idInt,
		OrderId:   r.OrderID,
		Reason:    r.Reason,
		Status:    statusStr,
		CreatedAt: &r.CreatedAt,
	}
}
