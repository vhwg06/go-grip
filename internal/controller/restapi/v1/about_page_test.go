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

func TestStaticPagePublicReflectionWithGallery(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin", "admin", true)
	require.NoError(t, err)

	v := &V1{
		content:    contentuc.New(persistent.NewContentRepo(nil)),
		jwtManager: jwtManager,
	}

	app := fiber.New()
	apiV1Group := app.Group("/v1")
	v.registerEcommerceRoutes(apiV1Group)

	createResp := testRequest(t, app, http.MethodPost, "/v1/content/pages", []byte(`{
		"title":"About Grip",
		"slug":"about",
		"body":"About body",
		"gallery":["https://cdn.example.com/a.png","https://cdn.example.com/b.png"],
		"template_key":"about-us",
		"status":"published"
	}`), "Bearer "+adminToken)
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	updateResp := testRequest(t, app, http.MethodPatch, "/v1/content/pages/about", []byte(`{
		"title":"About Grip Updated",
		"slug":"about",
		"body":"Updated body",
		"gallery":["https://cdn.example.com/c.png"],
		"template_key":"about-us",
		"status":"published"
	}`), "Bearer "+adminToken)
	require.Equal(t, http.StatusOK, updateResp.StatusCode)

	publicResp := testRequest(t, app, http.MethodGet, "/v1/public/content/pages/about", nil, "")
	require.Equal(t, http.StatusOK, publicResp.StatusCode)

	var page struct {
		Title   string   `json:"title"`
		Body    string   `json:"body"`
		Gallery []string `json:"gallery"`
	}
	require.NoError(t, json.NewDecoder(publicResp.Body).Decode(&page))
	require.Equal(t, "About Grip Updated", page.Title)
	require.Equal(t, "Updated body", page.Body)
	require.Equal(t, []string{"https://cdn.example.com/c.png"}, page.Gallery)
}
