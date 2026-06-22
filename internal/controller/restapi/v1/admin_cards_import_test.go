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

type adminCardsImportUseCaseStub struct {
	BaseAdminUseCaseStub
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

func TestAdminOperationalUtilityEndpoints(t *testing.T) {
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
