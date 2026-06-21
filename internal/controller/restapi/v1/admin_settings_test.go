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

type adminSettingsUseCaseStub struct {
	BaseAdminUseCaseStub
	listSettingsFunc  func(ctx context.Context, actor entity.Actor) ([]entity.Setting, error)
	setSettingFunc    func(ctx context.Context, actor entity.Actor, key, value string) error
	deleteSettingFunc func(ctx context.Context, actor entity.Actor, key string) error
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
