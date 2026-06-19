package persistent

import (
	"context"
	"errors"
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
		var product models.Product
		if err := forUpdate(tx).
			Where("id = ?", order.ProductID).
			Where("is_active = ?", true).
			First(&product).Error; err != nil {
			return fmt.Errorf("CheckoutRepo.CreateOrderWithReservation: select product: %w", err)
		}

		if product.StockCount-product.LockedCount < order.Quantity {
			return entity.ErrOutOfStock
		}

		if err := tx.Model(&models.Product{}).
			Where("id = ?", product.ID).
			UpdateColumn("locked_count", gorm.Expr("locked_count + ?", order.Quantity)).Error; err != nil {
			return fmt.Errorf("CheckoutRepo.CreateOrderWithReservation: update product stock: %w", err)
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
	return withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Where("order_id = ?", orderID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("CheckoutRepo.UpdateOrderStatus: select order: %w", err)
		}

		if order.Status == string(status) {
			return nil
		}

		currentStatus := entity.OrderStatus(order.Status)
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

		if err := tx.Model(&models.Order{}).Where("order_id = ?", orderID).Updates(updates).Error; err != nil {
			return fmt.Errorf("CheckoutRepo.UpdateOrderStatus: update order: %w", err)
		}

		// Adjust stock count, locked count, and sold count
		if status == entity.OrderStatusPaid || status == entity.OrderStatusDelivered {
			if currentStatus == entity.OrderStatusPending {
				if err := tx.Model(&models.Product{}).Where("id = ?", order.ProductID).
					Updates(map[string]any{
						"locked_count": gorm.Expr("locked_count - ?", order.Quantity),
						"stock_count":  gorm.Expr("stock_count - ?", order.Quantity),
						"sold_count":   gorm.Expr("sold_count + ?", order.Quantity),
					}).Error; err != nil {
					return fmt.Errorf("CheckoutRepo.UpdateOrderStatus: adjust product stock: %w", err)
				}
			}
		}

		return nil
	})
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
	return withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Where("order_id = ?", orderID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("CheckoutRepo.ReleaseReservation: select order: %w", err)
		}

		if err := tx.Model(&models.Product{}).
			Where("id = ?", order.ProductID).
			UpdateColumn("locked_count", gorm.Expr("locked_count - ?", order.Quantity)).Error; err != nil {
			return fmt.Errorf("CheckoutRepo.ReleaseReservation: update product: %w", err)
		}
		return nil
	})
}
