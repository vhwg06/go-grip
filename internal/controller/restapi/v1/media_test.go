package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// mediaUseCaseStub implements usecase.Media for testing.
type mediaUseCaseStub struct {
	storeFunc                func(ctx context.Context, media entity.MediaAsset) (entity.MediaAsset, error)
	listFunc                 func(ctx context.Context, page entity.Pagination) ([]entity.MediaAsset, int, error)
	deleteFunc               func(ctx context.Context, id string) error
	generatePresignedURLFunc func(ctx context.Context, fileName string, contentType string) (string, string, string, error)
}

func (s *mediaUseCaseStub) Store(ctx context.Context, media entity.MediaAsset) (entity.MediaAsset, error) {
	if s.storeFunc != nil {
		return s.storeFunc(ctx, media)
	}
	return media, nil
}

func (s *mediaUseCaseStub) List(ctx context.Context, page entity.Pagination) ([]entity.MediaAsset, int, error) {
	if s.listFunc != nil {
		return s.listFunc(ctx, page)
	}
	return nil, 0, nil
}

func (s *mediaUseCaseStub) Delete(ctx context.Context, id string) error {
	if s.deleteFunc != nil {
		return s.deleteFunc(ctx, id)
	}
	return nil
}

func (s *mediaUseCaseStub) GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (string, string, string, error) {
	if s.generatePresignedURLFunc != nil {
		return s.generatePresignedURLFunc(ctx, fileName, contentType)
	}
	return "http://upload-url", "http://public-url", "file-123", nil
}

func TestMediaEndpoints(t *testing.T) {
	t.Parallel()

	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin", "admin", true)
	require.NoError(t, err)
	userToken, err := jwtManager.GenerateTokenWithProfile("user", "user", false)
	require.NoError(t, err)

	setupApp := func(uc *mediaUseCaseStub) *fiber.App {
		v := &V1{
			media:      uc,
			jwtManager: jwtManager,
			adminUsers: "admin",
		}
		app := fiber.New()
		apiV1Group := app.Group("/v1")
		v.registerGripStoreRoutes(apiV1Group)
		v.registerEcommerceRoutes(apiV1Group)
		return app
	}

	t.Run("createMedia endpoint", func(t *testing.T) {
		uc := &mediaUseCaseStub{
			storeFunc: func(ctx context.Context, media entity.MediaAsset) (entity.MediaAsset, error) {
				if media.FileName == "invalid" {
					return entity.MediaAsset{}, errors.New("bad file")
				}
				media.ID = "media-id-123"
				return media, nil
			},
		}
		app := setupApp(uc)

		// Test unauthorized request
		resp := testRequest(t, app, http.MethodPost, "/v1/media", []byte(`{"file_name":"test.jpg"}`), "")
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Test valid request with authorized token
		resp = testRequest(t, app, http.MethodPost, "/v1/media", []byte(`{"file_name":"test.jpg"}`), "Bearer "+userToken)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		var result entity.MediaAsset
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, "media-id-123", result.ID)

		// Test invalid json body
		resp = testRequest(t, app, http.MethodPost, "/v1/media", []byte(`{invalid-json}`), "Bearer "+userToken)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Test usecase store error
		resp = testRequest(t, app, http.MethodPost, "/v1/media", []byte(`{"file_name":"invalid"}`), "Bearer "+userToken)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("listMedia endpoint", func(t *testing.T) {
		uc := &mediaUseCaseStub{
			listFunc: func(ctx context.Context, page entity.Pagination) ([]entity.MediaAsset, int, error) {
				if page.Limit == 999 {
					return nil, 0, errors.New("limit exceeded")
				}
				return []entity.MediaAsset{
					{ID: "id-1", FileName: "one.png"},
					{ID: "id-2", FileName: "two.png"},
				}, 2, nil
			},
		}
		app := setupApp(uc)

		// Test unauthorized request
		resp := testRequest(t, app, http.MethodGet, "/v1/media", nil, "")
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Test valid request
		resp = testRequest(t, app, http.MethodGet, "/v1/media?limit=10&offset=0", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var listResp struct {
			Data []entity.MediaAsset `json:"data"`
			Meta entity.Page         `json:"meta"`
		}
		err = json.NewDecoder(resp.Body).Decode(&listResp)
		require.NoError(t, err)
		require.Len(t, listResp.Data, 2)
		require.Equal(t, 2, listResp.Meta.Total)

		// Test error returned from usecase
		resp = testRequest(t, app, http.MethodGet, "/v1/media?limit=999", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("deleteMedia endpoint", func(t *testing.T) {
		var deletedID string
		uc := &mediaUseCaseStub{
			deleteFunc: func(ctx context.Context, id string) error {
				if id == "invalid-id" {
					return errors.New("delete error")
				}
				deletedID = id
				return nil
			},
		}
		app := setupApp(uc)

		// Test unauthorized request
		resp := testRequest(t, app, http.MethodDelete, "/v1/media/123", nil, "")
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Test valid request
		resp = testRequest(t, app, http.MethodDelete, "/v1/media/123", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		require.Equal(t, "123", deletedID)

		// Test error returned from usecase
		resp = testRequest(t, app, http.MethodDelete, "/v1/media/invalid-id", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("getPresignedURL endpoint", func(t *testing.T) {
		uc := &mediaUseCaseStub{
			generatePresignedURLFunc: func(ctx context.Context, fileName string, contentType string) (string, string, string, error) {
				if fileName == "fail.jpg" {
					return "", "", "", errors.New("presign error")
				}
				return "https://upload-url/" + fileName, "https://public-url/" + fileName, "mock-id", nil
			},
		}
		app := setupApp(uc)

		// Test unauthorized
		resp := testRequest(t, app, http.MethodGet, "/v1/admin/media/presigned?fileName=test.jpg&contentType=image/jpeg", nil, "")
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Test non-admin token
		resp = testRequest(t, app, http.MethodGet, "/v1/admin/media/presigned?fileName=test.jpg&contentType=image/jpeg", nil, "Bearer "+userToken)
		require.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Test admin token (success)
		resp = testRequest(t, app, http.MethodGet, "/v1/admin/media/presigned?fileName=test.jpg&contentType=image/jpeg", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var result map[string]string
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, "https://upload-url/test.jpg", result["upload_url"])
		require.Equal(t, "https://public-url/test.jpg", result["public_url"])
		require.Equal(t, "mock-id", result["id"])

		// Test missing query params
		resp = testRequest(t, app, http.MethodGet, "/v1/admin/media/presigned?fileName=test.jpg", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Test usecase generate error
		resp = testRequest(t, app, http.MethodGet, "/v1/admin/media/presigned?fileName=fail.jpg&contentType=image/jpeg", nil, "Bearer "+adminToken)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("simulateUpload endpoint", func(t *testing.T) {
		app := setupApp(&mediaUseCaseStub{})
		resp := testRequest(t, app, http.MethodPut, "/v1/media/simulate-upload/test.jpg", []byte("file-data"), "")
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
