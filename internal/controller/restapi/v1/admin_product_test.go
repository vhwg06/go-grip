package v1

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type adminProductUseCaseStub struct {
	BaseAdminUseCaseStub
	getProductFunc    func(ctx context.Context, actor entity.Actor, productID string) (entity.Product, error)
	upsertProductFunc func(ctx context.Context, actor entity.Actor, product entity.Product) (entity.Product, error)
}

func (s *adminProductUseCaseStub) GetProduct(ctx context.Context, actor entity.Actor, productID string) (entity.Product, error) {
	if s.getProductFunc != nil {
		return s.getProductFunc(ctx, actor, productID)
	}
	return entity.Product{}, nil
}

func (s *adminProductUseCaseStub) UpsertProduct(ctx context.Context, actor entity.Actor, product entity.Product) (entity.Product, error) {
	if s.upsertProductFunc != nil {
		return s.upsertProductFunc(ctx, actor, product)
	}
	return entity.Product{}, nil
}

func TestAdminProductStatusEndpoint(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin-id", "admin", true)
	require.NoError(t, err)
	userToken, err := jwtManager.GenerateTokenWithProfile("user-id", "user", false)
	require.NoError(t, err)

	setupApp := func(uc *adminProductUseCaseStub) *fiber.App {
		v := &V1{
			adminUC:    uc,
			jwtManager: jwtManager,
			adminUsers: "admin",
		}
		app := fiber.New()
		v.registerGripStoreRoutes(app.Group("/v1"))
		return app
	}

	t.Run("patch product status requires admin and forwards boolean state", func(t *testing.T) {
		var saved entity.Product
		app := setupApp(&adminProductUseCaseStub{
			getProductFunc: func(_ context.Context, actor entity.Actor, productID string) (entity.Product, error) {
				require.True(t, actor.IsAdmin)
				require.Equal(t, "p1", productID)
				return entity.Product{ID: productID, Title: "Grip Pad", IsActive: true}, nil
			},
			upsertProductFunc: func(_ context.Context, actor entity.Actor, product entity.Product) (entity.Product, error) {
				require.True(t, actor.IsAdmin)
				saved = product
				return product, nil
			},
		})

		require.Equal(t, http.StatusForbidden, testRequest(t, app, http.MethodPatch, "/v1/admin/products/p1/status", []byte(`{"isActive":false}`), "Bearer "+userToken).StatusCode)
		require.Equal(t, http.StatusBadRequest, testRequest(t, app, http.MethodPatch, "/v1/admin/products/p1/status", []byte(`{bad}`), "Bearer "+adminToken).StatusCode)
		require.Equal(t, http.StatusBadRequest, testRequest(t, app, http.MethodPatch, "/v1/admin/products/p1/status", []byte(`{}`), "Bearer "+adminToken).StatusCode)

		resp := testRequest(t, app, http.MethodPatch, "/v1/admin/products/p1/status", []byte(`{"isActive":false}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "p1", saved.ID)
		require.False(t, saved.IsActive)
	})
}
