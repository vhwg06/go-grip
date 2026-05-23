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

type GripOrderRepo struct {
	*postgres.Postgres
}

func NewGripOrderRepo(pg *postgres.Postgres) *GripOrderRepo {
	return &GripOrderRepo{Postgres: pg}
}

var _ repo.OrderRepository = (*GripOrderRepo)(nil)

func (r *GripOrderRepo) ListOrdersByOwner(ctx context.Context, userID, email string, page entity.Pagination) ([]entity.Order, int, error) {
	query := r.Gorm.WithContext(ctx).Model(&models.Order{})
	if userID != "" && email != "" {
		query = query.Where("user_id = ? OR email = ?", userID, email)
	} else if userID != "" {
		query = query.Where("user_id = ?", userID)
	} else {
		query = query.Where("email = ?", email)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("GripOrderRepo.ListOrdersByOwner: count: %w", err)
	}

	normalized := page.Normalize()
	var rows []models.Order
	if err := query.Order("created_at DESC").Limit(normalized.Limit).Offset(normalized.Offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("GripOrderRepo.ListOrdersByOwner: find: %w", err)
	}

	orders := make([]entity.Order, 0, len(rows))
	for _, row := range rows {
		order := models.OrderToEntity(row)
		if order.Status != entity.OrderStatusDelivered {
			order.CardKey = ""
		}
		orders = append(orders, order)
	}
	return orders, int(total), nil
}

func (r *GripOrderRepo) GetOrderByID(ctx context.Context, orderID string) (entity.Order, error) {
	var row models.Order
	if err := r.Gorm.WithContext(ctx).Where("order_id = ?", orderID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Order{}, entity.ErrOrderNotFound
		}
		return entity.Order{}, fmt.Errorf("GripOrderRepo.GetOrderByID: %w", err)
	}
	order := models.OrderToEntity(row)
	if order.Status != entity.OrderStatusDelivered {
		order.CardKey = ""
	}
	return order, nil
}

func (r *GripOrderRepo) CancelPendingOrder(ctx context.Context, actor entity.Actor, orderID string) error {
	current, err := r.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if actor.UserID != "" && current.UserID != "" && current.UserID != actor.UserID && !actor.IsAdmin {
		return entity.ErrForbidden
	}

	result := r.Gorm.WithContext(ctx).
		Model(&models.Order{}).
		Where("order_id = ? AND status = ?", orderID, string(entity.OrderStatusPending)).
		Updates(map[string]any{
			"status":     string(entity.OrderStatusCancelled),
			"updated_at": time.Now().UTC(),
		})
	if result.Error != nil {
		return fmt.Errorf("GripOrderRepo.CancelPendingOrder: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entity.ErrOrderStateConflict
	}
	return nil
}

func (r *GripOrderRepo) SubmitRefundRequest(ctx context.Context, refund entity.RefundRequest) error {
	return withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		model := models.RefundRequest{
			OrderID:       refund.OrderID,
			UserID:        refund.UserID,
			Username:      refund.Username,
			Reason:        refund.Reason,
			Status:        string(refund.Status),
			AdminUsername: refund.AdminUsername,
			AdminNote:     refund.AdminNote,
			ProcessedAt:   time.Time{},
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		}
		if err := tx.Create(&model).Error; err != nil {
			return fmt.Errorf("GripOrderRepo.SubmitRefundRequest(create): %w", err)
		}
		if err := tx.Model(&models.Order{}).
			Where("order_id = ?", refund.OrderID).
			Updates(map[string]any{
				"status":     string(entity.OrderStatusRefundPending),
				"updated_at": time.Now().UTC(),
			}).Error; err != nil {
			return fmt.Errorf("GripOrderRepo.SubmitRefundRequest(update order): %w", err)
		}
		return nil
	})
}
