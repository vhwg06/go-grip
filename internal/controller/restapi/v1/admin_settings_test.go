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

type adminSettingsUseCaseStub struct {
	listSettingsFunc  func(ctx context.Context, actor entity.Actor) ([]entity.Setting, error)
	setSettingFunc    func(ctx context.Context, actor entity.Actor, key, value string) error
	deleteSettingFunc func(ctx context.Context, actor entity.Actor, key string) error
}

func (s *adminSettingsUseCaseStub) ListProducts(context.Context, entity.Actor, entity.Pagination) ([]entity.Product, int, error) {
	return nil, 0, nil
}

func (s *adminSettingsUseCaseStub) GetProduct(context.Context, entity.Actor, string) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *adminSettingsUseCaseStub) UpsertProduct(context.Context, entity.Actor, entity.Product) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *adminSettingsUseCaseStub) DeleteProduct(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminSettingsUseCaseStub) ListCategories(context.Context, entity.Actor) ([]entity.Category, error) {
	return nil, nil
}

func (s *adminSettingsUseCaseStub) UpsertCategory(context.Context, entity.Actor, entity.Category) (entity.Category, error) {
	return entity.Category{}, nil
}

func (s *adminSettingsUseCaseStub) DeleteCategory(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminSettingsUseCaseStub) ListUsers(context.Context, entity.Actor, entity.Pagination) ([]entity.User, int, error) {
	return nil, 0, nil
}

func (s *adminSettingsUseCaseStub) UpdateUserStatus(context.Context, entity.Actor, string, entity.UserStatus) error {
	return nil
}

func (s *adminSettingsUseCaseStub) UpdateUserPoints(context.Context, entity.Actor, string, int) error {
	return nil
}

func (s *adminSettingsUseCaseStub) ListOrders(context.Context, entity.Actor, entity.Pagination, string, string) ([]entity.Order, int, error) {
	return nil, 0, nil
}

func (s *adminSettingsUseCaseStub) GetOrder(context.Context, entity.Actor, string) (entity.Order, error) {
	return entity.Order{}, nil
}

func (s *adminSettingsUseCaseStub) UpdateOrderStatus(context.Context, entity.Actor, string, entity.OrderStatus) error {
	return nil
}

func (s *adminSettingsUseCaseStub) DeleteOrder(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminSettingsUseCaseStub) ListRefunds(context.Context, entity.Actor, string) ([]entity.RefundRequest, error) {
	return nil, nil
}

func (s *adminSettingsUseCaseStub) DecideRefund(context.Context, entity.Actor, int64, bool, string) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}

func (s *adminSettingsUseCaseStub) ListReviews(context.Context, entity.Actor, entity.Pagination, string, string) ([]entity.Review, repo.ReviewModerationStats, int, error) {
	return nil, repo.ReviewModerationStats{}, 0, nil
}

func (s *adminSettingsUseCaseStub) UpdateReviewStatus(context.Context, entity.Actor, int64, entity.ReviewStatus) (entity.Review, error) {
	return entity.Review{}, nil
}

func (s *adminSettingsUseCaseStub) BulkPublishReviews(context.Context, entity.Actor, []int64) (int, error) {
	return 0, nil
}

func (s *adminSettingsUseCaseStub) DeleteReview(context.Context, entity.Actor, int64) error {
	return nil
}

func (s *adminSettingsUseCaseStub) RepairAggregates(context.Context, entity.Actor) error {
	return nil
}

func (s *adminSettingsUseCaseStub) ListSettings(ctx context.Context, actor entity.Actor) ([]entity.Setting, error) {
	if s.listSettingsFunc != nil {
		return s.listSettingsFunc(ctx, actor)
	}
	return nil, nil
}

func (s *adminSettingsUseCaseStub) SetSetting(ctx context.Context, actor entity.Actor, key, value string) error {
	if s.setSettingFunc != nil {
		return s.setSettingFunc(ctx, actor, key, value)
	}
	return nil
}

func (s *adminSettingsUseCaseStub) DeleteSetting(ctx context.Context, actor entity.Actor, key string) error {
	if s.deleteSettingFunc != nil {
		return s.deleteSettingFunc(ctx, actor, key)
	}
	return nil
}

func (s *adminSettingsUseCaseStub) SendBroadcast(context.Context, entity.Actor, string, string) error {
	return nil
}

func (s *adminSettingsUseCaseStub) SendTargeted(context.Context, entity.Actor, string, string, string) error {
	return nil
}

func TestAdminSettingsEndpoints(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin-id", "admin", true)
	require.NoError(t, err)
	userToken, err := jwtManager.GenerateTokenWithProfile("user-id", "user", false)
	require.NoError(t, err)

	setupApp := func(uc *adminSettingsUseCaseStub) *fiber.App {
		v := &V1{
			adminUC:    uc,
			jwtManager: jwtManager,
			adminUsers: "admin",
		}
		app := fiber.New()
		v.registerGripStoreRoutes(app.Group("/v1"))
		return app
	}

	t.Run("list requires admin and returns settings", func(t *testing.T) {
		app := setupApp(&adminSettingsUseCaseStub{
			listSettingsFunc: func(_ context.Context, actor entity.Actor) ([]entity.Setting, error) {
				require.True(t, actor.IsAdmin)
				return []entity.Setting{{Key: "shopName", Value: "Grip Store"}}, nil
			},
		})

		require.Equal(t, http.StatusUnauthorized, testRequest(t, app, http.MethodGet, "/v1/admin/settings", nil, "").StatusCode)
		require.Equal(t, http.StatusForbidden, testRequest(t, app, http.MethodGet, "/v1/admin/settings", nil, "Bearer "+userToken).StatusCode)

		resp := testRequest(t, app, http.MethodGet, "/v1/admin/settings", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body envelope
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		raw, err := json.Marshal(body.Data)
		require.NoError(t, err)

		var settings []entity.Setting
		require.NoError(t, json.Unmarshal(raw, &settings))
		require.Len(t, settings, 1)
		require.Equal(t, "shopName", settings[0].Key)
	})

	t.Run("put validates body and stores setting", func(t *testing.T) {
		var gotKey, gotValue string
		app := setupApp(&adminSettingsUseCaseStub{
			setSettingFunc: func(_ context.Context, actor entity.Actor, key, value string) error {
				require.True(t, actor.IsAdmin)
				gotKey = key
				gotValue = value
				return nil
			},
		})

		require.Equal(t, http.StatusBadRequest, testRequest(t, app, http.MethodPut, "/v1/admin/settings/shopName", []byte(`{invalid}`), "Bearer "+adminToken).StatusCode)

		resp := testRequest(t, app, http.MethodPut, "/v1/admin/settings/shopName", []byte(`{"value":"Updated Store"}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "shopName", gotKey)
		require.Equal(t, "Updated Store", gotValue)
	})

	t.Run("delete removes setting", func(t *testing.T) {
		var deletedKey string
		app := setupApp(&adminSettingsUseCaseStub{
			deleteSettingFunc: func(_ context.Context, actor entity.Actor, key string) error {
				require.True(t, actor.IsAdmin)
				deletedKey = key
				return nil
			},
		})

		resp := testRequest(t, app, http.MethodDelete, "/v1/admin/settings/announcement", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		require.Equal(t, "announcement", deletedKey)
	})
}
