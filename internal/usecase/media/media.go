package media

import (
	"context"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/google/uuid"
)

type UseCase struct {
	repo     repo.MediaRepo
	maxBytes int64
}

func New(r repo.MediaRepo, maxBytes int64) *UseCase {
	if maxBytes <= 0 {
		maxBytes = entity.MaxMediaUploadBytes
	}
	return &UseCase{repo: r, maxBytes: maxBytes}
}

func (uc *UseCase) Store(ctx context.Context, media entity.MediaAsset) (entity.MediaAsset, error) {
	if media.SizeBytes > uc.maxBytes || !entity.IsAllowedMediaMimeType(media.MimeType) {
		return entity.MediaAsset{}, entity.ErrInvalidInput
	}
	media.ID = uuid.New().String()
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
	return uc.repo.Delete(ctx, id)
}
