package order

import (
	"context"

	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// OrderRepo defines persistence operations for Order querying and refund submission.
type OrderRepo interface {
	ListOrdersByOwner(ctx context.Context, userID, email string, page pagination.Pagination) ([]Order, int, error)
	GetOrderByID(ctx context.Context, id string) (Order, error)
	SubmitRefundRequest(ctx context.Context, refund *RefundRequest) error
	CancelPendingOrder(ctx context.Context, actor Actor, orderID string) error
}

// CheckoutRepo defines persistence operations for order creation and payment management.
type CheckoutRepo interface {
	CreateOrderWithReservation(ctx context.Context, actor Actor, order Order) (Order, error)
	AttachPayment(ctx context.Context, payment Payment) error
	UpdateOrderStatus(ctx context.Context, orderID string, status OrderStatus) error
	ReleaseReservation(ctx context.Context, orderID string) error
}
