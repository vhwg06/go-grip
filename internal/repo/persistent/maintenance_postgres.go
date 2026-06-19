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

		if err := tx.WithContext(ctx).Exec(`
			UPDATE products p
			SET locked_count = p.locked_count - o.quantity
			FROM orders o
			WHERE p.id::text = o.product_id
			  AND o.order_id IN ?
		`, expiredOrders).Error; err != nil {
			return fmt.Errorf("MaintenanceRepo.CancelExpiredPendingOrders(release locked count): %w", err)
		}

		return nil
	})
}

func (r *MaintenanceRepo) SyncProductAggregates(ctx context.Context) error {
	return nil
}
