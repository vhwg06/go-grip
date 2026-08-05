package persistent

import (
	"context"
	"sync"

	"github.com/evrone/go-clean-template/internal/module/lead"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type LeadRepo struct {
	*postgres.Postgres
	mu    sync.RWMutex
	items map[string]lead.LeadSubmission
}

func NewLeadRepo(pg *postgres.Postgres) *LeadRepo {
	return &LeadRepo{Postgres: pg, items: map[string]lead.LeadSubmission{}}
}

func (r *LeadRepo) Store(ctx context.Context, item *lead.LeadSubmission) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[item.ID] = *item
	return nil
}

func (r *LeadRepo) Get(ctx context.Context, id string) (lead.LeadSubmission, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[id]
	if !ok {
		return lead.LeadSubmission{}, lead.ErrNotFound
	}
	return item, nil
}
