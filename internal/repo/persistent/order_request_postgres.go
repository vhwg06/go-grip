package persistent

import (
	"context"
	"sync"

	cartmodule "github.com/evrone/go-clean-template/internal/module/cart"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type OrderRequestRepo struct {
	*postgres.Postgres
	mu    sync.Mutex
	items map[string]cartmodule.OrderRequest
}

func NewOrderRequestRepo(pg *postgres.Postgres) *OrderRequestRepo {
	return &OrderRequestRepo{Postgres: pg, items: map[string]cartmodule.OrderRequest{}}
}

func (r *OrderRequestRepo) Store(ctx context.Context, order *cartmodule.OrderRequest) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[order.ID] = *order
	return nil
}

func (r *OrderRequestRepo) GetByID(ctx context.Context, id string) (cartmodule.OrderRequest, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.items[id]
	if !ok {
		return cartmodule.OrderRequest{}, cartmodule.ErrNotFound
	}
	return order, nil
}

func (r *OrderRequestRepo) UpdateStatus(ctx context.Context, id string, status cartmodule.WorkflowStatus) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.items[id]
	if !ok {
		return cartmodule.ErrNotFound
	}
	order.Status = status
	r.items[id] = order
	return nil
}
