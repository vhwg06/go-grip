package media

import (
	"context"

	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// MediaRepo defines the persistence port owned by the Media module.
type MediaRepo interface {
	Store(ctx context.Context, media *MediaAsset) error
	List(ctx context.Context, page pagination.Pagination, query string) ([]MediaAsset, int, error)
	Get(ctx context.Context, id string) (MediaAsset, error)
	Delete(ctx context.Context, id string) error
}

// MediaStorageProvider defines the external cloud/local object storage capability port.
type MediaStorageProvider interface {
	GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error)
	Delete(ctx context.Context, key string) error
}
