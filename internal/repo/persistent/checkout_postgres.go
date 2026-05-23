package persistent

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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
		var product models.Product
		if err := forUpdate(tx).
			Where("id = ?", order.ProductID).
			Where("is_active = ?", true).
			First(&product).Error; err != nil {
			return fmt.Errorf("CheckoutRepo.CreateOrderWithReservation: select product: %w", err)
		}

		reservedCards, err := r.reserveCardsTx(tx, order.ID, order.ProductID, order.Quantity, product.IsShared)
		if err != nil {
			return err
		}

		cardIDs := make([]string, 0, len(reservedCards))
		for _, card := range reservedCards {
			cardIDs = append(cardIDs, strconv.FormatInt(card.ID, 10))
		}
		orderModel.CardIDs = strings.Join(cardIDs, ",")
		if product.IsShared && len(reservedCards) > 0 {
			orderModel.CardKey = reservedCards[0].CardKey
		}

		if actor.UserID != "" && order.PointsUsed > 0 {
			pointsResult := tx.Model(&models.User{}).
				Where("id = ? AND points >= ?", actor.UserID, order.PointsUsed).
				Update("points", gorm.Expr("points - ?", order.PointsUsed))
			if pointsResult.Error != nil {
				return fmt.Errorf("CheckoutRepo.CreateOrderWithReservation: deduct points: %w", pointsResult.Error)
			}
			if pointsResult.RowsAffected == 0 {
				return entity.ErrPointsInsufficient
			}
		}

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
	processedAt := time.Time{}
	if payment.ProcessedAt != nil {
		processedAt = *payment.ProcessedAt
	}

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
		ProcessedAt:            processedAt,
		CreatedAt:              time.Now().UTC(),
	}

	db := r.Gorm.WithContext(ctx)
	if model.ProviderPaymentID != "" {
		var existing models.Payment
		if err := db.Where("provider_payment_id = ?", model.ProviderPaymentID).First(&existing).Error; err == nil {
			return nil
		}
	}

	if err := db.Create(&model).Error; err != nil {
		return fmt.Errorf("CheckoutRepo.AttachPayment: %w", err)
	}
	return nil
}

func (r *CheckoutRepo) UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error {
	query := r.Gorm.WithContext(ctx).
		Model(&models.Order{}).
		Where("order_id = ?", orderID)

	if status == entity.OrderStatusPaid {
		query = query.Where("status = ?", string(entity.OrderStatusPending))
	}
	if status == entity.OrderStatusDelivered {
		query = query.Where("status <> ?", string(entity.OrderStatusDelivered))
	}

	now := time.Now().UTC()
	updates := map[string]any{
		"status":     string(status),
		"updated_at": now,
	}
	if status == entity.OrderStatusPaid {
		updates["paid_at"] = now
	}
	if status == entity.OrderStatusDelivered {
		updates["delivered_at"] = now
	}

	result := query.Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("CheckoutRepo.UpdateOrderStatus: %w", result.Error)
	}
	if result.RowsAffected == 0 && status == entity.OrderStatusPaid {
		return nil
	}

	return nil
}

func (r *CheckoutRepo) ReserveCards(ctx context.Context, orderID, productID string, quantity int, isShared bool) ([]entity.Card, error) {
	reserved := make([]entity.Card, 0, quantity)

	err := withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		cards, err := r.reserveCardsTx(tx, orderID, productID, quantity, isShared)
		if err != nil {
			return err
		}
		reserved = append(reserved, cards...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return reserved, nil
}

func (r *CheckoutRepo) reserveCardsTx(tx *gorm.DB, orderID, productID string, quantity int, isShared bool) ([]entity.Card, error) {
	reserved := make([]entity.Card, 0, quantity)
	now := time.Now().UTC()
	expireBefore := now.Add(-5 * time.Minute)

	if err := tx.Model(&models.Card{}).
		Where("reserved_order_id <> ''").
		Where("reserved_at < ?", expireBefore).
		Updates(map[string]any{"reserved_order_id": "", "reserved_at": time.Time{}}).Error; err != nil {
		return nil, fmt.Errorf("CheckoutRepo.reserveCardsTx(expire old): %w", err)
	}

	if isShared {
		var card models.Card
		if err := forUpdate(tx).
			Where("product_id = ?", productID).
			Where("is_used = ?", false).
			Where("expires_at IS NULL OR expires_at = ? OR expires_at > ?", time.Time{}, now).
			Order("id ASC").
			First(&card).Error; err != nil {
			return nil, fmt.Errorf("CheckoutRepo.reserveCardsTx(shared): %w", err)
		}
		reserved = append(reserved, entity.Card{
			ID:        card.ID,
			ProductID: card.ProductID,
			CardKey:   card.CardKey,
			IsUsed:    card.IsUsed,
			CreatedAt: card.CreatedAt,
		})
		return reserved, nil
	}

	var cards []models.Card
	if err := forUpdate(tx).
		Where("product_id = ?", productID).
		Where("is_used = ?", false).
		Where("expires_at IS NULL OR expires_at = ? OR expires_at > ?", time.Time{}, now).
		Where("(reserved_order_id = '' OR reserved_order_id IS NULL OR reserved_at < ?)", expireBefore).
		Limit(quantity).
		Find(&cards).Error; err != nil {
		return nil, fmt.Errorf("CheckoutRepo.reserveCardsTx(find): %w", err)
	}
	if len(cards) < quantity {
		return nil, entity.ErrOutOfStock
	}

	for _, card := range cards {
		if err := tx.Model(&models.Card{}).
			Where("id = ?", card.ID).
			Updates(map[string]any{"reserved_order_id": orderID, "reserved_at": now}).Error; err != nil {
			return nil, fmt.Errorf("CheckoutRepo.reserveCardsTx(update): %w", err)
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
