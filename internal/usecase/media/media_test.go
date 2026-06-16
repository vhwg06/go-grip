package media

import (
	"context"
	"os"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestMediaUseCaseValidation(t *testing.T) {
	t.Parallel()
	uc := New(persistent.NewMediaRepo(nil), entity.MaxMediaUploadBytes)
	asset, err := uc.Store(context.Background(), entity.MediaAsset{FileName: "a.jpg", MimeType: "image/jpeg", SizeBytes: 10})
	require.NoError(t, err)
	require.NotEmpty(t, asset.ID)
	_, err = uc.Store(context.Background(), entity.MediaAsset{FileName: "a.gif", MimeType: "image/gif", SizeBytes: 10})
	require.ErrorIs(t, err, entity.ErrInvalidInput)
}

func TestGeneratePresignedURL(t *testing.T) {
	// Backup env to avoid interfering with other tests
	r2AccountID := os.Getenv("R2_ACCOUNT_ID")
	r2AccessKey := os.Getenv("R2_ACCESS_KEY_ID")
	r2SecretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	r2BucketName := os.Getenv("R2_BUCKET_NAME")
	r2PublicURL := os.Getenv("R2_PUBLIC_URL")

	defer func() {
		if r2AccountID != "" {
			os.Setenv("R2_ACCOUNT_ID", r2AccountID)
		} else {
			os.Unsetenv("R2_ACCOUNT_ID")
		}
		if r2AccessKey != "" {
			os.Setenv("R2_ACCESS_KEY_ID", r2AccessKey)
		} else {
			os.Unsetenv("R2_ACCESS_KEY_ID")
		}
		if r2SecretKey != "" {
			os.Setenv("R2_SECRET_ACCESS_KEY", r2SecretKey)
		} else {
			os.Unsetenv("R2_SECRET_ACCESS_KEY")
		}
		if r2BucketName != "" {
			os.Setenv("R2_BUCKET_NAME", r2BucketName)
		} else {
			os.Unsetenv("R2_BUCKET_NAME")
		}
		if r2PublicURL != "" {
			os.Setenv("R2_PUBLIC_URL", r2PublicURL)
		} else {
			os.Unsetenv("R2_PUBLIC_URL")
		}
	}()

	uc := New(persistent.NewMediaRepo(nil), entity.MaxMediaUploadBytes)

	t.Run("local simulation mode when R2 credentials are empty", func(t *testing.T) {
		os.Unsetenv("R2_ACCOUNT_ID")
		os.Unsetenv("R2_ACCESS_KEY_ID")
		os.Unsetenv("R2_SECRET_ACCESS_KEY")
		os.Unsetenv("R2_BUCKET_NAME")
		os.Unsetenv("R2_PUBLIC_URL")

		uploadURL, publicURL, fileID, err := uc.GeneratePresignedURL(context.Background(), "test.jpg", "image/jpeg")
		require.NoError(t, err)
		require.NotEmpty(t, fileID)
		require.Contains(t, uploadURL, "simulate-upload")
		require.Contains(t, publicURL, "static/uploads")
	})

	t.Run("production R2 client mode when R2 credentials are set", func(t *testing.T) {
		os.Setenv("R2_ACCOUNT_ID", "mock-account-id")
		os.Setenv("R2_ACCESS_KEY_ID", "mock-access-key-id")
		os.Setenv("R2_SECRET_ACCESS_KEY", "mock-secret-access-key")
		os.Setenv("R2_BUCKET_NAME", "mock-bucket-name")
		os.Setenv("R2_PUBLIC_URL", "https://cdn.example.com")

		uploadURL, publicURL, fileID, err := uc.GeneratePresignedURL(context.Background(), "test.jpg", "image/jpeg")
		require.NoError(t, err)
		require.NotEmpty(t, fileID)
		require.Contains(t, uploadURL, "mock-account-id.r2.cloudflarestorage.com")
		require.Contains(t, publicURL, "https://cdn.example.com")
	})

	t.Run("invalid mime type rejected", func(t *testing.T) {
		_, _, _, err := uc.GeneratePresignedURL(context.Background(), "test.gif", "image/gif")
		require.ErrorIs(t, err, entity.ErrInvalidInput)
	})
}

