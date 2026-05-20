package persistent

import (
	"context"
	"sync"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type OrderRequestRepo struct {
	*postgres.Postgres
	mu    sync.Mutex
	items map[string]entity.OrderRequest
}

func NewOrderRequestRepo(pg *postgres.Postgres) *OrderRequestRepo {
	return &OrderRequestRepo{Postgres: pg, items: map[string]entity.OrderRequest{}}
}

func (r *OrderRequestRepo) Store(ctx context.Context, order *entity.OrderRequest) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[order.ID] = *order
	return nil
}
