package v1

import (
	"context"
	"encoding/json"
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

type contentUseCaseStub struct {
	getArticleFunc func(ctx context.Context, idOrSlug string) (entity.ContentArticle, error)
}

func (s *contentUseCaseStub) CreateArticle(context.Context, entity.ContentArticle) (entity.ContentArticle, error) {
	return entity.ContentArticle{}, nil
}
func (s *contentUseCaseStub) UpdateArticle(context.Context, entity.ContentArticle) (entity.ContentArticle, error) {
	return entity.ContentArticle{}, nil
}
func (s *contentUseCaseStub) ListArticles(context.Context, entity.ArticleFilter) ([]entity.ContentArticle, int, error) {
	return nil, 0, nil
}
func (s *contentUseCaseStub) GetArticle(ctx context.Context, idOrSlug string) (entity.ContentArticle, error) {
	if s.getArticleFunc != nil {
		return s.getArticleFunc(ctx, idOrSlug)
	}
	return entity.ContentArticle{}, entity.ErrNotFound
}
func (s *contentUseCaseStub) DeleteArticle(context.Context, string) error { return nil }
func (s *contentUseCaseStub) CreatePage(context.Context, entity.StaticPage) (entity.StaticPage, error) {
	return entity.StaticPage{}, nil
}
func (s *contentUseCaseStub) UpdatePage(context.Context, entity.StaticPage) (entity.StaticPage, error) {
	return entity.StaticPage{}, nil
}
func (s *contentUseCaseStub) GetPage(context.Context, string) (entity.StaticPage, error) {
	return entity.StaticPage{}, nil
}
func (s *contentUseCaseStub) PublishDue(context.Context) (int, error) { return 0, nil }

func (s *adminProductUseCaseStub) UpsertProduct(ctx context.Context, actor entity.Actor, product entity.Product) (entity.Product, error) {
	if s.upsertProductFunc != nil {
		return s.upsertProductFunc(ctx, actor, product)
	}
	return entity.Product{}, nil
}

func TestAdminProductFormIncludesLinkedIntroArticle(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin-id", "admin", true)
	require.NoError(t, err)

	v := &V1{
		adminUC: &adminProductUseCaseStub{
			getProductFunc: func(_ context.Context, actor entity.Actor, productID string) (entity.Product, error) {
				require.True(t, actor.IsAdmin)
				require.Equal(t, "p1", productID)
				return entity.Product{ID: productID, Title: "Grip Pad", IntroArticleID: "art-1"}, nil
			},
		},
		content: &contentUseCaseStub{
			getArticleFunc: func(_ context.Context, idOrSlug string) (entity.ContentArticle, error) {
				require.Equal(t, "art-1", idOrSlug)
				return entity.ContentArticle{ID: "art-1", Title: "Intro Story", Status: entity.ContentStatusDraft}, nil
			},
		},
		jwtManager: jwtManager,
		adminUsers: "admin",
	}
	app := fiber.New()
	v.registerGripStoreRoutes(app.Group("/v1"))

	resp := testRequest(t, app, http.MethodGet, "/v1/admin/products/p1/form", nil, "Bearer "+adminToken)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var payload struct {
		Product entity.Product `json:"product"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
	require.Equal(t, "art-1", payload.Product.IntroArticleID)
	require.NotNil(t, payload.Product.IntroArticle)
	require.Equal(t, "Intro Story", payload.Product.IntroArticle.Title)
}

func TestAdminProductPatchSupportsClearingIntroArticle(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin-id", "admin", true)
	require.NoError(t, err)

	var saved entity.Product
	v := &V1{
		adminUC: &adminProductUseCaseStub{
			getProductFunc: func(_ context.Context, _ entity.Actor, productID string) (entity.Product, error) {
				return entity.Product{ID: productID, Title: "Grip Pad", IntroArticleID: "art-1", IsActive: true}, nil
			},
			upsertProductFunc: func(_ context.Context, _ entity.Actor, product entity.Product) (entity.Product, error) {
				saved = product
				return product, nil
			},
		},
		content:    &contentUseCaseStub{},
		jwtManager: jwtManager,
		adminUsers: "admin",
	}
	app := fiber.New()
	v.registerGripStoreRoutes(app.Group("/v1"))

	resp := testRequest(t, app, http.MethodPatch, "/v1/admin/products/p1", []byte(`{"introArticleId":null}`), "Bearer "+adminToken)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Empty(t, saved.IntroArticleID)
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
