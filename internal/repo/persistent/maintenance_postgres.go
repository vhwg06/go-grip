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

type MaintenanceRepo struct {
	*postgres.Postgres
}

func NewMaintenanceRepo(pg *postgres.Postgres) *MaintenanceRepo {
	return &MaintenanceRepo{Postgres: pg}
}

var _ repo.MaintenanceRepository = (*MaintenanceRepo)(nil)

func (r *MaintenanceRepo) CancelExpiredPendingOrders(ctx context.Context, olderThan time.Time) error {
	return withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		var expiredOrders []string
		if err := tx.WithContext(ctx).
			Model(&models.Order{}).
			Where("status = ?", string(entity.OrderStatusPending)).
			Where("created_at < ?", olderThan).
			Pluck("order_id", &expiredOrders).Error; err != nil {
			return fmt.Errorf("MaintenanceRepo.CancelExpiredPendingOrders(pluck): %w", err)
		}

		if len(expiredOrders) == 0 {
			return nil
		}

		if err := tx.WithContext(ctx).
			Model(&models.Order{}).
			Where("order_id IN ?", expiredOrders).
			Updates(map[string]any{
				"status":     string(entity.OrderStatusCancelled),
				"updated_at": time.Now().UTC(),
			}).Error; err != nil {
			return fmt.Errorf("MaintenanceRepo.CancelExpiredPendingOrders(update orders): %w", err)
		}

		if err := tx.WithContext(ctx).
			Model(&models.Card{}).
			Where("reserved_order_id IN ?", expiredOrders).
			Updates(map[string]any{
				"reserved_order_id": "",
				"reserved_at":       time.Time{},
			}).Error; err != nil {
			return fmt.Errorf("MaintenanceRepo.CancelExpiredPendingOrders(release cards): %w", err)
		}

		return nil
	})
}


func (r *MaintenanceRepo) SyncProductAggregates(ctx context.Context) error {
	sql := `
		UPDATE products p
		SET
			stock_count = COALESCE(src.stock_count, 0),
			locked_count = COALESCE(src.locked_count, 0),
			sold_count = COALESCE(src.sold_count, 0),
			updated_at = NOW()
		FROM (
			SELECT
				product_id,
				COUNT(*) FILTER (WHERE is_used = FALSE AND (reserved_order_id = '' OR reserved_order_id IS NULL)) AS stock_count,
				COUNT(*) FILTER (WHERE is_used = FALSE AND reserved_order_id <> '') AS locked_count,
				COUNT(*) FILTER (WHERE is_used = TRUE) AS sold_count
			FROM cards
			GROUP BY product_id
		) src
		WHERE p.id::text = src.product_id
	`
	if err := r.Gorm.WithContext(ctx).Exec(sql).Error; err != nil {
		return fmt.Errorf("MaintenanceRepo.SyncProductAggregates: %w", err)
	}
	return nil
}
