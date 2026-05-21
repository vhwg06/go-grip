package persistent

import (
	"context"
	"sync"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type HomepageRepo struct {
	*postgres.Postgres
	mu    sync.RWMutex
	items map[string]entity.HomepageBlock
}

func NewHomepageRepo(pg *postgres.Postgres) *HomepageRepo {
	return &HomepageRepo{Postgres: pg, items: map[string]entity.HomepageBlock{}}
}

func (r *HomepageRepo) Store(ctx context.Context, block *entity.HomepageBlock) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[block.ID] = *block
	return nil
}

func (r *HomepageRepo) List(ctx context.Context, activeOnly bool) ([]entity.HomepageBlock, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]entity.HomepageBlock, 0, len(r.items))
	for _, item := range r.items {
		if activeOnly && !item.IsActive {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *HomepageRepo) Update(ctx context.Context, block *entity.HomepageBlock) error {
	return r.Store(ctx, block)
}

func (r *HomepageRepo) Delete(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
