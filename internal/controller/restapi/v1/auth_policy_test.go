package v1

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/controller/restapi/middleware"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRouteAuthPolicyForCatalogCartAndCheckout(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	v := &V1{jwtManager: jwtManager}

	app := fiber.New()
	apiV1Group := app.Group("/v1")

	v.registerGripStoreRoutes(apiV1Group)
	v.registerEcommerceRoutes(apiV1Group)
	protected := apiV1Group.Group("", middleware.Auth(jwtManager))
	protected.Get("/user/profile", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	t.Run("catalog get stays public", func(t *testing.T) {
		resp := testRequest(t, app, http.MethodGet, "/v1/catalog/products", nil, "")
		require.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("checkout user-facing routes require auth", func(t *testing.T) {
		require.Equal(t, http.StatusUnauthorized, testRequest(t, app, http.MethodPost, "/v1/checkout/orders", []byte(`{}`), "").StatusCode)
		require.Equal(t, http.StatusUnauthorized, testRequest(t, app, http.MethodPost, "/v1/checkout/payment-orders", []byte(`{}`), "").StatusCode)
		require.Equal(t, http.StatusUnauthorized, testRequest(t, app, http.MethodGet, "/v1/checkout/orders/o1/status", nil, "").StatusCode)
	})

	t.Run("cart and order-request routes require auth", func(t *testing.T) {
		require.Equal(t, http.StatusUnauthorized, testRequest(t, app, http.MethodPost, "/v1/cart", []byte(`{"session_id":"guest"}`), "").StatusCode)
		require.Equal(t, http.StatusUnauthorized, testRequest(t, app, http.MethodGet, "/v1/cart/guest", nil, "").StatusCode)
		require.Equal(t, http.StatusUnauthorized, testRequest(t, app, http.MethodPost, "/v1/cart/guest/items", []byte(`{"product_id":"p1","quantity":1}`), "").StatusCode)
		require.Equal(t, http.StatusUnauthorized, testRequest(t, app, http.MethodPost, "/v1/order-requests", []byte(`{}`), "").StatusCode)
	})

	t.Run("payment callback endpoints stay public", func(t *testing.T) {
		require.NotEqual(t, http.StatusUnauthorized, testRequest(t, app, http.MethodPost, "/v1/checkout/notify", []byte(`{}`), "").StatusCode)
		require.NotEqual(t, http.StatusUnauthorized, testRequest(t, app, http.MethodGet, "/v1/checkout/callback/o1", nil, "").StatusCode)
	})
}

func TestCartHandlersUseActorUserIDAsCartKey(t *testing.T) {
	t.Parallel()

	cartUC := &cartUseCaseStub{}
	v := &V1{cart: cartUC}

	app := fiber.New()
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Locals("userID", "user-123")
		ctx.Locals("actor", entity.Actor{UserID: "user-123"})
		return ctx.Next()
	})

	app.Post("/cart", v.createCart)
	app.Get("/cart/:session_id", v.getCart)
	app.Post("/cart/:session_id/items", v.addCartItem)
	app.Patch("/cart/:session_id/items/:item_id", v.updateCartItem)
	app.Delete("/cart/:session_id/items/:item_id", v.removeCartItem)
	app.Post("/order-requests", v.submitOrder)

	require.Equal(t, http.StatusCreated, testRequest(t, app, http.MethodPost, "/cart", []byte(`{"session_id":"guest-session"}`), "").StatusCode)
	require.Equal(t, http.StatusOK, testRequest(t, app, http.MethodGet, "/cart/guest-session", nil, "").StatusCode)
	require.Equal(t, http.StatusOK, testRequest(t, app, http.MethodPost, "/cart/guest-session/items", []byte(`{"product_id":"p1","quantity":1}`), "").StatusCode)
	require.Equal(t, http.StatusOK, testRequest(t, app, http.MethodPatch, "/cart/guest-session/items/i1", []byte(`{"quantity":2}`), "").StatusCode)
	require.Equal(t, http.StatusOK, testRequest(t, app, http.MethodDelete, "/cart/guest-session/items/i1", nil, "").StatusCode)
	require.Equal(t, http.StatusCreated, testRequest(t, app, http.MethodPost, "/order-requests", []byte(`{"cart_id":"guest-session"}`), "").StatusCode)

	require.Equal(t, "user-123", cartUC.createdSessionID)
	require.Equal(t, "user-123", cartUC.gotSessionID)
	require.Equal(t, "user-123", cartUC.addedSessionID)
	require.Equal(t, "user-123", cartUC.updatedSessionID)
	require.Equal(t, "user-123", cartUC.removedSessionID)
	require.Equal(t, "user-123", cartUC.submittedOrderCartID)
}

type cartUseCaseStub struct {
	createdSessionID     string
	gotSessionID         string
	addedSessionID       string
	updatedSessionID     string
	removedSessionID     string
	submittedOrderCartID string
}

func (s *cartUseCaseStub) Create(_ context.Context, sessionID string) (entity.Cart, error) {
	s.createdSessionID = sessionID
	return entity.Cart{ID: "c1", SessionID: sessionID}, nil
}

func (s *cartUseCaseStub) Get(_ context.Context, sessionID string) (entity.Cart, error) {
	s.gotSessionID = sessionID
	return entity.Cart{ID: "c1", SessionID: sessionID}, nil
}

func (s *cartUseCaseStub) AddItem(_ context.Context, sessionID, _ string, _ int) (entity.Cart, error) {
	s.addedSessionID = sessionID
	return entity.Cart{ID: "c1", SessionID: sessionID}, nil
}

func (s *cartUseCaseStub) UpdateItem(_ context.Context, sessionID, _ string, _ int) (entity.Cart, error) {
	s.updatedSessionID = sessionID
	return entity.Cart{ID: "c1", SessionID: sessionID}, nil
}

func (s *cartUseCaseStub) RemoveItem(_ context.Context, sessionID, _ string) (entity.Cart, error) {
	s.removedSessionID = sessionID
	return entity.Cart{ID: "c1", SessionID: sessionID}, nil
}

func (s *cartUseCaseStub) SubmitOrder(_ context.Context, order entity.OrderRequest) (entity.OrderRequest, error) {
	s.submittedOrderCartID = order.CartID
	return order, nil
}

func testRequest(t *testing.T, app *fiber.App, method, target string, body []byte, authHeader string) *http.Response {
	t.Helper()

	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader([]byte{})
	} else {
		reader = bytes.NewReader(body)
	}

	req := httptest.NewRequest(method, target, reader)
	req.Header.Set("Content-Type", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := app.Test(req)
	require.NoError(t, err)
	return resp
}
