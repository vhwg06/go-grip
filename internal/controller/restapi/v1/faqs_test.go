package v1

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/repo/persistent"
	contentuc "github.com/evrone/go-clean-template/internal/usecase/content"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestFAQEndpoints(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin", "admin", true)
	require.NoError(t, err)
	userToken, err := jwtManager.GenerateTokenWithProfile("user", "user", false)
	require.NoError(t, err)

	v := &V1{
		homepage:   contentuc.NewHomepage(persistent.NewHomepageRepo(nil), persistent.NewSupportChannelRepo(nil)),
		jwtManager: jwtManager,
		adminUsers: "admin",
	}

	app := fiber.New()
	apiV1Group := app.Group("/v1")
	v.registerGripStoreRoutes(apiV1Group)

	t.Run("admin CRUD and public active projection", func(t *testing.T) {
		require.Equal(t, http.StatusUnauthorized, testRequest(t, app, http.MethodGet, "/v1/admin/faqs", nil, "").StatusCode)
		require.Equal(t, http.StatusForbidden, testRequest(t, app, http.MethodGet, "/v1/admin/faqs", nil, "Bearer "+userToken).StatusCode)

		respCreateA := testRequest(t, app, http.MethodPost, "/v1/admin/faqs", []byte(`{"question":"Q2","answer":"A2","sortOrder":20,"isActive":false}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, respCreateA.StatusCode)
		respCreateB := testRequest(t, app, http.MethodPost, "/v1/admin/faqs", []byte(`{"question":"Q1","answer":"A1","sortOrder":10,"isActive":true}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, respCreateB.StatusCode)

		respList := testRequest(t, app, http.MethodGet, "/v1/admin/faqs", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, respList.StatusCode)
		var adminItems []adminFAQEntry
		require.NoError(t, json.NewDecoder(respList.Body).Decode(&adminItems))
		require.Len(t, adminItems, 2)
		require.Equal(t, "Q1", adminItems[0].Question)
		require.Equal(t, "Q2", adminItems[1].Question)

		respPublic := testRequest(t, app, http.MethodGet, "/v1/faqs/active", nil, "")
		require.Equal(t, http.StatusOK, respPublic.StatusCode)
		var publicPayload struct {
			Items []struct {
				ID       string `json:"id"`
				Question string `json:"question"`
			} `json:"items"`
		}
		require.NoError(t, json.NewDecoder(respPublic.Body).Decode(&publicPayload))
		require.Len(t, publicPayload.Items, 1)
		require.Equal(t, "Q1", publicPayload.Items[0].Question)

		respUpdate := testRequest(t, app, http.MethodPost, "/v1/admin/faqs", []byte(`{"id":1,"question":"Q2 updated","answer":"A2 updated","sortOrder":5,"isActive":true}`), "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, respUpdate.StatusCode)

		respDelete := testRequest(t, app, http.MethodDelete, "/v1/admin/faqs/2", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, respDelete.StatusCode)

		respPublicAfter := testRequest(t, app, http.MethodGet, "/v1/faqs/active", nil, "")
		require.Equal(t, http.StatusOK, respPublicAfter.StatusCode)
		publicPayload = struct {
			Items []struct {
				ID       string `json:"id"`
				Question string `json:"question"`
			} `json:"items"`
		}{}
		require.NoError(t, json.NewDecoder(respPublicAfter.Body).Decode(&publicPayload))
		require.Len(t, publicPayload.Items, 1)
		require.Equal(t, "Q2 updated", publicPayload.Items[0].Question)
	})
}
