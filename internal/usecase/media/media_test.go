package media

import (
	"context"
	"fmt"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

type stubStorage struct {
	generatePresignedURLFn func(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error)
	deleteFn               func(ctx context.Context, key string) error
}

func (s *stubStorage) GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error) {
	if s.generatePresignedURLFn != nil {
		return s.generatePresignedURLFn(ctx, fileName, contentType)
	}
	return "", "", "", nil
}

func (s *stubStorage) Delete(ctx context.Context, key string) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, key)
	}
	return nil
}

func TestMediaUseCaseValidation(t *testing.T) {
	t.Parallel()
	storage := &stubStorage{}
	uc := New(persistent.NewMediaRepo(nil), storage, Config{MaxBytes: entity.MaxMediaUploadBytes})
	asset, err := uc.Store(context.Background(), entity.MediaAsset{FileName: "a.jpg", MimeType: "image/jpeg", SizeBytes: 10})
	require.NoError(t, err)
	require.NotEmpty(t, asset.ID)
	_, err = uc.Store(context.Background(), entity.MediaAsset{FileName: "a.gif", MimeType: "image/gif", SizeBytes: 10})
	require.ErrorIs(t, err, entity.ErrInvalidInput)
}

func TestGeneratePresignedURL(t *testing.T) {
	t.Parallel()

	t.Run("local simulation mode", func(t *testing.T) {
		storage := &stubStorage{
			generatePresignedURLFn: func(ctx context.Context, fileName string, contentType string) (string, string, string, error) {
				fileID := "test-uuid"
				return fmt.Sprintf("http://localhost:8080/v1/media/simulate-upload/%s.jpg", fileID),
					fmt.Sprintf("http://localhost:8080/static/uploads/%s.jpg", fileID),
					fileID, nil
			},
		}
		uc := New(persistent.NewMediaRepo(nil), storage, Config{MaxBytes: entity.MaxMediaUploadBytes})
		uploadURL, publicURL, fileID, err := uc.GeneratePresignedURL(context.Background(), "test.jpg", "image/jpeg")
		require.NoError(t, err)
		require.Equal(t, "test-uuid", fileID)
		require.Contains(t, uploadURL, "simulate-upload")
		require.Contains(t, publicURL, "static/uploads")
	})

	t.Run("production R2 client mode", func(t *testing.T) {
		storage := &stubStorage{
			generatePresignedURLFn: func(ctx context.Context, fileName string, contentType string) (string, string, string, error) {
				fileID := "test-uuid"
				return fmt.Sprintf("https://mock-account-id.r2.cloudflarestorage.com/mock-bucket-name/%s.jpg", fileID),
					fmt.Sprintf("https://cdn.example.com/%s.jpg", fileID),
					fileID, nil
			},
		}
		uc := New(persistent.NewMediaRepo(nil), storage, Config{MaxBytes: entity.MaxMediaUploadBytes})
		uploadURL, publicURL, fileID, err := uc.GeneratePresignedURL(context.Background(), "test.jpg", "image/jpeg")
		require.NoError(t, err)
		require.Equal(t, "test-uuid", fileID)
		require.Contains(t, uploadURL, "mock-account-id.r2.cloudflarestorage.com")
		require.Contains(t, publicURL, "https://cdn.example.com")
	})

	t.Run("invalid mime type rejected", func(t *testing.T) {
		storage := &stubStorage{}
		uc := New(persistent.NewMediaRepo(nil), storage, Config{MaxBytes: entity.MaxMediaUploadBytes})
		_, _, _, err := uc.GeneratePresignedURL(context.Background(), "test.gif", "image/gif")
		require.ErrorIs(t, err, entity.ErrInvalidInput)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("invalid uuid returns ErrInvalidInput", func(t *testing.T) {
		uc := New(persistent.NewMediaRepo(nil), &stubStorage{}, Config{MaxBytes: entity.MaxMediaUploadBytes})
		err := uc.Delete(context.Background(), "invalid-uuid")
		require.ErrorIs(t, err, entity.ErrInvalidInput)
	})

	t.Run("non-existent media returns nil (idempotent)", func(t *testing.T) {
		uc := New(persistent.NewMediaRepo(nil), &stubStorage{}, Config{MaxBytes: entity.MaxMediaUploadBytes})
		err := uc.Delete(context.Background(), "c56b9074-1234-5678-abcd-1234567890ab")
		require.NoError(t, err)
	})

	t.Run("existing media calls storage delete and repo delete", func(t *testing.T) {
		repo := persistent.NewMediaRepo(nil)
		mediaID := "c56b9074-1234-5678-abcd-1234567890ab"
		err := repo.Store(context.Background(), &entity.MediaAsset{
			ID:       mediaID,
			FileName: "test.jpg",
			URL:      "https://cdn.example.com/some-file-key.jpg",
		})
		require.NoError(t, err)

		var deletedKey string
		storage := &stubStorage{
			deleteFn: func(ctx context.Context, key string) error {
				deletedKey = key
				return nil
			},
		}

		uc := New(repo, storage, Config{MaxBytes: entity.MaxMediaUploadBytes})
		err = uc.Delete(context.Background(), mediaID)
		require.NoError(t, err)
		require.Equal(t, "some-file-key.jpg", deletedKey)

		// Assert it's deleted from repo
		_, err = repo.Get(context.Background(), mediaID)
		require.Error(t, err)
	})
}
