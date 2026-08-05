package persistent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	cartmodule "github.com/evrone/go-clean-template/internal/module/cart"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm"
)

type CartRepo struct {
	*postgres.Postgres
	mu          sync.RWMutex
	items       map[string]cartmodule.Cart
	useInMemory bool
}

func NewCartRepo(pg *postgres.Postgres) *CartRepo {
	repo := &CartRepo{
		Postgres:    pg,
		items:       map[string]cartmodule.Cart{},
		useInMemory: pg == nil || pg.Gorm == nil,
	}
	return repo
}

func (r *CartRepo) Store(ctx context.Context, cart *cartmodule.Cart) error {
	if r.useInMemory {
		_ = ctx
		r.mu.Lock()
		defer r.mu.Unlock()
		r.items[cart.SessionID] = *cart
		return nil
	}

	now := time.Now().UTC()
	if cart.CreatedAt.IsZero() {
		cart.CreatedAt = now
	}
	cart.UpdatedAt = now

	model := models.Cart{
		ID:        cart.ID,
		SessionID: cart.SessionID,
		Status:    string(cart.Status),
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}

	if err := r.Gorm.WithContext(ctx).
		Where("session_id = ?", cart.SessionID).
		Assign(model).
		FirstOrCreate(&model).Error; err != nil {
		return fmt.Errorf("CartRepo.Store: %w", err)
	}

	cart.ID = model.ID
	cart.CreatedAt = model.CreatedAt
	cart.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *CartRepo) GetBySession(ctx context.Context, sessionID string) (cartmodule.Cart, error) {
	if r.useInMemory {
		_ = ctx
		r.mu.RLock()
		defer r.mu.RUnlock()
		c, ok := r.items[sessionID]
		if !ok {
			return cartmodule.Cart{}, cartmodule.ErrNotFound
		}
		return c, nil
	}

	var cartModel models.Cart
	if err := r.Gorm.WithContext(ctx).Where("session_id = ?", sessionID).First(&cartModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return cartmodule.Cart{}, cartmodule.ErrNotFound
		}
		return cartmodule.Cart{}, fmt.Errorf("CartRepo.GetBySession(cart): %w", err)
	}

	var itemModels []models.CartItem
	if err := r.Gorm.WithContext(ctx).Where("cart_id = ?", cartModel.ID).Order("created_at ASC").Find(&itemModels).Error; err != nil {
		return cartmodule.Cart{}, fmt.Errorf("CartRepo.GetBySession(items): %w", err)
	}

	c := cartmodule.Cart{
		ID:        cartModel.ID,
		SessionID: cartModel.SessionID,
		Status:    cartmodule.CartStatus(cartModel.Status),
		CreatedAt: cartModel.CreatedAt,
		UpdatedAt: cartModel.UpdatedAt,
		Items:     make([]cartmodule.CartItem, 0, len(itemModels)),
	}
	for _, item := range itemModels {
		snapshot := map[string]any{}
		if item.ProductSnapshot != "" {
			_ = json.Unmarshal([]byte(item.ProductSnapshot), &snapshot)
		}
		c.Items = append(c.Items, cartmodule.CartItem{
			ID:              item.ID,
			CartID:          item.CartID,
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			UnitPrice:       item.UnitPrice,
			ProductSnapshot: snapshot,
			Blocked:         item.Blocked,
		})
	}

	return c, nil
}

func (r *CartRepo) AddItem(ctx context.Context, sessionID string, item *cartmodule.CartItem) error {
	if r.useInMemory {
		_ = ctx
		r.mu.Lock()
		defer r.mu.Unlock()
		c := r.items[sessionID]
		c.Items = append(c.Items, *item)
		r.items[sessionID] = c
		return nil
	}

	var cartModel models.Cart
	if err := r.Gorm.WithContext(ctx).Where("session_id = ?", sessionID).First(&cartModel).Error; err != nil {
		return fmt.Errorf("CartRepo.AddItem(cart): %w", err)
	}

	var existing models.CartItem
	err := r.Gorm.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartModel.ID, item.ProductID).
		First(&existing).Error
	if err == nil {
		if err := r.Gorm.WithContext(ctx).Model(&existing).
			Updates(map[string]any{
				"quantity":   existing.Quantity + item.Quantity,
				"updated_at": time.Now().UTC(),
			}).Error; err != nil {
			return fmt.Errorf("CartRepo.AddItem(update existing): %w", err)
		}
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("CartRepo.AddItem(find existing): %w", err)
	}

	unitPrice := item.UnitPrice
	snapshot := item.ProductSnapshot

	var product models.Product
	if err := r.Gorm.WithContext(ctx).Where("id = ?", item.ProductID).First(&product).Error; err == nil {
		if unitPrice == 0 {
			unitPrice = product.Price
		}
		if len(snapshot) == 0 {
			snapshot = map[string]any{
				"id":    product.ID,
				"name":  product.Name,
				"title": product.Name,
				"price": product.Price,
				"image": product.Image,
			}
		}
	}

	snapshotJSON, _ := json.Marshal(snapshot)
	model := models.CartItem{
		ID:              item.ID,
		CartID:          cartModel.ID,
		ProductID:       item.ProductID,
		Quantity:        item.Quantity,
		UnitPrice:       unitPrice,
		ProductSnapshot: string(snapshotJSON),
		Blocked:         item.Blocked,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	if err := r.Gorm.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("CartRepo.AddItem(create): %w", err)
	}
	return nil
}

func (r *CartRepo) UpdateItem(ctx context.Context, sessionID, itemID string, quantity int) error {
	if r.useInMemory {
		_ = ctx
		r.mu.Lock()
		defer r.mu.Unlock()
		c := r.items[sessionID]
		for i := range c.Items {
			if c.Items[i].ID == itemID {
				c.Items[i].Quantity = quantity
			}
		}
		r.items[sessionID] = c
		return nil
	}

	var cartModel models.Cart
	if err := r.Gorm.WithContext(ctx).Where("session_id = ?", sessionID).First(&cartModel).Error; err != nil {
		return fmt.Errorf("CartRepo.UpdateItem(cart): %w", err)
	}
	if err := r.Gorm.WithContext(ctx).
		Model(&models.CartItem{}).
		Where("id = ? AND cart_id = ?", itemID, cartModel.ID).
		Updates(map[string]any{
			"quantity":   quantity,
			"updated_at": time.Now().UTC(),
		}).Error; err != nil {
		return fmt.Errorf("CartRepo.UpdateItem(update): %w", err)
	}
	return nil
}

func (r *CartRepo) RemoveItem(ctx context.Context, sessionID, itemID string) error {
	if r.useInMemory {
		_ = ctx
		r.mu.Lock()
		defer r.mu.Unlock()
		c := r.items[sessionID]
		next := c.Items[:0]
		for _, item := range c.Items {
			if item.ID != itemID {
				next = append(next, item)
			}
		}
		c.Items = next
		r.items[sessionID] = c
		return nil
	}

	var cartModel models.Cart
	if err := r.Gorm.WithContext(ctx).Where("session_id = ?", sessionID).First(&cartModel).Error; err != nil {
		return fmt.Errorf("CartRepo.RemoveItem(cart): %w", err)
	}
	if err := r.Gorm.WithContext(ctx).
		Where("id = ? AND cart_id = ?", itemID, cartModel.ID).
		Delete(&models.CartItem{}).Error; err != nil {
		return fmt.Errorf("CartRepo.RemoveItem(delete): %w", err)
	}
	return nil
}

func (r *CartRepo) Convert(ctx context.Context, cartID string) error {
	if r.useInMemory {
		_ = ctx
		r.mu.Lock()
		defer r.mu.Unlock()
		for sessionID, c := range r.items {
			if c.ID == cartID {
				c.Status = cartmodule.CartStatusConverted
				r.items[sessionID] = c
				return nil
			}
		}
		return cartmodule.ErrNotFound
	}

	result := r.Gorm.WithContext(ctx).Model(&models.Cart{}).
		Where("id = ?", cartID).
		Updates(map[string]any{
			"status":     string(cartmodule.CartStatusConverted),
			"updated_at": time.Now().UTC(),
		})
	if result.Error != nil {
		return fmt.Errorf("CartRepo.Convert: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return cartmodule.ErrNotFound
	}
	return nil
}
