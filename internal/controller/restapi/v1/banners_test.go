package v1

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type homepageUseCaseStub struct {
	blocks     []entity.HomepageBlock
	storeBlock func(ctx context.Context, block entity.HomepageBlock) (entity.HomepageBlock, error)
	updateBlock func(ctx context.Context, block entity.HomepageBlock) (entity.HomepageBlock, error)
	deleteBlock func(ctx context.Context, id string) error
}

func (s *homepageUseCaseStub) StoreBlock(ctx context.Context, block entity.HomepageBlock) (entity.HomepageBlock, error) {
	if s.storeBlock != nil {
		return s.storeBlock(ctx, block)
	}
	block.ID = "new-uuid"
	s.blocks = append(s.blocks, block)
	return block, nil
}

func (s *homepageUseCaseStub) ListBlocks(ctx context.Context, activeOnly bool) ([]entity.HomepageBlock, error) {
	return s.blocks, nil
}

func (s *homepageUseCaseStub) UpdateBlock(ctx context.Context, block entity.HomepageBlock) (entity.HomepageBlock, error) {
	if s.updateBlock != nil {
		return s.updateBlock(ctx, block)
	}
	for i := range s.blocks {
		if s.blocks[i].BlockType == block.BlockType {
			s.blocks[i] = block
			return block, nil
		}
	}
	return block, nil
}

func (s *homepageUseCaseStub) DeleteBlock(ctx context.Context, id string) error {
	if s.deleteBlock != nil {
		return s.deleteBlock(ctx, id)
	}
	return nil
}

func (s *homepageUseCaseStub) ListSupport(ctx context.Context, enabledOnly bool) ([]entity.SupportChannel, error) {
	return nil, nil
}

func (s *homepageUseCaseStub) UpdateSupport(ctx context.Context, channel entity.SupportChannel) (entity.SupportChannel, error) {
	return channel, nil
}

func TestBannerEndpoints(t *testing.T) {
	jwtManager := jwt.New("secret", time.Hour)
	adminToken, err := jwtManager.GenerateTokenWithProfile("admin", "admin", true)
	require.NoError(t, err)

	setupApp := func(uc *homepageUseCaseStub) *fiber.App {
		v := &V1{
			homepage:   uc,
			jwtManager: jwtManager,
			adminUsers: "admin",
		}
		app := fiber.New()
		apiV1Group := app.Group("/v1")
		v.registerGripStoreRoutes(apiV1Group)
		return app
	}

	t.Run("Create and List banners", func(t *testing.T) {
		uc := &homepageUseCaseStub{}
		app := setupApp(uc)

		// 1. Create a banner via multipart/form-data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("title", "Test Title")
		_ = writer.WriteField("subtitle", "Test Subtitle")
		_ = writer.WriteField("image", "test-image.png")
		_ = writer.WriteField("sortOrder", "10")
		_ = writer.WriteField("isActive", "true")
		_ = writer.Close()

		req := httptest.NewRequest("POST", "/v1/admin/banners", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+adminToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var actionResult map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&actionResult)
		require.True(t, actionResult["success"].(bool))

		// 2. List the banners
		reqList := httptest.NewRequest("GET", "/v1/admin/banners", nil)
		reqList.Header.Set("Authorization", "Bearer "+adminToken)
		respList, err := app.Test(reqList)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, respList.StatusCode)

		var slides []AdminBannerSlide
		_ = json.NewDecoder(respList.Body).Decode(&slides)
		require.Len(t, slides, 1)
		require.Equal(t, "Test Title", slides[0].Title)
		require.Equal(t, 1, slides[0].ID) // First created should get ID 1

		// 3. Update the banner
		bodyUpdate := &bytes.Buffer{}
		writerUpdate := multipart.NewWriter(bodyUpdate)
		_ = writerUpdate.WriteField("id", "1")
		_ = writerUpdate.WriteField("title", "Updated Title")
		_ = writerUpdate.WriteField("image", "test-image-updated.png")
		_ = writerUpdate.WriteField("sortOrder", "5")
		_ = writerUpdate.WriteField("isActive", "true")
		_ = writerUpdate.Close()

		reqUpdate := httptest.NewRequest("POST", "/v1/admin/banners", bodyUpdate)
		reqUpdate.Header.Set("Content-Type", writerUpdate.FormDataContentType())
		reqUpdate.Header.Set("Authorization", "Bearer "+adminToken)

		respUpdate, err := app.Test(reqUpdate)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, respUpdate.StatusCode)

		// Verify update
		reqList2 := httptest.NewRequest("GET", "/v1/admin/banners", nil)
		reqList2.Header.Set("Authorization", "Bearer "+adminToken)
		respList2, _ := app.Test(reqList2)
		var slides2 []AdminBannerSlide
		_ = json.NewDecoder(respList2.Body).Decode(&slides2)
		require.Len(t, slides2, 1)
		require.Equal(t, "Updated Title", slides2[0].Title)
		require.Equal(t, 5, slides2[0].SortOrder)

		// 4. Delete the banner
		reqDelete := httptest.NewRequest("DELETE", "/v1/admin/banners/1", nil)
		reqDelete.Header.Set("Authorization", "Bearer "+adminToken)
		respDelete, err := app.Test(reqDelete)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, respDelete.StatusCode)

		// Verify deletion
		reqList3 := httptest.NewRequest("GET", "/v1/admin/banners", nil)
		reqList3.Header.Set("Authorization", "Bearer "+adminToken)
		respList3, _ := app.Test(reqList3)
		var slides3 []AdminBannerSlide
		_ = json.NewDecoder(respList3.Body).Decode(&slides3)
		require.Len(t, slides3, 0)
	})
}
