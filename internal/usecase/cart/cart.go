package cart

import (
	"context"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/google/uuid"
)

type UseCase struct {
	carts         repo.CartRepo
	orders        repo.OrderRequestRepo
	notifications usecase.Notification
}

func New(carts repo.CartRepo, orders repo.OrderRequestRepo, notifications usecase.Notification) *UseCase {
	return &UseCase{carts: carts, orders: orders, notifications: notifications}
}

func (uc *UseCase) Create(ctx context.Context, sessionID string) (entity.Cart, error) {
	now := time.Now().UTC()
	cart := entity.Cart{ID: uuid.New().String(), SessionID: sessionID, Status: entity.CartStatusActive, CreatedAt: now, UpdatedAt: now}
	if err := uc.carts.Store(ctx, &cart); err != nil {
		return entity.Cart{}, err
	}
	return cart, nil
}

func (uc *UseCase) Get(ctx context.Context, sessionID string) (entity.Cart, error) {
	return uc.carts.GetBySession(ctx, sessionID)
}

func (uc *UseCase) AddItem(ctx context.Context, sessionID, productID string, quantity int) (entity.Cart, error) {
	if quantity <= 0 {
		return entity.Cart{}, entity.ErrInvalidInput
	}
	cart, err := uc.ensureCart(ctx, sessionID)
	if err != nil {
		return entity.Cart{}, err
	}
	item := entity.CartItem{ID: uuid.New().String(), CartID: cart.ID, ProductID: productID, Quantity: quantity}
	if err = uc.carts.AddItem(ctx, sessionID, &item); err != nil {
		return entity.Cart{}, err
	}
	return uc.Get(ctx, sessionID)
}

func (uc *UseCase) UpdateItem(ctx context.Context, sessionID, itemID string, quantity int) (entity.Cart, error) {
	if err := uc.carts.UpdateItem(ctx, sessionID, itemID, quantity); err != nil {
		return entity.Cart{}, err
	}
	return uc.Get(ctx, sessionID)
}

func (uc *UseCase) RemoveItem(ctx context.Context, sessionID, itemID string) (entity.Cart, error) {
	if err := uc.carts.RemoveItem(ctx, sessionID, itemID); err != nil {
		return entity.Cart{}, err
	}
	return uc.Get(ctx, sessionID)
}

func (uc *UseCase) SubmitOrder(ctx context.Context, order entity.OrderRequest) (entity.OrderRequest, error) {
	cart, err := uc.carts.GetBySession(ctx, order.CartID)
	if err != nil {
		return entity.OrderRequest{}, err
	}
	for _, item := range cart.Items {
		if item.Blocked {
			return entity.OrderRequest{}, entity.ErrCartBlocked
		}
	}
	order.ID = uuid.New().String()
	order.CartID = cart.ID
	order.Status = entity.WorkflowStatusNew
	order.CreatedAt = time.Now().UTC()
	if err = uc.orders.Store(ctx, &order); err != nil {
		return entity.OrderRequest{}, err
	}
	_ = uc.carts.Convert(ctx, cart.ID)
	if uc.notifications != nil {
		_ = uc.notifications.Dispatch(ctx, entity.Notification{Channel: "email", To: order.CustomerEmail, Subject: "Order request received"})
	}
	return order, nil
}

func (uc *UseCase) ensureCart(ctx context.Context, sessionID string) (entity.Cart, error) {
	cart, err := uc.carts.GetBySession(ctx, sessionID)
	if err == nil {
		return cart, nil
	}
	return uc.Create(ctx, sessionID)
}
