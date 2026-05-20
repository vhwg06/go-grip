package persistent

import (
	"context"
	"sync"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type LeadRepo struct {
	*postgres.Postgres
	mu    sync.RWMutex
	items map[string]entity.LeadSubmission
}

func NewLeadRepo(pg *postgres.Postgres) *LeadRepo {
	return &LeadRepo{Postgres: pg, items: map[string]entity.LeadSubmission{}}
}

func (r *LeadRepo) Store(ctx context.Context, lead *entity.LeadSubmission) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[lead.ID] = *lead
	return nil
}

func (r *LeadRepo) Get(ctx context.Context, id string) (entity.LeadSubmission, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	lead, ok := r.items[id]
	if !ok {
		return entity.LeadSubmission{}, entity.ErrNotFound
	}
	return lead, nil
}
