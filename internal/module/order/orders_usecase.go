package order

import (
	"context"
	"fmt"
	"time"

	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// OrdersUseCase defines application services for querying orders and requesting refunds.
type OrdersUseCase interface {
	List(ctx context.Context, actor usermodule.Actor, email string, page pagination.Pagination) ([]Order, int, error)
	Get(ctx context.Context, actor usermodule.Actor, orderID string) (Order, error)
	RequestRefund(ctx context.Context, actor usermodule.Actor, orderID, reason string) (RefundRequest, error)
}

type ordersUseCase struct {
	repo OrderRepo
}

// NewOrdersUseCase constructs a new OrdersUseCase instance.
func NewOrdersUseCase(r OrderRepo) OrdersUseCase {
	return &ordersUseCase{repo: r}
}

func (uc *ordersUseCase) List(ctx context.Context, actor usermodule.Actor, email string, page pagination.Pagination) ([]Order, int, error) {
	orders, total, err := uc.repo.ListOrdersByOwner(ctx, actor.UserID, email, page)
	if err != nil {
		return nil, 0, fmt.Errorf("OrdersUseCase - List - repo.ListOrdersByOwner: %w", err)
	}
	for i := range orders {
		orders[i].StatusText, orders[i].StatusColor = mapStatus(orders[i].Status)
	}
	return orders, total, nil
}

func (uc *ordersUseCase) Get(ctx context.Context, actor usermodule.Actor, orderID string) (Order, error) {
	o, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return Order{}, fmt.Errorf("OrdersUseCase - Get - repo.GetOrderByID: %w", err)
	}
	if actor.UserID != "" && o.UserID != "" && actor.UserID != o.UserID && !actor.IsAdmin {
		return Order{}, ErrForbidden
	}
	o.StatusText, o.StatusColor = mapStatus(o.Status)
	return o, nil
}

func (uc *ordersUseCase) RequestRefund(ctx context.Context, actor usermodule.Actor, orderID, reason string) (RefundRequest, error) {
	o, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return RefundRequest{}, fmt.Errorf("OrdersUseCase - RequestRefund - repo.GetOrderByID: %w", err)
	}
	if actor.UserID == "" {
		return RefundRequest{}, ErrUnauthorized
	}
	if o.UserID != "" && o.UserID != actor.UserID && !actor.IsAdmin {
		return RefundRequest{}, ErrForbidden
	}
	if reason == "" {
		return RefundRequest{}, ErrInvalidInput
	}
	if !o.CanRequestRefund() {
		return RefundRequest{}, ErrRefundNotAllowed
	}

	refund := RefundRequest{
		OrderID:   orderID,
		UserID:    actor.UserID,
		Username:  actor.Username,
		Reason:    reason,
		Status:    RefundStatusPending,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := uc.repo.SubmitRefundRequest(ctx, &refund); err != nil {
		return RefundRequest{}, fmt.Errorf("OrdersUseCase - RequestRefund - repo.SubmitRefundRequest: %w", err)
	}
	return refund, nil
}

func mapStatus(status OrderStatus) (text string, color string) {
	switch status {
	case OrderStatusPending:
		return "Cho thanh toan", "#f59e0b"
	case OrderStatusPaid:
		return "Da thanh toan", "#3b82f6"
	case OrderStatusDelivered:
		return "Da giao hang", "#10b981"
	case OrderStatusCancelled:
		return "Da huy", "#6b7280"
	case OrderStatusFailed:
		return "That bai", "#ef4444"
	case OrderStatusRefundPending:
		return "Cho hoan tien", "#8b5cf6"
	case OrderStatusRefunded:
		return "Da hoan tien", "#14b8a6"
	default:
		return "Khong xac dinh", "#94a3b8"
	}
}
