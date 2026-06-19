package v1

import (
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

const bearerParts = 2

func (r *V1) gripActor(ctx *fiber.Ctx) entity.Actor {
	if actor, ok := ctx.Locals("actor").(entity.Actor); ok {
		return actor
	}

	actor := entity.Actor{TrustLevel: -1}
	if header := ctx.Get("Authorization"); header != "" && r.jwtManager != nil {
		parts := strings.SplitN(header, " ", bearerParts)
		if len(parts) == bearerParts && parts[0] == "Bearer" {
			userID, err := r.jwtManager.ParseToken(parts[1])
			if err == nil {
				actor.UserID = userID
				actor.TrustLevel = 0
			}
		}
	}

	return actor
}

func gripPage(ctx *fiber.Ctx) entity.Pagination {
	return entity.Pagination{
		Limit:  ctx.QueryInt("limit", 20),
		Offset: ctx.QueryInt("offset", 0),
	}
}

func gripDecorateOrderStatus(order entity.Order) entity.Order {
	switch order.Status {
	case entity.OrderStatusPending:
		order.StatusText, order.StatusColor = "Cho thanh toan", "#f59e0b"
	case entity.OrderStatusPaid:
		order.StatusText, order.StatusColor = "Da thanh toan", "#3b82f6"
	case entity.OrderStatusDelivered:
		order.StatusText, order.StatusColor = "Da giao hang", "#10b981"
	case entity.OrderStatusCancelled:
		order.StatusText, order.StatusColor = "Da huy", "#6b7280"
	case entity.OrderStatusFailed:
		order.StatusText, order.StatusColor = "That bai", "#ef4444"
	case entity.OrderStatusRefundPending:
		order.StatusText, order.StatusColor = "Cho hoan tien", "#8b5cf6"
	case entity.OrderStatusRefunded:
		order.StatusText, order.StatusColor = "Da hoan tien", "#14b8a6"
	default:
		order.StatusText, order.StatusColor = "Khong xac dinh", "#94a3b8"
	}

	return order
}
