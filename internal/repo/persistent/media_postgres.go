package persistent

import (
	"context"
	"sync"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type MediaRepo struct {
	*postgres.Postgres
	mu    sync.RWMutex
	items map[string]entity.MediaAsset
}

func NewMediaRepo(pg *postgres.Postgres) *MediaRepo {
	return &MediaRepo{Postgres: pg, items: map[string]entity.MediaAsset{}}
}

func (r *MediaRepo) Store(ctx context.Context, media *entity.MediaAsset) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[media.ID] = *media
	return nil
}

func (r *MediaRepo) List(ctx context.Context, page entity.Pagination) ([]entity.MediaAsset, int, error) {
	_ = ctx
	_ = page.Normalize()
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]entity.MediaAsset, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	return items, len(items), nil
}

func (r *MediaRepo) Delete(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
