package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type adminOrderRefundUseCaseStub struct {
	updateOrderFunc          func(ctx context.Context, actor entity.Actor, orderID string, status entity.OrderStatus) error
	deleteOrderFunc          func(ctx context.Context, actor entity.Actor, orderID string) error
	listRefundsFunc          func(ctx context.Context, actor entity.Actor, status string) ([]entity.RefundRequest, error)
	decideRefundFunc         func(ctx context.Context, actor entity.Actor, refundID int64, approve bool, note string) (entity.RefundRequest, error)
	getRefundFunc            func(ctx context.Context, actor entity.Actor, refundID int64) (entity.RefundRequest, error)
	getOrderRefundStatusFunc func(ctx context.Context, actor entity.Actor, orderID string) (entity.RefundRequest, error)
}

func (s *adminOrderRefundUseCaseStub) ListUsers(context.Context, entity.Actor, entity.Pagination) ([]entity.User, int, error) {
	return nil, 0, nil
}

func (s *adminOrderRefundUseCaseStub) UpdateUserStatus(context.Context, entity.Actor, string, entity.UserStatus) error {
	return nil
}

func (s *adminOrderRefundUseCaseStub) UpdateUserPoints(context.Context, entity.Actor, string, int) error {
	return nil
}

func (s *adminOrderRefundUseCaseStub) ListOrders(context.Context, entity.Actor, entity.Pagination, string, string) ([]entity.Order, int, error) {
	return nil, 0, nil
}

func (s *adminOrderRefundUseCaseStub) GetOrder(context.Context, entity.Actor, string) (entity.Order, error) {
	return entity.Order{}, nil
}

func (s *adminOrderRefundUseCaseStub) UpdateOrderStatus(ctx context.Context, actor entity.Actor, orderID string, status entity.OrderStatus) error {
	if s.updateOrderFunc != nil {
		return s.updateOrderFunc(ctx, actor, orderID, status)
	}
	return nil
}

func (s *adminOrderRefundUseCaseStub) DeleteOrder(ctx context.Context, actor entity.Actor, orderID string) error {
	if s.deleteOrderFunc != nil {
		return s.deleteOrderFunc(ctx, actor, orderID)
	}
	return nil
}

func (s *adminOrderRefundUseCaseStub) ListRefunds(ctx context.Context, actor entity.Actor, status string) ([]entity.RefundRequest, error) {
	if s.listRefundsFunc != nil {
		return s.listRefundsFunc(ctx, actor, status)
	}
	return nil, nil
}

func (s *adminOrderRefundUseCaseStub) DecideRefund(ctx context.Context, actor entity.Actor, refundID int64, approve bool, note string) (entity.RefundRequest, error) {
	if s.decideRefundFunc != nil {
		return s.decideRefundFunc(ctx, actor, refundID, approve, note)
	}
	return entity.RefundRequest{}, nil
}

func (s *adminOrderRefundUseCaseStub) GetRefund(ctx context.Context, actor entity.Actor, refundID int64) (entity.RefundRequest, error) {
	if s.getRefundFunc != nil {
		return s.getRefundFunc(ctx, actor, refundID)
	}
	return entity.RefundRequest{}, nil
}

func (s *adminOrderRefundUseCaseStub) GetOrderRefundStatus(ctx context.Context, actor entity.Actor, orderID string) (entity.RefundRequest, error) {
	if s.getOrderRefundStatusFunc != nil {
		return s.getOrderRefundStatusFunc(ctx, actor, orderID)
	}
	return entity.RefundRequest{}, nil
}

func (s *adminOrderRefundUseCaseStub) ListCards(ctx context.Context, actor entity.Actor) ([]entity.Card, error) {
	return nil, nil
}

func (s *adminOrderRefundUseCaseStub) ListReviews(context.Context, entity.Actor, entity.Pagination, string, string) ([]entity.Review, repo.ReviewModerationStats, int, error) {
	return nil, repo.ReviewModerationStats{}, 0, nil
}

func (s *adminOrderRefundUseCaseStub) UpdateReviewStatus(context.Context, entity.Actor, int64, entity.ReviewStatus) (entity.Review, error) {
	return entity.Review{}, nil
}

func (s *adminOrderRefundUseCaseStub) BulkPublishReviews(context.Context, entity.Actor, []int64) (int, error) {
	return 0, nil
}

func (s *adminOrderRefundUseCaseStub) DeleteReview(context.Context, entity.Actor, int64) error {
	return nil
}

func (s *adminOrderRefundUseCaseStub) RepairAggregates(context.Context, entity.Actor) error {
	return nil
}

func (s *adminOrderRefundUseCaseStub) ListProducts(context.Context, entity.Actor, entity.Pagination) ([]entity.Product, int, error) {
	return nil, 0, nil
}

func (s *adminOrderRefundUseCaseStub) GetProduct(context.Context, entity.Actor, string) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *adminOrderRefundUseCaseStub) UpsertProduct(context.Context, entity.Actor, entity.Product) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *adminOrderRefundUseCaseStub) DeleteProduct(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminOrderRefundUseCaseStub) ListCategories(context.Context, entity.Actor) ([]entity.Category, error) {
	return nil, nil
}

func (s *adminOrderRefundUseCaseStub) UpsertCategory(context.Context, entity.Actor, entity.Category) (entity.Category, error) {
	return entity.Category{}, nil
}

func (s *adminOrderRefundUseCaseStub) DeleteCategory(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminOrderRefundUseCaseStub) ListSettings(context.Context, entity.Actor) ([]entity.Setting, error) {
	return nil, nil
}

func (s *adminOrderRefundUseCaseStub) SetSetting(context.Context, entity.Actor, string, string) error {
	return nil
}

func (s *adminOrderRefundUseCaseStub) DeleteSetting(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminOrderRefundUseCaseStub) SendBroadcast(context.Context, entity.Actor, string, string) error {
	return nil
}

func (s *adminOrderRefundUseCaseStub) SendTargeted(context.Context, entity.Actor, string, string, string) error {
	return nil
}

func TestAdminOrderAndRefundEndpoints(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin-id", "admin", true)
	require.NoError(t, err)
	userToken, err := jwtManager.GenerateTokenWithProfile("user-id", "user", false)
	require.NoError(t, err)

	setupApp := func(uc *adminOrderRefundUseCaseStub) *fiber.App {
		v := &V1{
			adminUC:    uc,
			jwtManager: jwtManager,
			adminUsers: "admin",
		}
		app := fiber.New()
		v.registerGripStoreRoutes(app.Group("/v1"))
		return app
	}

	t.Run("patch order status requires admin and forwards payload", func(t *testing.T) {
		var gotStatus entity.OrderStatus
		app := setupApp(&adminOrderRefundUseCaseStub{
			updateOrderFunc: func(_ context.Context, actor entity.Actor, orderID string, status entity.OrderStatus) error {
				require.True(t, actor.IsAdmin)
				require.Equal(t, "o1", orderID)
				gotStatus = status
				return nil
			},
		})

		require.Equal(t, http.StatusForbidden, testRequest(t, app, http.MethodPatch, "/v1/admin/orders/o1", []byte(`{"status":"cancelled"}`), "Bearer "+userToken).StatusCode)
		require.Equal(t, http.StatusBadRequest, testRequest(t, app, http.MethodPatch, "/v1/admin/orders/o1", []byte(`{bad}`), "Bearer "+adminToken).StatusCode)
		require.Equal(t, http.StatusNoContent, testRequest(t, app, http.MethodPatch, "/v1/admin/orders/o1", []byte(`{"status":"cancelled"}`), "Bearer "+adminToken).StatusCode)
		require.Equal(t, entity.OrderStatusCancelled, gotStatus)
	})

	t.Run("delete order forwards target", func(t *testing.T) {
		var deleted string
		app := setupApp(&adminOrderRefundUseCaseStub{
			deleteOrderFunc: func(_ context.Context, actor entity.Actor, orderID string) error {
				require.True(t, actor.IsAdmin)
				deleted = orderID
				return nil
			},
		})

		resp := testRequest(t, app, http.MethodDelete, "/v1/admin/orders/o2", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		require.Equal(t, "o2", deleted)
	})

	t.Run("list refunds returns data", func(t *testing.T) {
		app := setupApp(&adminOrderRefundUseCaseStub{
			listRefundsFunc: func(_ context.Context, actor entity.Actor, status string) ([]entity.RefundRequest, error) {
				require.True(t, actor.IsAdmin)
				require.Equal(t, "pending", status)
				return []entity.RefundRequest{{ID: 10, Status: entity.RefundStatusPending, OrderID: "o1"}}, nil
			},
		})

		resp := testRequest(t, app, http.MethodGet, "/v1/admin/refunds?status=pending", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body envelope
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		raw, err := json.Marshal(body.Data)
		require.NoError(t, err)
		var refunds []entity.RefundRequest
		require.NoError(t, json.Unmarshal(raw, &refunds))
		require.Len(t, refunds, 1)
	})

	t.Run("approve and reject refund parse body and route decision", func(t *testing.T) {
		var decisions []struct {
			ID      int64
			Approve bool
			Note    string
		}
		app := setupApp(&adminOrderRefundUseCaseStub{
			decideRefundFunc: func(_ context.Context, actor entity.Actor, refundID int64, approve bool, note string) (entity.RefundRequest, error) {
				require.True(t, actor.IsAdmin)
				decisions = append(decisions, struct {
					ID      int64
					Approve bool
					Note    string
				}{ID: refundID, Approve: approve, Note: note})
				status := entity.RefundStatusRejected
				if approve {
					status = entity.RefundStatusApproved
				}
				return entity.RefundRequest{ID: refundID, Status: status, AdminNote: note}, nil
			},
		})

		require.Equal(t, http.StatusBadRequest, testRequest(t, app, http.MethodPost, "/v1/admin/refunds/not-a-number/approve", []byte(`{}`), "Bearer "+adminToken).StatusCode)
		require.Equal(t, http.StatusBadRequest, testRequest(t, app, http.MethodPost, "/v1/admin/refunds/4/approve", []byte(`{bad}`), "Bearer "+adminToken).StatusCode)

		require.Equal(t, http.StatusOK, testRequest(t, app, http.MethodPost, "/v1/admin/refunds/4/approve", []byte(`{"note":"ok"}`), "Bearer "+adminToken).StatusCode)
		require.Equal(t, http.StatusOK, testRequest(t, app, http.MethodPost, "/v1/admin/refunds/4/reject", []byte(`{"note":"no"}`), "Bearer "+adminToken).StatusCode)
		require.Len(t, decisions, 2)
		require.True(t, decisions[0].Approve)
		require.False(t, decisions[1].Approve)
	})
}
