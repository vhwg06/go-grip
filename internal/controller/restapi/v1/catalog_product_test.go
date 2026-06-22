package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/usecase/catalog"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type visibleCatalogRepoStub struct {
	product entity.Product
	err     error
}

func (s *visibleCatalogRepoStub) ListVisibleProducts(context.Context, entity.Actor, repo.ProductFilter) ([]entity.Product, int, error) {
	return nil, 0, nil
}
func (s *visibleCatalogRepoStub) GetVisibleProduct(context.Context, entity.Actor, string) (entity.Product, error) {
	return s.product, s.err
}
func (s *visibleCatalogRepoStub) ListCategories(context.Context) ([]entity.Category, error) {
	return nil, nil
}
func (s *visibleCatalogRepoStub) ListSettings(context.Context) ([]entity.Setting, error) {
	return nil, nil
}
func (s *visibleCatalogRepoStub) GetSetting(context.Context, string) (entity.Setting, error) {
	return entity.Setting{}, nil
}

func TestCatalogProductDetailHydratesPublishedIntroArticle(t *testing.T) {
	t.Parallel()

	v := &V1{
		catalog: catalog.NewGrip(&visibleCatalogRepoStub{
			product: entity.Product{ID: "p1", Title: "Grip Pad", IntroArticleID: "art-1", IsActive: true},
		}),
		content: &contentUseCaseStub{
			getArticleFunc: func(_ context.Context, idOrSlug string) (entity.ContentArticle, error) {
				require.Equal(t, "art-1", idOrSlug)
				return entity.ContentArticle{ID: "art-1", Title: "Intro Story", Status: entity.ContentStatusPublished}, nil
			},
		},
	}
	app := fiber.New()
	v.registerGripStoreRoutes(app.Group("/v1"))

	resp := testRequest(t, app, http.MethodGet, "/v1/catalog/products/p1", nil, "")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var payload struct {
		Data entity.Product `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
	require.NotNil(t, payload.Data.IntroArticle)
	require.Equal(t, "Intro Story", payload.Data.IntroArticle.Title)
}

func TestCatalogProductDetailHidesDraftIntroArticle(t *testing.T) {
	t.Parallel()

	v := &V1{
		catalog: catalog.NewGrip(&visibleCatalogRepoStub{
			product: entity.Product{ID: "p1", Title: "Grip Pad", IntroArticleID: "art-1", IsActive: true},
		}),
		content: &contentUseCaseStub{
			getArticleFunc: func(_ context.Context, _ string) (entity.ContentArticle, error) {
				return entity.ContentArticle{ID: "art-1", Title: "Draft Intro", Status: entity.ContentStatusDraft}, nil
			},
		},
	}
	app := fiber.New()
	v.registerGripStoreRoutes(app.Group("/v1"))

	resp := testRequest(t, app, http.MethodGet, "/v1/catalog/products/p1", nil, "")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var payload struct {
		Data entity.Product `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
	require.Nil(t, payload.Data.IntroArticle)
}
