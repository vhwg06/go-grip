package persistent

import (
	"context"
	"sync"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type SEORepo struct {
	*postgres.Postgres
	mu    sync.RWMutex
	items map[string]entity.SeoMetadata
}

func NewSEORepo(pg *postgres.Postgres) *SEORepo {
	return &SEORepo{Postgres: pg, items: map[string]entity.SeoMetadata{}}
}

func (r *SEORepo) Store(ctx context.Context, meta *entity.SeoMetadata) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[meta.OwnerType+":"+meta.OwnerID] = *meta
	return nil
}

func (r *SEORepo) GetByOwner(ctx context.Context, ownerType, ownerID string) (entity.SeoMetadata, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[ownerType+":"+ownerID]
	if !ok {
		return entity.SeoMetadata{}, entity.ErrNotFound
	}
	return item, nil
}
