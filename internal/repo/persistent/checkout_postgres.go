package persistent

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm"
)

type CheckoutRepo struct {
	*postgres.Postgres
}

func NewCheckoutRepo(pg *postgres.Postgres) *CheckoutRepo {
	return &CheckoutRepo{Postgres: pg}
}

var _ repo.CheckoutRepository = (*CheckoutRepo)(nil)

func (r *CheckoutRepo) CreateOrderWithReservation(ctx context.Context, actor entity.Actor, order entity.Order) (entity.Order, error) {
	orderModel := models.EntityToOrder(order)
	orderModel.UserID = actor.UserID
	orderModel.Username = actor.Username
	orderModel.Status = string(entity.OrderStatusPending)
	now := time.Now().UTC()
	orderModel.CreatedAt = now
	orderModel.UpdatedAt = now

	if err := withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		if err := tx.Create(&orderModel).Error; err != nil {
			return fmt.Errorf("CheckoutRepo.CreateOrderWithReservation: create order: %w", err)
		}
		return nil
	}); err != nil {
		return entity.Order{}, err
	}

	return models.OrderToEntity(orderModel), nil
}

func (r *CheckoutRepo) AttachPayment(ctx context.Context, payment entity.Payment) error {
	model := models.Payment{
		ID:                     payment.ID,
		OrderID:                payment.OrderID,
		Provider:               payment.Provider,
		ProviderPaymentID:      payment.ProviderPaymentID,
		Amount:                 int64(payment.Amount),
		Status:                 payment.Status,
		RequestPayloadSummary:  payment.RequestPayloadSummary,
		CallbackPayloadSummary: payment.CallbackPayloadSummary,
		IsSignatureValid:       payment.IsSignatureValid,
		ProcessedAt:            denullTime(payment.ProcessedAt),
		CreatedAt:              time.Now().UTC(),
	}
	if err := r.Gorm.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("CheckoutRepo.AttachPayment: %w", err)
	}
	return nil
}

func (r *CheckoutRepo) UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error {
	err := r.Gorm.WithContext(ctx).
		Model(&models.Order{}).
		Where("order_id = ?", orderID).
		Updates(map[string]any{
			"status":     string(status),
			"updated_at": time.Now().UTC(),
		}).Error
	if err != nil {
		return fmt.Errorf("CheckoutRepo.UpdateOrderStatus: %w", err)
	}

	return nil
}

func (r *CheckoutRepo) ReserveCards(ctx context.Context, orderID, productID string, quantity int, isShared bool) ([]entity.Card, error) {
	reserved := make([]entity.Card, 0, quantity)
	now := time.Now().UTC()

	err := withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		if isShared {
			var card models.Card
			if err := tx.Where("product_id = ?", productID).Where("is_used = ?", false).Order("id ASC").First(&card).Error; err != nil {
				return fmt.Errorf("CheckoutRepo.ReserveCards(shared): %w", err)
			}
			reserved = append(reserved, entity.Card{
				ID:        card.ID,
				ProductID: card.ProductID,
				CardKey:   card.CardKey,
				IsUsed:    card.IsUsed,
				CreatedAt: card.CreatedAt,
			})
			return nil
		}

		var cards []models.Card
		if err := forUpdate(tx).
			Where("product_id = ?", productID).
			Where("is_used = ?", false).
			Where("reserved_order_id = '' OR reserved_order_id IS NULL").
			Limit(quantity).
			Find(&cards).Error; err != nil {
			return fmt.Errorf("CheckoutRepo.ReserveCards(find): %w", err)
		}
		if len(cards) < quantity {
			return entity.ErrOutOfStock
		}

		for _, card := range cards {
			if err := tx.Model(&models.Card{}).
				Where("id = ?", card.ID).
				Updates(map[string]any{"reserved_order_id": orderID, "reserved_at": now}).Error; err != nil {
				return fmt.Errorf("CheckoutRepo.ReserveCards(update): %w", err)
			}
			reserved = append(reserved, entity.Card{
				ID:              card.ID,
				ProductID:       card.ProductID,
				CardKey:         card.CardKey,
				IsUsed:          card.IsUsed,
				ReservedOrderID: orderID,
				ReservedAt:      &now,
				CreatedAt:       card.CreatedAt,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return reserved, nil
}

func (r *CheckoutRepo) DeductPoints(ctx context.Context, userID string, points int) error {
	result := r.Gorm.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ? AND points >= ?", userID, points).
		UpdateColumn("points", gorm.Expr("points - ?", points))

	if result.Error != nil {
		return fmt.Errorf("CheckoutRepo.DeductPoints: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entity.ErrPointsInsufficient
	}
	return nil
}

func (r *CheckoutRepo) ReleaseReservation(ctx context.Context, orderID string) error {
	if err := r.Gorm.WithContext(ctx).
		Model(&models.Card{}).
		Where("reserved_order_id = ?", orderID).
		Updates(map[string]any{
			"reserved_order_id": "",
			"reserved_at":       time.Time{},
		}).Error; err != nil {
		return fmt.Errorf("CheckoutRepo.ReleaseReservation: %w", err)
	}
	return nil
}
