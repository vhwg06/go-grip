package cart

import "context"

// CartRepo defines persistence port for Cart aggregate operations.
type CartRepo interface {
	Store(ctx context.Context, cart *Cart) error
	GetBySession(ctx context.Context, sessionID string) (Cart, error)
	AddItem(ctx context.Context, sessionID string, item *CartItem) error
	UpdateItem(ctx context.Context, sessionID, itemID string, quantity int) error
	RemoveItem(ctx context.Context, sessionID, itemID string) error
	Convert(ctx context.Context, cartID string) error
}

// OrderRequestRepo defines persistence port for OrderRequest operations.
type OrderRequestRepo interface {
	Store(ctx context.Context, order *OrderRequest) error
	GetByID(ctx context.Context, id string) (OrderRequest, error)
	UpdateStatus(ctx context.Context, id string, status WorkflowStatus) error
}
