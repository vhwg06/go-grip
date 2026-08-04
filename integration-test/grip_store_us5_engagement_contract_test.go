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
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type engagementEnvelope struct {
	Data  json.RawMessage `json:"data"`
	Error string          `json:"error"`
}

type engagementWishlistStub struct {
	listFunc         func(ctx context.Context, page entity.Pagination) ([]entity.WishlistItem, int, error)
	createFunc       func(ctx context.Context, actor entity.Actor, title, description string) (entity.WishlistItem, error)
	updateFunc       func(ctx context.Context, actor entity.Actor, item entity.WishlistItem) (entity.WishlistItem, error)
	deleteFunc       func(ctx context.Context, actor entity.Actor, itemID int64) error
	toggleVoteFunc   func(ctx context.Context, actor entity.Actor, itemID int64) error
	createReviewFunc func(ctx context.Context, actor entity.Actor, review entity.Review) (entity.Review, error)
	listReviewsFunc  func(ctx context.Context, productID string) ([]entity.Review, error)
}

func (s *engagementWishlistStub) List(ctx context.Context, page entity.Pagination) ([]entity.WishlistItem, int, error) {
	if s.listFunc != nil {
		return s.listFunc(ctx, page)
	}
	return nil, 0, nil
}

func (s *engagementWishlistStub) Create(ctx context.Context, actor entity.Actor, title, description string) (entity.WishlistItem, error) {
	if s.createFunc != nil {
		return s.createFunc(ctx, actor, title, description)
	}
	return entity.WishlistItem{}, nil
}

func (s *engagementWishlistStub) Update(ctx context.Context, actor entity.Actor, item entity.WishlistItem) (entity.WishlistItem, error) {
	if s.updateFunc != nil {
		return s.updateFunc(ctx, actor, item)
	}
	return entity.WishlistItem{}, nil
}

func (s *engagementWishlistStub) Delete(ctx context.Context, actor entity.Actor, itemID int64) error {
	if s.deleteFunc != nil {
		return s.deleteFunc(ctx, actor, itemID)
	}
	return nil
}

func (s *engagementWishlistStub) ToggleVote(ctx context.Context, actor entity.Actor, itemID int64) error {
	if s.toggleVoteFunc != nil {
		return s.toggleVoteFunc(ctx, actor, itemID)
	}
	return nil
}

func (s *engagementWishlistStub) CreateReview(ctx context.Context, actor entity.Actor, review entity.Review) (entity.Review, error) {
	if s.createReviewFunc != nil {
		return s.createReviewFunc(ctx, actor, review)
	}
	return entity.Review{}, nil
}

func (s *engagementWishlistStub) ListReviews(ctx context.Context, productID string) ([]entity.Review, error) {
	if s.listReviewsFunc != nil {
		return s.listReviewsFunc(ctx, productID)
	}
	return nil, nil
}

type engagementNotifyStub struct {
	inboxFunc       func(ctx context.Context, actor entity.Actor, page entity.Pagination) ([]entity.UserNotification, int, error)
	unreadCountFunc func(ctx context.Context, actor entity.Actor) (int, error)
	markReadFunc    func(ctx context.Context, actor entity.Actor, notificationID int64) error
	markAllFunc     func(ctx context.Context, actor entity.Actor) error
	clearFunc       func(ctx context.Context, actor entity.Actor) error
}

func (s *engagementNotifyStub) Inbox(ctx context.Context, actor entity.Actor, page entity.Pagination) ([]entity.UserNotification, int, error) {
	if s.inboxFunc != nil {
		return s.inboxFunc(ctx, actor, page)
	}
	return nil, 0, nil
}

func (s *engagementNotifyStub) UnreadCount(ctx context.Context, actor entity.Actor) (int, error) {
	if s.unreadCountFunc != nil {
		return s.unreadCountFunc(ctx, actor)
	}
	return 0, nil
}

func (s *engagementNotifyStub) MarkRead(ctx context.Context, actor entity.Actor, notificationID int64) error {
	if s.markReadFunc != nil {
		return s.markReadFunc(ctx, actor, notificationID)
	}
	return nil
}

func (s *engagementNotifyStub) MarkAllRead(ctx context.Context, actor entity.Actor) error {
	if s.markAllFunc != nil {
		return s.markAllFunc(ctx, actor)
	}
	return nil
}

func (s *engagementNotifyStub) Clear(ctx context.Context, actor entity.Actor) error {
	if s.clearFunc != nil {
		return s.clearFunc(ctx, actor)
	}
	return nil
}

func TestUS5_EngagementContract_TDD(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	userToken, err := jwtManager.GenerateTokenWithProfile("user-1", "alice", false)
	require.NoError(t, err)

	wishlistUC := &engagementWishlistStub{
		listFunc: func(_ context.Context, page entity.Pagination) ([]entity.WishlistItem, int, error) {
			require.Equal(t, 0, page.Offset)
			return []entity.WishlistItem{{ID: 1, Title: "Need game key", VoteCount: 3}}, 1, nil
		},
		createFunc: func(_ context.Context, actor entity.Actor, title, description string) (entity.WishlistItem, error) {
			require.Equal(t, "user-1", actor.UserID)
			require.Equal(t, "Need game key", title)
			require.Equal(t, "For weekend", description)
			return entity.WishlistItem{ID: 2, Title: title, Description: description, UserID: actor.UserID, Username: actor.Username}, nil
		},
		updateFunc: func(_ context.Context, actor entity.Actor, item entity.WishlistItem) (entity.WishlistItem, error) {
			require.Equal(t, "user-1", actor.UserID)
			require.Equal(t, int64(2), item.ID)
			return item, nil
		},
		deleteFunc: func(_ context.Context, actor entity.Actor, itemID int64) error {
			require.Equal(t, "user-1", actor.UserID)
			require.Equal(t, int64(2), itemID)
			return nil
		},
		toggleVoteFunc: func(_ context.Context, actor entity.Actor, itemID int64) error {
			require.Equal(t, "user-1", actor.UserID)
			require.Equal(t, int64(1), itemID)
			return nil
		},
		createReviewFunc: func(_ context.Context, actor entity.Actor, review entity.Review) (entity.Review, error) {
			require.Equal(t, "user-1", actor.UserID)
			require.Equal(t, "product-1", review.ProductID)
			require.Equal(t, "order-1", review.OrderID)
			require.Equal(t, 5, review.Rating)
			require.Equal(t, "Great", review.Comment)
			review.ID = 7
			review.Username = actor.Username
			review.CreatedAt = time.Unix(1700000000, 0).UTC()
			return review, nil
		},
		listReviewsFunc: func(_ context.Context, productID string) ([]entity.Review, error) {
			require.Equal(t, "product-1", productID)
			return []entity.Review{{ID: 7, ProductID: productID, Rating: 5, Comment: "Great", Username: "alice", CreatedAt: time.Unix(1700000000, 0).UTC()}}, nil
		},
	}

	notifyUC := &engagementNotifyStub{
		inboxFunc: func(_ context.Context, actor entity.Actor, page entity.Pagination) ([]entity.UserNotification, int, error) {
			require.Equal(t, "user-1", actor.UserID)
			require.Equal(t, 0, page.Offset)
			return []entity.UserNotification{{ID: 11, UserID: actor.UserID, TitleKey: "promo.title", ContentKey: "promo.body", IsRead: false}}, 1, nil
		},
		unreadCountFunc: func(_ context.Context, actor entity.Actor) (int, error) {
			require.Equal(t, "user-1", actor.UserID)
			return 4, nil
		},
		markReadFunc: func(_ context.Context, actor entity.Actor, notificationID int64) error {
			require.Equal(t, "user-1", actor.UserID)
			require.Equal(t, int64(11), notificationID)
			return nil
		},
		markAllFunc: func(_ context.Context, actor entity.Actor) error {
			require.Equal(t, "user-1", actor.UserID)
			return nil
		},
		clearFunc: func(_ context.Context, actor entity.Actor) error {
			require.Equal(t, "user-1", actor.UserID)
			return nil
		},
	}

	app := fiber.New()
	v1.NewRoutes(
		app.Group("/v1"),
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		wishlistUC, notifyUC, nil, nil, nil, nil, nil, nil,
		jwtManager,
		"",
		logger.New("error"),
	)

	t.Run("wishlist routes", func(t *testing.T) {
		resp := engagementRequest(t, app, http.MethodGet, "/v1/wishlist?limit=20&offset=0", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = engagementRequest(t, app, http.MethodPost, "/v1/wishlist", []byte(`{"title":"Need game key","description":"For weekend"}`), "")
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		resp = engagementRequest(t, app, http.MethodPost, "/v1/wishlist", []byte(`{"title":"Need game key","description":"For weekend"}`), "Bearer "+userToken)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		resp = engagementRequest(t, app, http.MethodPatch, "/v1/wishlist/2", []byte(`{"title":"Updated","description":"Later"}`), "Bearer "+userToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = engagementRequest(t, app, http.MethodPost, "/v1/wishlist/1/vote", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp = engagementRequest(t, app, http.MethodDelete, "/v1/wishlist/2", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("review routes", func(t *testing.T) {
		resp := engagementRequest(t, app, http.MethodGet, "/v1/reviews?product_id=product-1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = engagementRequest(t, app, http.MethodPost, "/v1/reviews", []byte(`{"productId":"product-1","orderId":"order-1","rating":5,"comment":"Great"}`), "")
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		resp = engagementRequest(t, app, http.MethodPost, "/v1/reviews", []byte(`{"productId":"product-1","orderId":"order-1","rating":5,"comment":"Great"}`), "Bearer "+userToken)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("notification routes", func(t *testing.T) {
		resp := engagementRequest(t, app, http.MethodGet, "/v1/notifications", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp = engagementRequest(t, app, http.MethodGet, "/v1/notifications/unread-count", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var env engagementEnvelope
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&env))
		require.NotEmpty(t, env.Data)

		resp = engagementRequest(t, app, http.MethodPost, "/v1/notifications/11/read", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp = engagementRequest(t, app, http.MethodPost, "/v1/notifications/read-all", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp = engagementRequest(t, app, http.MethodDelete, "/v1/notifications", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func engagementRequest(t *testing.T, app *fiber.App, method, target string, body []byte, authHeader string) *http.Response {
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
