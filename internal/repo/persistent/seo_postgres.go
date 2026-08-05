package persistent

import (
	"context"
	"sync"

	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type SEORepo struct {
	*postgres.Postgres
	mu    sync.RWMutex
	items map[string]contentmodule.SeoMetadata
}

func NewSEORepo(pg *postgres.Postgres) *SEORepo {
	return &SEORepo{Postgres: pg, items: map[string]contentmodule.SeoMetadata{}}
}

func (r *SEORepo) Store(ctx context.Context, meta *contentmodule.SeoMetadata) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[meta.OwnerType+":"+meta.OwnerID] = *meta
	return nil
}

func (r *SEORepo) GetByOwner(ctx context.Context, ownerType, ownerID string) (contentmodule.SeoMetadata, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[ownerType+":"+ownerID]
	if !ok {
		return contentmodule.SeoMetadata{}, contentmodule.ErrNotFound
	}
	return item, nil
}
