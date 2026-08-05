package media

import (
	"context"
	"path/filepath"
	"time"

	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/google/uuid"
)

// Config holds runtime configuration options for MediaUseCase.
type Config struct {
	MaxBytes int64
}

// MediaUseCase defines the application service interface for Media asset management.
type MediaUseCase interface {
	Store(ctx context.Context, media MediaAsset) (MediaAsset, error)
	List(ctx context.Context, page pagination.Pagination, q string) ([]MediaAsset, int, error)
	Delete(ctx context.Context, id string) error
	GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error)
}

type mediaUseCase struct {
	repo    MediaRepo
	storage MediaStorageProvider
	cfg     Config
}

// NewMediaUseCase constructs a new MediaUseCase application service.
func NewMediaUseCase(r MediaRepo, s MediaStorageProvider, cfg Config) MediaUseCase {
	if cfg.MaxBytes <= 0 {
		cfg.MaxBytes = MaxMediaUploadBytes
	}
	return &mediaUseCase{
		repo:    r,
		storage: s,
		cfg:     cfg,
	}
}

// Store validates and saves media asset metadata.
func (uc *mediaUseCase) Store(ctx context.Context, media MediaAsset) (MediaAsset, error) {
	if media.SizeBytes > uc.cfg.MaxBytes || !IsAllowedMediaMimeType(media.MimeType) {
		return MediaAsset{}, ErrInvalidInput
	}
	if _, err := uuid.Parse(media.ID); err != nil {
		media.ID = uuid.New().String()
	}
	media.CreatedAt = time.Now().UTC()
	if err := uc.repo.Store(ctx, &media); err != nil {
		return MediaAsset{}, err
	}
	return media, nil
}

// List queries paginated media assets.
func (uc *mediaUseCase) List(ctx context.Context, page pagination.Pagination, q string) ([]MediaAsset, int, error) {
	return uc.repo.List(ctx, page.Normalize(), q)
}

// Delete removes a media asset from object storage and database.
func (uc *mediaUseCase) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidInput
	}

	mediaAsset, err := uc.repo.Get(ctx, id)
	if err != nil {
		// Keep delete idempotent. If not found in database, return nil
		return nil
	}

	// Delete from storage
	fileKey := filepath.Base(mediaAsset.URL)
	if fileKey != "" && fileKey != "." && fileKey != "/" {
		_ = uc.storage.Delete(ctx, fileKey)
	}

	// Delete from database
	return uc.repo.Delete(ctx, id)
}

// GeneratePresignedURL retrieves presigned upload parameters from storage provider.
func (uc *mediaUseCase) GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error) {
	if !IsAllowedMediaMimeType(contentType) {
		return "", "", "", ErrInvalidInput
	}

	return uc.storage.GeneratePresignedURL(ctx, fileName, contentType)
}
