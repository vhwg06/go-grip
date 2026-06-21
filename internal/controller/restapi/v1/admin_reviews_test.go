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

type adminReviewsUseCaseStub struct {
	listFunc   func(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error)
	updateFunc func(ctx context.Context, actor entity.Actor, reviewID int64, status entity.ReviewStatus) (entity.Review, error)
	bulkFunc   func(ctx context.Context, actor entity.Actor, reviewIDs []int64) (int, error)
	deleteFunc func(ctx context.Context, actor entity.Actor, reviewID int64) error
}

func (s *adminReviewsUseCaseStub) ListUsers(context.Context, entity.Actor, entity.Pagination) ([]entity.User, int, error) {
	return nil, 0, nil
}
func (s *adminReviewsUseCaseStub) UpdateUserStatus(context.Context, entity.Actor, string, entity.UserStatus) error {
	return nil
}
func (s *adminReviewsUseCaseStub) UpdateUserPoints(context.Context, entity.Actor, string, int) error {
	return nil
}
func (s *adminReviewsUseCaseStub) ListProducts(context.Context, entity.Actor, entity.Pagination) ([]entity.Product, int, error) {
	return nil, 0, nil
}
func (s *adminReviewsUseCaseStub) GetProduct(context.Context, entity.Actor, string) (entity.Product, error) {
	return entity.Product{}, nil
}
func (s *adminReviewsUseCaseStub) UpsertProduct(context.Context, entity.Actor, entity.Product) (entity.Product, error) {
	return entity.Product{}, nil
}
func (s *adminReviewsUseCaseStub) DeleteProduct(context.Context, entity.Actor, string) error {
	return nil
}
func (s *adminReviewsUseCaseStub) ListCategories(context.Context, entity.Actor) ([]entity.Category, error) {
	return nil, nil
}
func (s *adminReviewsUseCaseStub) UpsertCategory(context.Context, entity.Actor, entity.Category) (entity.Category, error) {
	return entity.Category{}, nil
}
func (s *adminReviewsUseCaseStub) DeleteCategory(context.Context, entity.Actor, string) error {
	return nil
}
func (s *adminReviewsUseCaseStub) ListOrders(context.Context, entity.Actor, entity.Pagination, string, string) ([]entity.Order, int, error) {
	return nil, 0, nil
}
func (s *adminReviewsUseCaseStub) GetOrder(context.Context, entity.Actor, string) (entity.Order, error) {
	return entity.Order{}, nil
}
func (s *adminReviewsUseCaseStub) UpdateOrderStatus(context.Context, entity.Actor, string, entity.OrderStatus) error {
	return nil
}
func (s *adminReviewsUseCaseStub) DeleteOrder(context.Context, entity.Actor, string) error {
	return nil
}
func (s *adminReviewsUseCaseStub) ListRefunds(context.Context, entity.Actor, string) ([]entity.RefundRequest, error) {
	return nil, nil
}
func (s *adminReviewsUseCaseStub) DecideRefund(context.Context, entity.Actor, int64, bool, string) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}
func (s *adminReviewsUseCaseStub) GetRefund(context.Context, entity.Actor, int64) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}
func (s *adminReviewsUseCaseStub) GetOrderRefundStatus(context.Context, entity.Actor, string) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}
func (s *adminReviewsUseCaseStub) ListCards(context.Context, entity.Actor) ([]entity.Card, error) {
	return nil, nil
}
func (s *adminReviewsUseCaseStub) ListSettings(context.Context, entity.Actor) ([]entity.Setting, error) {
	return nil, nil
}
func (s *adminReviewsUseCaseStub) SetSetting(context.Context, entity.Actor, string, string) error {
	return nil
}
func (s *adminReviewsUseCaseStub) DeleteSetting(context.Context, entity.Actor, string) error {
	return nil
}
func (s *adminReviewsUseCaseStub) SendBroadcast(context.Context, entity.Actor, string, string) error {
	return nil
}
func (s *adminReviewsUseCaseStub) SendTargeted(context.Context, entity.Actor, string, string, string) error {
	return nil
}
func (s *adminReviewsUseCaseStub) ListReviews(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error) {
	if s.listFunc != nil {
		return s.listFunc(ctx, actor, page, query, status)
	}
	return nil, repo.ReviewModerationStats{}, 0, nil
}
func (s *adminReviewsUseCaseStub) UpdateReviewStatus(ctx context.Context, actor entity.Actor, reviewID int64, status entity.ReviewStatus) (entity.Review, error) {
	if s.updateFunc != nil {
		return s.updateFunc(ctx, actor, reviewID, status)
	}
	return entity.Review{}, nil
}
func (s *adminReviewsUseCaseStub) BulkPublishReviews(ctx context.Context, actor entity.Actor, reviewIDs []int64) (int, error) {
	if s.bulkFunc != nil {
		return s.bulkFunc(ctx, actor, reviewIDs)
	}
	return len(reviewIDs), nil
}
func (s *adminReviewsUseCaseStub) DeleteReview(ctx context.Context, actor entity.Actor, reviewID int64) error {
	if s.deleteFunc != nil {
		return s.deleteFunc(ctx, actor, reviewID)
	}
	return nil
}
func (s *adminReviewsUseCaseStub) RepairAggregates(context.Context, entity.Actor) error { return nil }

func TestAdminReviewModerationEndpoints(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin-id", "admin", true)
	require.NoError(t, err)
	userToken, err := jwtManager.GenerateTokenWithProfile("user-id", "user", false)
	require.NoError(t, err)

	setupApp := func(uc *adminReviewsUseCaseStub) *fiber.App {
		v := &V1{
			adminUC:    uc,
			jwtManager: jwtManager,
			adminUsers: "admin",
		}
		app := fiber.New()
		v.registerGripStoreRoutes(app.Group("/v1"))
		return app
	}

	t.Run("list returns queue and stats with admin auth", func(t *testing.T) {
		app := setupApp(&adminReviewsUseCaseStub{
			listFunc: func(_ context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error) {
				require.True(t, actor.IsAdmin)
				require.Equal(t, "PENDING", status)
				require.Equal(t, 0, page.Normalize().Offset)
				return []entity.Review{{ID: 101, ProductID: "p1", ProductName: "Grip", OrderID: "o1", UserID: "u1", Username: "alice", Rating: 5, Comment: "Great", Status: entity.ReviewStatusPending, CreatedAt: time.Unix(1700000000, 0).UTC()}}, repo.ReviewModerationStats{Pending: 1, Featured: 2, Hidden: 3}, 1, nil
			},
		})

		require.Equal(t, http.StatusUnauthorized, testRequest(t, app, http.MethodGet, "/v1/admin/reviews?status=PENDING", nil, "").StatusCode)
		require.Equal(t, http.StatusForbidden, testRequest(t, app, http.MethodGet, "/v1/admin/reviews?status=PENDING", nil, "Bearer "+userToken).StatusCode)

		resp := testRequest(t, app, http.MethodGet, "/v1/admin/reviews?status=PENDING", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body envelope
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		raw, err := json.Marshal(body.Data)
		require.NoError(t, err)
		var payload struct {
			Reviews []map[string]any           `json:"reviews"`
			Stats   repo.ReviewModerationStats `json:"stats"`
			Total   int                        `json:"total"`
		}
		require.NoError(t, json.Unmarshal(raw, &payload))
		require.Len(t, payload.Reviews, 1)
		require.Equal(t, float64(101), payload.Reviews[0]["id"])
		require.Equal(t, 1, payload.Stats.Pending)
		require.Equal(t, 1, payload.Total)
	})

	t.Run("approve hide feature bulk and delete route mutations", func(t *testing.T) {
		var bulkIDs []int64
		var deletedID int64
		app := setupApp(&adminReviewsUseCaseStub{
			updateFunc: func(_ context.Context, actor entity.Actor, reviewID int64, status entity.ReviewStatus) (entity.Review, error) {
				require.True(t, actor.IsAdmin)
				return entity.Review{ID: reviewID, Status: status}, nil
			},
			bulkFunc: func(_ context.Context, actor entity.Actor, reviewIDs []int64) (int, error) {
				require.True(t, actor.IsAdmin)
				bulkIDs = reviewIDs
				return len(reviewIDs), nil
			},
			deleteFunc: func(_ context.Context, actor entity.Actor, reviewID int64) error {
				require.True(t, actor.IsAdmin)
				deletedID = reviewID
				return nil
			},
		})

		require.Equal(t, http.StatusOK, testRequest(t, app, http.MethodPut, "/v1/admin/reviews/101/approve", []byte(`{}`), "Bearer "+adminToken).StatusCode)
		require.Equal(t, http.StatusOK, testRequest(t, app, http.MethodPut, "/v1/admin/reviews/101/hide", []byte(`{}`), "Bearer "+adminToken).StatusCode)
		require.Equal(t, http.StatusBadRequest, testRequest(t, app, http.MethodPut, "/v1/admin/reviews/101/feature", []byte(`{}`), "Bearer "+adminToken).StatusCode)
		require.Equal(t, http.StatusOK, testRequest(t, app, http.MethodPut, "/v1/admin/reviews/101/feature", []byte(`{"isFeatured":true}`), "Bearer "+adminToken).StatusCode)
		require.Equal(t, http.StatusOK, testRequest(t, app, http.MethodPost, "/v1/admin/reviews/publish-selected", []byte(`{"ids":[101,102]}`), "Bearer "+adminToken).StatusCode)
		require.Equal(t, []int64{101, 102}, bulkIDs)
		require.Equal(t, http.StatusOK, testRequest(t, app, http.MethodDelete, "/v1/admin/reviews/101", nil, "Bearer "+adminToken).StatusCode)
		require.Equal(t, int64(101), deletedID)
	})
}
