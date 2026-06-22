package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "github.com/evrone/go-clean-template/internal/controller/restapi/v1"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type adminContractEnvelope struct {
	Data  json.RawMessage `json:"data"`
	Error string          `json:"error"`
}

type adminContractAdminStub struct {
	listOrdersFunc        func(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Order, int, error)
	getOrderFunc          func(ctx context.Context, actor entity.Actor, orderID string) (entity.Order, error)
	updateOrderStatusFunc func(ctx context.Context, actor entity.Actor, orderID string, status entity.OrderStatus) error
	deleteOrderFunc       func(ctx context.Context, actor entity.Actor, orderID string) error
	listRefundsFunc       func(ctx context.Context, actor entity.Actor, status string) ([]entity.RefundRequest, error)
	decideRefundFunc      func(ctx context.Context, actor entity.Actor, refundID int64, approve bool, note string) (entity.RefundRequest, error)

	listSettingsFunc      func(ctx context.Context, actor entity.Actor) ([]entity.Setting, error)
	setSettingFunc        func(ctx context.Context, actor entity.Actor, key, value string) error
	deleteSettingFunc     func(ctx context.Context, actor entity.Actor, key string) error
	broadcastFunc         func(ctx context.Context, actor entity.Actor, title, body string) error
	targetedFunc          func(ctx context.Context, actor entity.Actor, userID, title, body string) error
	listReviewsFunc       func(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error)
	updateReviewFunc      func(ctx context.Context, actor entity.Actor, reviewID int64, status entity.ReviewStatus) (entity.Review, error)
	bulkReviewFunc        func(ctx context.Context, actor entity.Actor, reviewIDs []int64) (int, error)
	deleteReviewFunc      func(ctx context.Context, actor entity.Actor, reviewID int64) error
}

func (s *adminContractAdminStub) ListUsers(context.Context, entity.Actor, entity.Pagination) ([]entity.User, int, error) {
	return nil, 0, nil
}

func (s *adminContractAdminStub) UpdateUserStatus(context.Context, entity.Actor, string, entity.UserStatus) error {
	return nil
}

func (s *adminContractAdminStub) ListOrders(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Order, int, error) {
	if s.listOrdersFunc != nil {
		return s.listOrdersFunc(ctx, actor, page, query, status)
	}
	return nil, 0, nil
}

func (s *adminContractAdminStub) GetOrder(ctx context.Context, actor entity.Actor, orderID string) (entity.Order, error) {
	if s.getOrderFunc != nil {
		return s.getOrderFunc(ctx, actor, orderID)
	}
	return entity.Order{}, nil
}

func (s *adminContractAdminStub) UpdateOrderStatus(ctx context.Context, actor entity.Actor, orderID string, status entity.OrderStatus) error {
	if s.updateOrderStatusFunc != nil {
		return s.updateOrderStatusFunc(ctx, actor, orderID, status)
	}
	return nil
}

func (s *adminContractAdminStub) DeleteOrder(ctx context.Context, actor entity.Actor, orderID string) error {
	if s.deleteOrderFunc != nil {
		return s.deleteOrderFunc(ctx, actor, orderID)
	}
	return nil
}

func (s *adminContractAdminStub) ListRefunds(ctx context.Context, actor entity.Actor, status string) ([]entity.RefundRequest, error) {
	if s.listRefundsFunc != nil {
		return s.listRefundsFunc(ctx, actor, status)
	}
	return nil, nil
}

func (s *adminContractAdminStub) DecideRefund(ctx context.Context, actor entity.Actor, refundID int64, approve bool, note string) (entity.RefundRequest, error) {
	if s.decideRefundFunc != nil {
		return s.decideRefundFunc(ctx, actor, refundID, approve, note)
	}
	return entity.RefundRequest{}, nil
}

func (s *adminContractAdminStub) GetRefund(ctx context.Context, actor entity.Actor, refundID int64) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}

func (s *adminContractAdminStub) GetOrderRefundStatus(ctx context.Context, actor entity.Actor, orderID string) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}

func (s *adminContractAdminStub) ListAdminMessages(ctx context.Context, actor entity.Actor) ([]entity.AdminMessage, error) {
	return nil, nil
}

func (s *adminContractAdminStub) RepairAggregates(context.Context, entity.Actor) error {
	return nil
}

func (s *adminContractAdminStub) ListProducts(context.Context, entity.Actor, entity.Pagination) ([]entity.Product, int, error) {
	return nil, 0, nil
}

func (s *adminContractAdminStub) GetProduct(context.Context, entity.Actor, string) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *adminContractAdminStub) UpsertProduct(context.Context, entity.Actor, entity.Product) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *adminContractAdminStub) DeleteProduct(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminContractAdminStub) ListCategories(context.Context, entity.Actor) ([]entity.Category, error) {
	return nil, nil
}

func (s *adminContractAdminStub) UpsertCategory(context.Context, entity.Actor, entity.Category) (entity.Category, error) {
	return entity.Category{}, nil
}

func (s *adminContractAdminStub) DeleteCategory(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminContractAdminStub) ListSettings(ctx context.Context, actor entity.Actor) ([]entity.Setting, error) {
	if s.listSettingsFunc != nil {
		return s.listSettingsFunc(ctx, actor)
	}
	return nil, nil
}

func (s *adminContractAdminStub) SetSetting(ctx context.Context, actor entity.Actor, key, value string) error {
	if s.setSettingFunc != nil {
		return s.setSettingFunc(ctx, actor, key, value)
	}
	return nil
}

func (s *adminContractAdminStub) DeleteSetting(ctx context.Context, actor entity.Actor, key string) error {
	if s.deleteSettingFunc != nil {
		return s.deleteSettingFunc(ctx, actor, key)
	}
	return nil
}

func (s *adminContractAdminStub) SendBroadcast(ctx context.Context, actor entity.Actor, title, body string) error {
	if s.broadcastFunc != nil {
		return s.broadcastFunc(ctx, actor, title, body)
	}
	return nil
}

func (s *adminContractAdminStub) SendTargeted(ctx context.Context, actor entity.Actor, userID, title, body string) error {
	if s.targetedFunc != nil {
		return s.targetedFunc(ctx, actor, userID, title, body)
	}
	return nil
}

func (s *adminContractAdminStub) ListReviews(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error) {
	if s.listReviewsFunc != nil {
		return s.listReviewsFunc(ctx, actor, page, query, status)
	}
	return nil, repo.ReviewModerationStats{}, 0, nil
}

func (s *adminContractAdminStub) UpdateReviewStatus(ctx context.Context, actor entity.Actor, reviewID int64, status entity.ReviewStatus) (entity.Review, error) {
	if s.updateReviewFunc != nil {
		return s.updateReviewFunc(ctx, actor, reviewID, status)
	}
	return entity.Review{}, nil
}

func (s *adminContractAdminStub) BulkPublishReviews(ctx context.Context, actor entity.Actor, reviewIDs []int64) (int, error) {
	if s.bulkReviewFunc != nil {
		return s.bulkReviewFunc(ctx, actor, reviewIDs)
	}
	return len(reviewIDs), nil
}

func (s *adminContractAdminStub) DeleteReview(ctx context.Context, actor entity.Actor, reviewID int64) error {
	if s.deleteReviewFunc != nil {
		return s.deleteReviewFunc(ctx, actor, reviewID)
	}
	return nil
}

type adminContractImporterStub struct {
	importFunc func(ctx context.Context, items []entity.ImportItem) (entity.ImportResult, error)
}

func (s *adminContractImporterStub) Import(ctx context.Context, items []entity.ImportItem) (entity.ImportResult, error) {
	if s.importFunc != nil {
		return s.importFunc(ctx, items)
	}
	return entity.ImportResult{}, nil
}

func TestUS4_AdminContract_TDD(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin-id", "admin", true)
	require.NoError(t, err)
	userToken, err := jwtManager.GenerateTokenWithProfile("user-id", "user", false)
	require.NoError(t, err)

	adminUC := &adminContractAdminStub{
		listOrdersFunc: func(_ context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Order, int, error) {
			require.True(t, actor.IsAdmin)
			require.Equal(t, "search", query)
			require.Equal(t, "pending", status)
			require.Equal(t, 0, page.Offset)
			return []entity.Order{{ID: "o1", ProductName: "Product", Status: entity.OrderStatusPending}}, 1, nil
		},
		getOrderFunc: func(_ context.Context, actor entity.Actor, orderID string) (entity.Order, error) {
			require.True(t, actor.IsAdmin)
			require.Equal(t, "o1", orderID)
			return entity.Order{ID: "o1", ProductName: "Product", Status: entity.OrderStatusDelivered}, nil
		},
		updateOrderStatusFunc: func(_ context.Context, actor entity.Actor, orderID string, status entity.OrderStatus) error {
			require.True(t, actor.IsAdmin)
			require.Equal(t, "o1", orderID)
			require.Equal(t, entity.OrderStatusCancelled, status)
			return nil
		},
		deleteOrderFunc: func(_ context.Context, actor entity.Actor, orderID string) error {
			require.True(t, actor.IsAdmin)
			require.Equal(t, "o1", orderID)
			return nil
		},
		listRefundsFunc: func(_ context.Context, actor entity.Actor, status string) ([]entity.RefundRequest, error) {
			require.True(t, actor.IsAdmin)
			require.Equal(t, "pending", status)
			return []entity.RefundRequest{{ID: 9, OrderID: "o1", Status: entity.RefundStatusPending}}, nil
		},
		decideRefundFunc: func(_ context.Context, actor entity.Actor, refundID int64, approve bool, note string) (entity.RefundRequest, error) {
			require.True(t, actor.IsAdmin)
			require.Equal(t, int64(9), refundID)
			status := entity.RefundStatusRejected
			if approve {
				status = entity.RefundStatusApproved
				require.Equal(t, "ok", note)
			} else {
				require.Equal(t, "no", note)
			}
			return entity.RefundRequest{ID: refundID, Status: status, AdminNote: note}, nil
		},

		listSettingsFunc: func(_ context.Context, actor entity.Actor) ([]entity.Setting, error) {
			require.True(t, actor.IsAdmin)
			return []entity.Setting{{Key: "shopName", Value: "Grip Store"}}, nil
		},
		setSettingFunc: func(_ context.Context, actor entity.Actor, key, value string) error {
			require.True(t, actor.IsAdmin)
			require.Equal(t, "shopName", key)
			require.Equal(t, "Updated Store", value)
			return nil
		},
		deleteSettingFunc: func(_ context.Context, actor entity.Actor, key string) error {
			require.True(t, actor.IsAdmin)
			require.Equal(t, "shopName", key)
			return nil
		},
		broadcastFunc: func(_ context.Context, actor entity.Actor, title, body string) error {
			require.True(t, actor.IsAdmin)
			require.Equal(t, "Broadcast", title)
			require.Equal(t, "Hello", body)
			return nil
		},
		targetedFunc: func(_ context.Context, actor entity.Actor, userID, title, body string) error {
			require.True(t, actor.IsAdmin)
			require.Equal(t, "u1", userID)
			require.Equal(t, "Direct", title)
			require.Equal(t, "Hi", body)
			return nil
		},
	}
	importerUC := &adminContractImporterStub{
		importFunc: func(_ context.Context, items []entity.ImportItem) (entity.ImportResult, error) {
			require.Len(t, items, 1)
			require.Equal(t, entity.ImportItemProduct, items[0].Type)
			return entity.ImportResult{Imported: 1}, nil
		},
	}

	app := fiber.New()
	v1.NewRoutes(
		app.Group("/v1"),
		nil, nil, nil, nil, nil, nil, nil, nil,
		adminUC,
		nil, nil, nil, nil, nil, nil, nil, importerUC,
		jwtManager,
		"admin",
		logger.New("error"),
	)

	t.Run("auth policy", func(t *testing.T) {
		require.Equal(t, http.StatusUnauthorized, adminContractRequest(t, app, http.MethodGet, "/v1/admin/settings", nil, "").StatusCode)
		require.Equal(t, http.StatusForbidden, adminContractRequest(t, app, http.MethodGet, "/v1/admin/settings", nil, "Bearer "+userToken).StatusCode)
	})

	t.Run("settings routes", func(t *testing.T) {
		resp := adminContractRequest(t, app, http.MethodGet, "/v1/admin/settings", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body adminContractEnvelope
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		require.NotEmpty(t, body.Data)

		resp = adminContractRequest(t, app, http.MethodPut, "/v1/admin/settings/shopName", []byte(`{"value":"Updated Store"}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = adminContractRequest(t, app, http.MethodDelete, "/v1/admin/settings/shopName", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})


	t.Run("orders and refunds routes", func(t *testing.T) {
		resp := adminContractRequest(t, app, http.MethodGet, "/v1/admin/orders?q=search&status=pending", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = adminContractRequest(t, app, http.MethodGet, "/v1/admin/orders/o1", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = adminContractRequest(t, app, http.MethodPatch, "/v1/admin/orders/o1", []byte(`{"status":"cancelled"}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp = adminContractRequest(t, app, http.MethodDelete, "/v1/admin/orders/o1", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp = adminContractRequest(t, app, http.MethodGet, "/v1/admin/refunds?status=pending", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = adminContractRequest(t, app, http.MethodPost, "/v1/admin/refunds/9/approve", []byte(`{"note":"ok"}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = adminContractRequest(t, app, http.MethodPost, "/v1/admin/refunds/9/reject", []byte(`{"note":"no"}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("messages and ops routes", func(t *testing.T) {
		resp := adminContractRequest(t, app, http.MethodPost, "/v1/admin/messages/broadcast", []byte(`{"title":"Broadcast","body":"Hello"}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp = adminContractRequest(t, app, http.MethodPost, "/v1/admin/messages/targeted", []byte(`{"userId":"u1","title":"Direct","body":"Hi"}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp = adminContractRequest(t, app, http.MethodPost, "/v1/admin/notifications/test", []byte(`{"channel":"email","to":"user@example.com"}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = adminContractRequest(t, app, http.MethodPost, "/v1/admin/data/import", []byte(`{"items":[{"type":"product","data":{"title":"P","sku":"S"}}]}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func adminContractRequest(t *testing.T, app *fiber.App, method, target string, body []byte, authHeader string) *http.Response {
	t.Helper()

	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := app.Test(req)
	require.NoError(t, err)
	return resp
}
