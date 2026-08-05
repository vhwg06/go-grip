package cart

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// toCartItemResponse maps entity.CartItem to openapi.CartItemResponse DTO.
func toCartItemResponse(item entity.CartItem) openapi.CartItemResponse {
	priceInt := int(item.UnitPrice)
	blocked := item.Blocked
	return openapi.CartItemResponse{
		Id:        item.ID,
		ProductId: item.ProductID,
		Quantity:  item.Quantity,
		UnitPrice: priceInt,
		Blocked:   &blocked,
	}
}

// toCartResponse maps entity.Cart to openapi.CartResponse DTO.
func toCartResponse(c entity.Cart) openapi.CartResponse {
	itemsDTO := make([]openapi.CartItemResponse, len(c.Items))
	for i, item := range c.Items {
		itemsDTO[i] = toCartItemResponse(item)
	}

	statusStr := string(c.Status)
	sessionID := c.SessionID

	return openapi.CartResponse{
		Id:        c.ID,
		SessionId: &sessionID,
		Status:    statusStr,
		Items:     itemsDTO,
		CreatedAt: &c.CreatedAt,
		UpdatedAt: &c.UpdatedAt,
	}
}
