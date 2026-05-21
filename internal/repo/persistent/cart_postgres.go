package persistent

import (
	"context"
	"sync"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type CartRepo struct {
	*postgres.Postgres
	mu    sync.RWMutex
	items map[string]entity.Cart
}

func NewCartRepo(pg *postgres.Postgres) *CartRepo {
	return &CartRepo{Postgres: pg, items: map[string]entity.Cart{}}
}

func (r *CartRepo) Store(ctx context.Context, cart *entity.Cart) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[cart.SessionID] = *cart
	return nil
}

func (r *CartRepo) GetBySession(ctx context.Context, sessionID string) (entity.Cart, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	cart, ok := r.items[sessionID]
	if !ok {
		return entity.Cart{}, entity.ErrNotFound
	}
	return cart, nil
}

func (r *CartRepo) AddItem(ctx context.Context, sessionID string, item *entity.CartItem) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	cart := r.items[sessionID]
	cart.Items = append(cart.Items, *item)
	r.items[sessionID] = cart
	return nil
}

func (r *CartRepo) UpdateItem(ctx context.Context, sessionID, itemID string, quantity int) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	cart := r.items[sessionID]
	for i := range cart.Items {
		if cart.Items[i].ID == itemID {
			cart.Items[i].Quantity = quantity
		}
	}
	r.items[sessionID] = cart
	return nil
}

func (r *CartRepo) RemoveItem(ctx context.Context, sessionID, itemID string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	cart := r.items[sessionID]
	next := cart.Items[:0]
	for _, item := range cart.Items {
		if item.ID != itemID {
			next = append(next, item)
		}
	}
	cart.Items = next
	r.items[sessionID] = cart
	return nil
}

func (r *CartRepo) Convert(ctx context.Context, cartID string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for sessionID, cart := range r.items {
		if cart.ID == cartID {
			cart.Status = entity.CartStatusConverted
			r.items[sessionID] = cart
			return nil
		}
	}
	return entity.ErrNotFound
}
