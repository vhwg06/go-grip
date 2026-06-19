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

type adminCardsImportUseCaseStub struct {
	listCardsFunc   func(ctx context.Context, actor entity.Actor, productID string) ([]entity.Card, error)
	createCardFunc  func(ctx context.Context, actor entity.Actor, productID, cardKey string) (entity.Card, error)
	deleteCardFunc  func(ctx context.Context, actor entity.Actor, cardID int64) error
	importCardsFunc func(ctx context.Context, actor entity.Actor, productID string, keys []string) (int, error)
}

func (s *adminCardsImportUseCaseStub) ListUsers(context.Context, entity.Actor, entity.Pagination) ([]entity.User, int, error) {
	return nil, 0, nil
}

func (s *adminCardsImportUseCaseStub) UpdateUserStatus(context.Context, entity.Actor, string, entity.UserStatus) error {
	return nil
}

func (s *adminCardsImportUseCaseStub) UpdateUserPoints(context.Context, entity.Actor, string, int) error {
	return nil
}

func (s *adminCardsImportUseCaseStub) ListOrders(context.Context, entity.Actor, entity.Pagination, string, string) ([]entity.Order, int, error) {
	return nil, 0, nil
}

func (s *adminCardsImportUseCaseStub) GetOrder(context.Context, entity.Actor, string) (entity.Order, error) {
	return entity.Order{}, nil
}

func (s *adminCardsImportUseCaseStub) UpdateOrderStatus(context.Context, entity.Actor, string, entity.OrderStatus) error {
	return nil
}

func (s *adminCardsImportUseCaseStub) DeleteOrder(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminCardsImportUseCaseStub) ListRefunds(context.Context, entity.Actor, string) ([]entity.RefundRequest, error) {
	return nil, nil
}

func (s *adminCardsImportUseCaseStub) DecideRefund(context.Context, entity.Actor, int64, bool, string) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}

func (s *adminCardsImportUseCaseStub) ListReviews(context.Context, entity.Actor, entity.Pagination, string, string) ([]entity.Review, repo.ReviewModerationStats, int, error) {
	return nil, repo.ReviewModerationStats{}, 0, nil
}

func (s *adminCardsImportUseCaseStub) UpdateReviewStatus(context.Context, entity.Actor, int64, entity.ReviewStatus) (entity.Review, error) {
	return entity.Review{}, nil
}

func (s *adminCardsImportUseCaseStub) BulkPublishReviews(context.Context, entity.Actor, []int64) (int, error) {
	return 0, nil
}

func (s *adminCardsImportUseCaseStub) DeleteReview(context.Context, entity.Actor, int64) error {
	return nil
}

func (s *adminCardsImportUseCaseStub) ListCards(ctx context.Context, actor entity.Actor, productID string) ([]entity.Card, error) {
	if s.listCardsFunc != nil {
		return s.listCardsFunc(ctx, actor, productID)
	}
	return nil, nil
}

func (s *adminCardsImportUseCaseStub) CreateCard(ctx context.Context, actor entity.Actor, productID, cardKey string) (entity.Card, error) {
	if s.createCardFunc != nil {
		return s.createCardFunc(ctx, actor, productID, cardKey)
	}
	return entity.Card{}, nil
}

func (s *adminCardsImportUseCaseStub) DeleteCard(ctx context.Context, actor entity.Actor, cardID int64) error {
	if s.deleteCardFunc != nil {
		return s.deleteCardFunc(ctx, actor, cardID)
	}
	return nil
}

func (s *adminCardsImportUseCaseStub) ImportCards(ctx context.Context, actor entity.Actor, productID string, keys []string) (int, error) {
	if s.importCardsFunc != nil {
		return s.importCardsFunc(ctx, actor, productID, keys)
	}
	return 0, nil
}

func (s *adminCardsImportUseCaseStub) RepairAggregates(context.Context, entity.Actor) error {
	return nil
}

func (s *adminCardsImportUseCaseStub) ListProducts(context.Context, entity.Actor, entity.Pagination) ([]entity.Product, int, error) {
	return nil, 0, nil
}

func (s *adminCardsImportUseCaseStub) GetProduct(context.Context, entity.Actor, string) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *adminCardsImportUseCaseStub) UpsertProduct(context.Context, entity.Actor, entity.Product) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *adminCardsImportUseCaseStub) DeleteProduct(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminCardsImportUseCaseStub) ListCategories(context.Context, entity.Actor) ([]entity.Category, error) {
	return nil, nil
}

func (s *adminCardsImportUseCaseStub) UpsertCategory(context.Context, entity.Actor, entity.Category) (entity.Category, error) {
	return entity.Category{}, nil
}

func (s *adminCardsImportUseCaseStub) DeleteCategory(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminCardsImportUseCaseStub) ListSettings(context.Context, entity.Actor) ([]entity.Setting, error) {
	return nil, nil
}

func (s *adminCardsImportUseCaseStub) SetSetting(context.Context, entity.Actor, string, string) error {
	return nil
}

func (s *adminCardsImportUseCaseStub) DeleteSetting(context.Context, entity.Actor, string) error {
	return nil
}

func (s *adminCardsImportUseCaseStub) SendBroadcast(context.Context, entity.Actor, string, string) error {
	return nil
}

func (s *adminCardsImportUseCaseStub) SendTargeted(context.Context, entity.Actor, string, string, string) error {
	return nil
}

type importerStub struct {
	importFunc func(ctx context.Context, items []entity.ImportItem) (entity.ImportResult, error)
}

func (s *importerStub) Import(ctx context.Context, items []entity.ImportItem) (entity.ImportResult, error) {
	if s.importFunc != nil {
		return s.importFunc(ctx, items)
	}
	return entity.ImportResult{}, nil
}

func TestAdminCardsAndImportEndpoints(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin-id", "admin", true)
	require.NoError(t, err)

	setupApp := func(adminUC *adminCardsImportUseCaseStub, importerUC *importerStub) *fiber.App {
		v := &V1{
			adminUC:    adminUC,
			importer:   importerUC,
			jwtManager: jwtManager,
			adminUsers: "admin",
		}
		app := fiber.New()
		v.registerGripStoreRoutes(app.Group("/v1"))
		return app
	}

	t.Run("card routes perform CRUD and import", func(t *testing.T) {
		var deletedCardID int64
		var createdCardKey string
		var importedProductID string
		app := setupApp(&adminCardsImportUseCaseStub{
			listCardsFunc: func(_ context.Context, actor entity.Actor, productID string) ([]entity.Card, error) {
				require.True(t, actor.IsAdmin)
				require.Equal(t, "p1", productID)
				return []entity.Card{{ID: 1, ProductID: "p1", CardKey: "k1"}}, nil
			},
			createCardFunc: func(_ context.Context, actor entity.Actor, productID, cardKey string) (entity.Card, error) {
				require.True(t, actor.IsAdmin)
				createdCardKey = cardKey
				return entity.Card{ID: 2, ProductID: productID, CardKey: cardKey}, nil
			},
			deleteCardFunc: func(_ context.Context, actor entity.Actor, cardID int64) error {
				require.True(t, actor.IsAdmin)
				deletedCardID = cardID
				return nil
			},
			importCardsFunc: func(_ context.Context, actor entity.Actor, productID string, keys []string) (int, error) {
				require.True(t, actor.IsAdmin)
				importedProductID = productID
				require.Equal(t, []string{"a", "b"}, keys)
				return 2, nil
			},
		}, &importerStub{})

		resp := testRequest(t, app, http.MethodGet, "/v1/admin/cards?productId=p1", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = testRequest(t, app, http.MethodPost, "/v1/admin/cards", []byte(`{"productId":"p1","cardKey":"new-key"}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		require.Equal(t, "new-key", createdCardKey)

		resp = testRequest(t, app, http.MethodDelete, "/v1/admin/cards/7", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		require.Equal(t, int64(7), deletedCardID)

		resp = testRequest(t, app, http.MethodPost, "/v1/admin/cards/import", []byte(`{"productId":"p1","keys":["a","b"]}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "p1", importedProductID)

		resp = testRequest(t, app, http.MethodPost, "/v1/admin/cards/replenish", []byte(`{"productId":"p1","keys":["a","b"]}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("notification test and data import return success envelopes", func(t *testing.T) {
		app := setupApp(&adminCardsImportUseCaseStub{}, &importerStub{
			importFunc: func(_ context.Context, items []entity.ImportItem) (entity.ImportResult, error) {
				require.Len(t, items, 1)
				return entity.ImportResult{Imported: 1}, nil
			},
		})

		resp := testRequest(t, app, http.MethodPost, "/v1/admin/notifications/test", []byte(`{"channel":"email","to":"user@example.com"}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = testRequest(t, app, http.MethodPost, "/v1/admin/data/import", []byte(`{"items":[{"type":"product","data":{"title":"P","sku":"S"}}]}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body envelope
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		raw, err := json.Marshal(body.Data)
		require.NoError(t, err)
		var result entity.ImportResult
		require.NoError(t, json.Unmarshal(raw, &result))
		require.Equal(t, 1, result.Imported)
	})
}
