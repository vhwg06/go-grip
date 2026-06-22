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

type storeSettingsCatalogStub struct {
	listSettingsFunc func(ctx context.Context) ([]entity.Setting, error)
}

func (s *storeSettingsCatalogStub) CreateProduct(context.Context, entity.Product) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *storeSettingsCatalogStub) ListProducts(context.Context, entity.ProductFilter) ([]entity.Product, int, error) {
	return nil, 0, nil
}

func (s *storeSettingsCatalogStub) GetProduct(context.Context, string) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *storeSettingsCatalogStub) UpdateProduct(context.Context, entity.Product) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *storeSettingsCatalogStub) DeleteProduct(context.Context, string) error {
	return nil
}

func (s *storeSettingsCatalogStub) CreateCategory(context.Context, entity.Category) (entity.Category, error) {
	return entity.Category{}, nil
}

func (s *storeSettingsCatalogStub) ListCategories(context.Context) ([]entity.Category, error) {
	return nil, nil
}

func (s *storeSettingsCatalogStub) CreateTag(context.Context, entity.Tag) (entity.Tag, error) {
	return entity.Tag{}, nil
}

func (s *storeSettingsCatalogStub) ListTags(context.Context) ([]entity.Tag, error) {
	return nil, nil
}

func (s *storeSettingsCatalogStub) ListVisibleProducts(context.Context, entity.Actor, entity.ProductFilter) ([]entity.Product, int, error) {
	return nil, 0, nil
}

func (s *storeSettingsCatalogStub) GetVisibleProduct(context.Context, entity.Actor, string) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *storeSettingsCatalogStub) ListPublicSettings(ctx context.Context) ([]entity.Setting, error) {
	if s.listSettingsFunc != nil {
		return s.listSettingsFunc(ctx)
	}
	return nil, nil
}

func (s *storeSettingsCatalogStub) GetPublicSetting(context.Context, string) (entity.Setting, error) {
	return entity.Setting{}, entity.ErrNotFound
}

func TestStoreSettingsContractEndpoints(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin-id", "admin", true)
	require.NoError(t, err)
	userToken, err := jwtManager.GenerateTokenWithProfile("user-id", "user", false)
	require.NoError(t, err)

	settingsFixture := []entity.Setting{
		{Key: "shopName", Value: "Grip QA"},
		{Key: "shopDescription", Value: "QA storefront"},
		{Key: "shopLogo", Value: "https://cdn.example.com/logo.webp"},
		{Key: "themeColor", Value: "amber"},
		{Key: "stickyBarAddress", Value: "Hanoi"},
		{Key: "stickyBarHotline", Value: "+84 903 117 742"},
		{Key: "contactEmail", Value: "qa@example.com"},
		{Key: "homepageBlocks", Value: `[{"key":"hero","enabled":true,"order":1},{"key":"categories","enabled":true,"order":2}]`},
		{Key: "homepageNewsCount", Value: "4"},
		{Key: "footerColumns", Value: `[{"id":"products","title":"Products","links":[{"label":"Door Handles","url":"/products"}]}]`},
		{Key: "footerCopyright", Value: "Copyright 2026 Grip QA"},
		{Key: "socialLinks", Value: `{"facebook":"https://facebook.com/gripqa"}`},
		{Key: "floatingSupport", Value: `[{"key":"zalo","enabled":true,"target":"https://zalo.me/gripqa"},{"key":"scroll_to_top","enabled":true,"target":null}]`},
		{Key: "noIndexEnabled", Value: "true"},
		{Key: "wishlistEnabled", Value: "false"},
		{Key: "registryOptIn", Value: "true"},
		{Key: "registryHideNav", Value: "true"},
	}

	setupApp := func(adminUC *adminSettingsUseCaseStub, catalogUC *storeSettingsCatalogStub) *fiber.App {
		v := &V1{
			adminUC:    adminUC,
			catalog:    catalogUC,
			jwtManager: jwtManager,
			adminUsers: "admin",
		}
		app := fiber.New()
		v.registerGripStoreRoutes(app.Group("/v1"))
		return app
	}

	t.Run("reads structured admin store settings payload", func(t *testing.T) {
		app := setupApp(&adminSettingsUseCaseStub{
			listSettingsFunc: func(_ context.Context, actor entity.Actor) ([]entity.Setting, error) {
				require.True(t, actor.IsAdmin)
				return settingsFixture, nil
			},
		}, &storeSettingsCatalogStub{})

		require.Equal(t, http.StatusUnauthorized, testRequest(t, app, http.MethodGet, "/v1/admin/store-settings", nil, "").StatusCode)
		require.Equal(t, http.StatusForbidden, testRequest(t, app, http.MethodGet, "/v1/admin/store-settings", nil, "Bearer "+userToken).StatusCode)

		resp := testRequest(t, app, http.MethodGet, "/v1/admin/store-settings", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body envelope
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		raw, err := json.Marshal(body.Data)
		require.NoError(t, err)

		var payload storeSettingsAdminResponse
		require.NoError(t, json.Unmarshal(raw, &payload))
		require.Equal(t, "Grip QA", payload.Config.Brand.ShopName)
		require.Equal(t, "qa@example.com", payload.Config.Contact.ContactEmail)
		require.Len(t, payload.Config.Homepage.Blocks, 2)
		require.Equal(t, 4, payload.Config.Homepage.NewsCount)
		require.Len(t, payload.Config.FloatSupport, 2)
		require.True(t, payload.Config.Visibility.NoIndexEnabled)
		require.True(t, payload.Config.Registry.Joined)
	})

	t.Run("validates and persists section writes", func(t *testing.T) {
		var saved map[string]string
		app := setupApp(&adminSettingsUseCaseStub{
			setSettingFunc: func(_ context.Context, actor entity.Actor, key, value string) error {
				require.True(t, actor.IsAdmin)
				if saved == nil {
					saved = map[string]string{}
				}
				saved[key] = value
				return nil
			},
		}, &storeSettingsCatalogStub{})

		require.Equal(t, http.StatusBadRequest, testRequest(t, app, http.MethodPut, "/v1/admin/store-settings/brand", []byte(`{"shopName":"","shopLogo":"bad-url"}`), "Bearer "+adminToken).StatusCode)
		require.Equal(t, http.StatusBadRequest, testRequest(t, app, http.MethodPut, "/v1/admin/store-settings/homepage", []byte(`{"blocks":[{"key":"hero","enabled":true,"order":1},{"key":"hero","enabled":true,"order":2}],"newsCount":-1}`), "Bearer "+adminToken).StatusCode)

		resp := testRequest(t, app, http.MethodPut, "/v1/admin/store-settings/visibility", []byte(`{"noIndexEnabled":true,"wishlistEnabled":false}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "true", saved["noIndexEnabled"])
		require.Equal(t, "false", saved["wishlistEnabled"])
	})

	t.Run("reflects same source of truth through public routes", func(t *testing.T) {
		catalogUC := &storeSettingsCatalogStub{
			listSettingsFunc: func(context.Context) ([]entity.Setting, error) {
				return settingsFixture, nil
			},
		}
		app := setupApp(&adminSettingsUseCaseStub{}, catalogUC)

		resp := testRequest(t, app, http.MethodGet, "/v1/site-config", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body envelope
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		raw, err := json.Marshal(body.Data)
		require.NoError(t, err)

		var siteConfig storeSettingsSiteConfigResponse
		require.NoError(t, json.Unmarshal(raw, &siteConfig))
		require.Equal(t, "Grip QA", siteConfig.Brand.ShopName)
		require.Equal(t, "https://cdn.example.com/logo.webp", siteConfig.Brand.ShopLogo)
		require.Equal(t, "Hanoi", siteConfig.Contact.StickyBarAddress)
		require.False(t, siteConfig.Visibility.WishlistEnabled)

		resp = testRequest(t, app, http.MethodGet, "/v1/catalog/settings", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		raw, err = json.Marshal(body.Data)
		require.NoError(t, err)

		var projection map[string]any
		require.NoError(t, json.Unmarshal(raw, &projection))
		require.Equal(t, "Grip QA", projection["shopName"])
		require.Equal(t, "amber", projection["themeColor"])
		require.Equal(t, false, projection["wishlistEnabled"])
	})
}
