package cart

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// NotificationService defines cross-module notification contract consumed by Cart use case.
type NotificationService interface {
	Dispatch(ctx context.Context, channel, to, subject string) error
}

// CartUseCase defines application service for Cart operations.
type CartUseCase interface {
	Create(ctx context.Context, sessionID string) (Cart, error)
	Get(ctx context.Context, sessionID string) (Cart, error)
	AddItem(ctx context.Context, sessionID, productID string, quantity int) (Cart, error)
	UpdateItem(ctx context.Context, sessionID, itemID string, quantity int) (Cart, error)
	RemoveItem(ctx context.Context, sessionID, itemID string) (Cart, error)
	SubmitOrder(ctx context.Context, order OrderRequest) (OrderRequest, error)
}

type cartUseCase struct {
	carts         CartRepo
	orders        OrderRequestRepo
	notifications NotificationService
}

// NewCartUseCase constructs a new CartUseCase instance.
func NewCartUseCase(carts CartRepo, orders OrderRequestRepo, notifications NotificationService) CartUseCase {
	return &cartUseCase{carts: carts, orders: orders, notifications: notifications}
}

func (uc *cartUseCase) Create(ctx context.Context, sessionID string) (Cart, error) {
	now := time.Now().UTC()
	c := Cart{ID: uuid.New().String(), SessionID: sessionID, Status: CartStatusActive, CreatedAt: now, UpdatedAt: now}
	if err := uc.carts.Store(ctx, &c); err != nil {
		return Cart{}, err
	}
	return c, nil
}

func (uc *cartUseCase) Get(ctx context.Context, sessionID string) (Cart, error) {
	return uc.carts.GetBySession(ctx, sessionID)
}

func (uc *cartUseCase) AddItem(ctx context.Context, sessionID, productID string, quantity int) (Cart, error) {
	if quantity <= 0 {
		return Cart{}, ErrInvalidInput
	}
	c, err := uc.ensureCart(ctx, sessionID)
	if err != nil {
		return Cart{}, err
	}
	item := CartItem{ID: uuid.New().String(), CartID: c.ID, ProductID: productID, Quantity: quantity}
	if err = uc.carts.AddItem(ctx, sessionID, &item); err != nil {
		return Cart{}, err
	}
	return uc.Get(ctx, sessionID)
}

func (uc *cartUseCase) UpdateItem(ctx context.Context, sessionID, itemID string, quantity int) (Cart, error) {
	if err := uc.carts.UpdateItem(ctx, sessionID, itemID, quantity); err != nil {
		return Cart{}, err
	}
	return uc.Get(ctx, sessionID)
}

func (uc *cartUseCase) RemoveItem(ctx context.Context, sessionID, itemID string) (Cart, error) {
	if err := uc.carts.RemoveItem(ctx, sessionID, itemID); err != nil {
		return Cart{}, err
	}
	return uc.Get(ctx, sessionID)
}

func (uc *cartUseCase) SubmitOrder(ctx context.Context, order OrderRequest) (OrderRequest, error) {
	c, err := uc.carts.GetBySession(ctx, order.CartID)
	if err != nil {
		return OrderRequest{}, err
	}
	for _, item := range c.Items {
		if item.Blocked {
			return OrderRequest{}, ErrCartBlocked
		}
	}
	order.ID = uuid.New().String()
	order.CartID = c.ID
	order.Status = WorkflowStatusNew
	order.CreatedAt = time.Now().UTC()
	if err = uc.orders.Store(ctx, &order); err != nil {
		return OrderRequest{}, err
	}
	_ = uc.carts.Convert(ctx, c.ID)
	if uc.notifications != nil {
		_ = uc.notifications.Dispatch(ctx, "email", order.CustomerEmail, "Order request received")
	}
	return order, nil
}

func (uc *cartUseCase) ensureCart(ctx context.Context, sessionID string) (Cart, error) {
	c, err := uc.carts.GetBySession(ctx, sessionID)
	if err == nil {
		return c, nil
	}
	return uc.Create(ctx, sessionID)
}
