package v1

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	contentuc "github.com/evrone/go-clean-template/internal/usecase/content"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestContentArticlePublicReflectionAfterUpdate(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin", "admin", true)
	require.NoError(t, err)

	v := &V1{
		content:    contentuc.New(persistent.NewContentRepo(nil)),
		jwtManager: jwtManager,
		adminUsers: "admin",
	}

	app := fiber.New()
	apiV1Group := app.Group("/v1")
	v.registerGripStoreRoutes(apiV1Group)
	v.registerEcommerceRoutes(apiV1Group)

	createResp := testRequest(t, app, http.MethodPost, "/v1/content/articles", []byte(`{
		"title":"API Test Article",
		"slug":"api-test-article",
		"body":"This is a body test.",
		"status":"published",
		"image_url":"https://example.com/test.png",
		"topic":"tech",
		"tags":["playwright","api"],
		"priority":42
	}`), "Bearer "+adminToken)
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	var created entity.ContentArticle
	require.NoError(t, json.NewDecoder(createResp.Body).Decode(&created))
	require.NotEmpty(t, created.ID)

	updateResp := testRequest(t, app, http.MethodPatch, "/v1/content/articles/"+created.ID, []byte(`{
		"id":"`+created.ID+`",
		"title":"Updated API Test Article",
		"slug":"api-test-article",
		"body":"This is a body test.",
		"status":"published",
		"image_url":"https://example.com/test.png",
		"topic":"tech",
		"tags":["playwright","api"],
		"priority":99
	}`), "Bearer "+adminToken)
	require.Equal(t, http.StatusOK, updateResp.StatusCode)

	listResp := testRequest(t, app, http.MethodGet, "/v1/public/content/articles", nil, "")
	require.Equal(t, http.StatusOK, listResp.StatusCode)

	var listed struct {
		Data []entity.ContentArticle `json:"data"`
	}
	require.NoError(t, json.NewDecoder(listResp.Body).Decode(&listed))
	require.Len(t, listed.Data, 1)
	require.Equal(t, created.ID, listed.Data[0].ID)
	require.Equal(t, "Updated API Test Article", listed.Data[0].Title)

	detailResp := testRequest(t, app, http.MethodGet, "/v1/public/content/articles/"+created.ID, nil, "")
	require.Equal(t, http.StatusOK, detailResp.StatusCode)

	var detail entity.ContentArticle
	require.NoError(t, json.NewDecoder(detailResp.Body).Decode(&detail))
	require.Equal(t, created.ID, detail.ID)
	require.Equal(t, "Updated API Test Article", detail.Title)
}
