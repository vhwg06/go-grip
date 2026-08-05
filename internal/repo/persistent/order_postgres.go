package persistent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm"
)

type GripOrderRepo struct {
	*postgres.Postgres
}

func NewGripOrderRepo(pg *postgres.Postgres) *GripOrderRepo {
	return &GripOrderRepo{Postgres: pg}
}

func (r *GripOrderRepo) ListOrdersByOwner(ctx context.Context, userID, email string, page pagination.Pagination) ([]ordermodule.Order, int, error) {
	if r.Postgres == nil || r.Gorm == nil {
		return nil, 0, nil
	}
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

	orders := make([]ordermodule.Order, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, models.OrderToModule(row))
	}
	return orders, int(total), nil
}

func (r *GripOrderRepo) GetOrderByID(ctx context.Context, orderID string) (ordermodule.Order, error) {
	if r.Postgres == nil || r.Gorm == nil {
		return ordermodule.Order{}, ordermodule.ErrNotFound
	}
	var row models.Order
	if err := r.Gorm.WithContext(ctx).Where("order_id = ?", orderID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ordermodule.Order{}, ordermodule.ErrNotFound
		}
		return ordermodule.Order{}, fmt.Errorf("GripOrderRepo.GetOrderByID: %w", err)
	}
	return models.OrderToModule(row), nil
}

func (r *GripOrderRepo) CancelPendingOrder(ctx context.Context, actor usermodule.Actor, orderID string) error {
	if r.Postgres == nil || r.Gorm == nil {
		return nil
	}
	current, err := r.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if actor.UserID != "" && current.UserID != "" && current.UserID != actor.UserID && !actor.IsAdmin {
		return ordermodule.ErrForbidden
	}

	result := r.Gorm.WithContext(ctx).
		Model(&models.Order{}).
		Where("order_id = ? AND status = ?", orderID, string(ordermodule.OrderStatusPending)).
		Updates(map[string]any{
			"status":     string(ordermodule.OrderStatusCancelled),
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

func (r *GripOrderRepo) SubmitRefundRequest(ctx context.Context, refund *ordermodule.RefundRequest) error {
	if r.Postgres == nil || r.Gorm == nil {
		return nil
	}
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
		refund.ID = model.ID
		if err := tx.Model(&models.Order{}).
			Where("order_id = ?", refund.OrderID).
			Updates(map[string]any{
				"status":     string(ordermodule.OrderStatusRefundPending),
				"updated_at": time.Now().UTC(),
			}).Error; err != nil {
			return fmt.Errorf("GripOrderRepo.SubmitRefundRequest(update order): %w", err)
		}
		return nil
	})
}
