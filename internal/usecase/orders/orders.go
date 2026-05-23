package orders

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/usecase"
)

type UseCase struct {
	repo repo.OrderRepository
}

func New(r repo.OrderRepository) *UseCase {
	return &UseCase{repo: r}
}

var _ usecase.Orders = (*UseCase)(nil)

func (uc *UseCase) List(ctx context.Context, actor entity.Actor, email string, page entity.Pagination) ([]entity.Order, int, error) {
	orders, total, err := uc.repo.ListOrdersByOwner(ctx, actor.UserID, email, page)
	if err != nil {
		return nil, 0, fmt.Errorf("OrdersUseCase - List - repo.ListOrdersByOwner: %w", err)
	}
	for i := range orders {
		orders[i].StatusText, orders[i].StatusColor = mapStatus(orders[i].Status)
	}
	return orders, total, nil
}

func (uc *UseCase) Get(ctx context.Context, actor entity.Actor, orderID string) (entity.Order, error) {
	order, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrdersUseCase - Get - repo.GetOrderByID: %w", err)
	}
	if actor.UserID != "" && order.UserID != "" && actor.UserID != order.UserID && !actor.IsAdmin {
		return entity.Order{}, entity.ErrForbidden
	}
	order.StatusText, order.StatusColor = mapStatus(order.Status)
	return order, nil
}

func (uc *UseCase) RequestRefund(ctx context.Context, actor entity.Actor, orderID, reason string) (entity.RefundRequest, error) {
	order, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return entity.RefundRequest{}, fmt.Errorf("OrdersUseCase - RequestRefund - repo.GetOrderByID: %w", err)
	}
	if !order.CanRequestRefund() {
		return entity.RefundRequest{}, entity.ErrRefundNotAllowed
	}

	refund := entity.RefundRequest{
		OrderID:   orderID,
		UserID:    actor.UserID,
		Username:  actor.Username,
		Reason:    reason,
		Status:    entity.RefundStatusPending,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := uc.repo.SubmitRefundRequest(ctx, refund); err != nil {
		return entity.RefundRequest{}, fmt.Errorf("OrdersUseCase - RequestRefund - repo.SubmitRefundRequest: %w", err)
	}
	return refund, nil
}

func mapStatus(status entity.OrderStatus) (text string, color string) {
	switch status {
	case entity.OrderStatusPending:
		return "Cho thanh toan", "#f59e0b"
	case entity.OrderStatusPaid:
		return "Da thanh toan", "#3b82f6"
	case entity.OrderStatusDelivered:
		return "Da giao hang", "#10b981"
	case entity.OrderStatusCancelled:
		return "Da huy", "#6b7280"
	case entity.OrderStatusFailed:
		return "That bai", "#ef4444"
	case entity.OrderStatusRefundPending:
		return "Cho hoan tien", "#8b5cf6"
	case entity.OrderStatusRefunded:
		return "Da hoan tien", "#14b8a6"
	default:
		return "Khong xac dinh", "#94a3b8"
	}
}
