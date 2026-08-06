package order_test

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/module/order"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/stretchr/testify/require"
)

type orderRepoStub struct {
	listFunc         func(ctx context.Context, userID, email string, page pagination.Pagination) ([]order.Order, int, error)
	getFunc          func(ctx context.Context, orderID string) (order.Order, error)
	submitRefundFunc func(ctx context.Context, refund *order.RefundRequest) error
}

func (s *orderRepoStub) ListOrdersByOwner(ctx context.Context, userID, email string, page pagination.Pagination) ([]order.Order, int, error) {
	if s.listFunc != nil {
		return s.listFunc(ctx, userID, email, page)
	}
	return nil, 0, nil
}

func (s *orderRepoStub) GetOrderByID(ctx context.Context, orderID string) (order.Order, error) {
	if s.getFunc != nil {
		return s.getFunc(ctx, orderID)
	}
	return order.Order{}, nil
}

func (s *orderRepoStub) CancelPendingOrder(context.Context, order.Actor, string) error {
	return nil
}

func (s *orderRepoStub) SubmitRefundRequest(ctx context.Context, refund *order.RefundRequest) error {
	if s.submitRefundFunc != nil {
		return s.submitRefundFunc(ctx, refund)
	}
	return nil
}

func TestOrdersUseCase_Lifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	owner := order.Actor{UserID: "u1", Username: "alice"}
	admin := order.Actor{UserID: "admin", Username: "root", IsAdmin: true}
	other := order.Actor{UserID: "u2", Username: "bob"}

	t.Run("list decorates statuses", func(t *testing.T) {
		uc := order.NewOrdersUseCase(&orderRepoStub{
			listFunc: func(_ context.Context, userID, email string, page pagination.Pagination) ([]order.Order, int, error) {
				require.Equal(t, "u1", userID)
				require.Equal(t, "alice@example.com", email)
				require.Equal(t, 20, page.Limit)
				return []order.Order{{ID: "o1", Status: order.OrderStatusPending}}, 1, nil
			},
		})

		orders, total, err := uc.List(ctx, owner, "alice@example.com", pagination.Pagination{Limit: 20})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Len(t, orders, 1)
		require.Equal(t, "Cho thanh toan", orders[0].StatusText)
		require.Equal(t, "#f59e0b", orders[0].StatusColor)
	})

	t.Run("get blocks foreign non-admin actor", func(t *testing.T) {
		uc := order.NewOrdersUseCase(&orderRepoStub{
			getFunc: func(_ context.Context, orderID string) (order.Order, error) {
				require.Equal(t, "o1", orderID)
				return order.Order{ID: orderID, UserID: "u1", Status: order.OrderStatusDelivered}, nil
			},
		})

		_, err := uc.Get(ctx, other, "o1")
		require.ErrorIs(t, err, order.ErrForbidden)

		ord, err := uc.Get(ctx, admin, "o1")
		require.NoError(t, err)
		require.Equal(t, "Da giao hang", ord.StatusText)
	})

	t.Run("refund request requires auth ownership reason and delivered status", func(t *testing.T) {
		uc := order.NewOrdersUseCase(&orderRepoStub{
			getFunc: func(_ context.Context, orderID string) (order.Order, error) {
				return order.Order{ID: orderID, UserID: "u1", Status: order.OrderStatusPending}, nil
			},
		})

		_, err := uc.RequestRefund(ctx, order.Actor{}, "o1", "late")
		require.ErrorIs(t, err, order.ErrUnauthorized)

		_, err = uc.RequestRefund(ctx, other, "o1", "late")
		require.ErrorIs(t, err, order.ErrForbidden)

		_, err = uc.RequestRefund(ctx, owner, "o1", "")
		require.ErrorIs(t, err, order.ErrInvalidInput)

		_, err = uc.RequestRefund(ctx, owner, "o1", "late")
		require.ErrorIs(t, err, order.ErrRefundNotAllowed)
	})

	t.Run("refund request persists pending refund for delivered order", func(t *testing.T) {
		var submitted order.RefundRequest
		uc := order.NewOrdersUseCase(&orderRepoStub{
			getFunc: func(_ context.Context, orderID string) (order.Order, error) {
				return order.Order{ID: orderID, UserID: "u1", Status: order.OrderStatusDelivered}, nil
			},
			submitRefundFunc: func(_ context.Context, refund *order.RefundRequest) error {
				submitted = *refund
				return nil
			},
		})

		refund, err := uc.RequestRefund(ctx, owner, "o2", "wrong item")
		require.NoError(t, err)
		require.Equal(t, "o2", submitted.OrderID)
		require.Equal(t, "u1", submitted.UserID)
		require.Equal(t, "alice", submitted.Username)
		require.Equal(t, order.RefundStatusPending, submitted.Status)
		require.Equal(t, order.RefundStatusPending, refund.Status)
	})
}
