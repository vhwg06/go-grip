package media

import (
	"context"
	"path/filepath"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/google/uuid"
)

type Config struct {
	MaxBytes int64
}

type UseCase struct {
	repo    repo.MediaRepo
	storage usecase.MediaStorage
	cfg     Config
}

func New(r repo.MediaRepo, s usecase.MediaStorage, cfg Config) *UseCase {
	if cfg.MaxBytes <= 0 {
		cfg.MaxBytes = entity.MaxMediaUploadBytes
	}
	return &UseCase{
		repo:    r,
		storage: s,
		cfg:     cfg,
	}
}

func (uc *UseCase) Store(ctx context.Context, media entity.MediaAsset) (entity.MediaAsset, error) {
	if media.SizeBytes > uc.cfg.MaxBytes || !entity.IsAllowedMediaMimeType(media.MimeType) {
		return entity.MediaAsset{}, entity.ErrInvalidInput
	}
	if _, err := uuid.Parse(media.ID); err != nil {
		media.ID = uuid.New().String()
	}
	media.CreatedAt = time.Now().UTC()
	if err := uc.repo.Store(ctx, &media); err != nil {
		return entity.MediaAsset{}, err
	}
	return media, nil
}

func (uc *UseCase) List(ctx context.Context, page entity.Pagination) ([]entity.MediaAsset, int, error) {
	return uc.repo.List(ctx, page.Normalize())
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return entity.ErrInvalidInput
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

func (uc *UseCase) GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error) {
	if !entity.IsAllowedMediaMimeType(contentType) {
		return "", "", "", entity.ErrInvalidInput
	}

	return uc.storage.GeneratePresignedURL(ctx, fileName, contentType)
}
