package orders

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

type orderRepoStub struct {
	listFunc         func(ctx context.Context, userID, email string, page entity.Pagination) ([]entity.Order, int, error)
	getFunc          func(ctx context.Context, orderID string) (entity.Order, error)
	submitRefundFunc func(ctx context.Context, refund entity.RefundRequest) error
}

func (s *orderRepoStub) ListOrdersByOwner(ctx context.Context, userID, email string, page entity.Pagination) ([]entity.Order, int, error) {
	if s.listFunc != nil {
		return s.listFunc(ctx, userID, email, page)
	}
	return nil, 0, nil
}

func (s *orderRepoStub) GetOrderByID(ctx context.Context, orderID string) (entity.Order, error) {
	if s.getFunc != nil {
		return s.getFunc(ctx, orderID)
	}
	return entity.Order{}, nil
}

func (s *orderRepoStub) CancelPendingOrder(context.Context, entity.Actor, string) error {
	return nil
}

func (s *orderRepoStub) SubmitRefundRequest(ctx context.Context, refund entity.RefundRequest) error {
	if s.submitRefundFunc != nil {
		return s.submitRefundFunc(ctx, refund)
	}
	return nil
}

func TestUseCase_US3OrderLifecycle_TDD(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	owner := entity.Actor{UserID: "u1", Username: "alice"}
	admin := entity.Actor{UserID: "admin", Username: "root", IsAdmin: true}
	other := entity.Actor{UserID: "u2", Username: "bob"}

	t.Run("list decorates statuses", func(t *testing.T) {
		uc := New(&orderRepoStub{
			listFunc: func(_ context.Context, userID, email string, page entity.Pagination) ([]entity.Order, int, error) {
				require.Equal(t, "u1", userID)
				require.Equal(t, "alice@example.com", email)
				require.Equal(t, 20, page.Limit)
				return []entity.Order{{ID: "o1", Status: entity.OrderStatusPending}}, 1, nil
			},
		})

		orders, total, err := uc.List(ctx, owner, "alice@example.com", entity.Pagination{Limit: 20})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Len(t, orders, 1)
		require.Equal(t, "Cho thanh toan", orders[0].StatusText)
		require.Equal(t, "#f59e0b", orders[0].StatusColor)
	})

	t.Run("get blocks foreign non-admin actor", func(t *testing.T) {
		uc := New(&orderRepoStub{
			getFunc: func(_ context.Context, orderID string) (entity.Order, error) {
				require.Equal(t, "o1", orderID)
				return entity.Order{ID: orderID, UserID: "u1", Status: entity.OrderStatusDelivered}, nil
			},
		})

		_, err := uc.Get(ctx, other, "o1")
		require.ErrorIs(t, err, entity.ErrForbidden)

		order, err := uc.Get(ctx, admin, "o1")
		require.NoError(t, err)
		require.Equal(t, "Da giao hang", order.StatusText)
	})

	t.Run("refund request requires auth ownership reason and delivered status", func(t *testing.T) {
		uc := New(&orderRepoStub{
			getFunc: func(_ context.Context, orderID string) (entity.Order, error) {
				return entity.Order{ID: orderID, UserID: "u1", Status: entity.OrderStatusPending}, nil
			},
		})

		_, err := uc.RequestRefund(ctx, entity.Actor{}, "o1", "late")
		require.ErrorIs(t, err, entity.ErrUnauthorized)

		_, err = uc.RequestRefund(ctx, other, "o1", "late")
		require.ErrorIs(t, err, entity.ErrForbidden)

		_, err = uc.RequestRefund(ctx, owner, "o1", "")
		require.ErrorIs(t, err, entity.ErrInvalidInput)

		_, err = uc.RequestRefund(ctx, owner, "o1", "late")
		require.ErrorIs(t, err, entity.ErrRefundNotAllowed)
	})

	t.Run("refund request persists pending refund for delivered order", func(t *testing.T) {
		var submitted entity.RefundRequest
		uc := New(&orderRepoStub{
			getFunc: func(_ context.Context, orderID string) (entity.Order, error) {
				return entity.Order{ID: orderID, UserID: "u1", Status: entity.OrderStatusDelivered}, nil
			},
			submitRefundFunc: func(_ context.Context, refund entity.RefundRequest) error {
				submitted = refund
				return nil
			},
		})

		refund, err := uc.RequestRefund(ctx, owner, "o2", "wrong item")
		require.NoError(t, err)
		require.Equal(t, "o2", submitted.OrderID)
		require.Equal(t, "u1", submitted.UserID)
		require.Equal(t, "alice", submitted.Username)
		require.Equal(t, entity.RefundStatusPending, submitted.Status)
		require.Equal(t, entity.RefundStatusPending, refund.Status)
	})
}
