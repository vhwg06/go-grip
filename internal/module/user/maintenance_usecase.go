package user

import (
	"context"
	"fmt"
	"time"
)

// MaintenanceRepository defines the consumer-owned persistence contract required for background maintenance operations.
type MaintenanceRepository interface {
	// CancelExpiredPendingOrders marks pending orders older than the specified timestamp as cancelled.
	CancelExpiredPendingOrders(ctx context.Context, olderThan time.Time) error
	// SyncProductAggregates recalculates and updates aggregated metrics for products.
	SyncProductAggregates(ctx context.Context) error
}

// MaintenanceUseCase coordinates scheduled system maintenance tasks such as order expiration and product aggregate synchronization.
type MaintenanceUseCase struct {
	repo          MaintenanceRepository
	pendingWindow time.Duration
}

// NewMaintenanceUseCase constructs a new MaintenanceUseCase instance with the provided repository and pending order timeout window.
func NewMaintenanceUseCase(repo MaintenanceRepository, pendingWindow time.Duration) *MaintenanceUseCase {
	if pendingWindow <= 0 {
		pendingWindow = 5 * time.Minute
	}
	return &MaintenanceUseCase{
		repo:          repo,
		pendingWindow: pendingWindow,
	}
}

// CancelExpiredPendingOrders identifies and cancels pending orders created prior to the configured pending window duration.
func (uc *MaintenanceUseCase) CancelExpiredPendingOrders(ctx context.Context) error {
	olderThan := time.Now().UTC().Add(-uc.pendingWindow)
	if err := uc.repo.CancelExpiredPendingOrders(ctx, olderThan); err != nil {
		return fmt.Errorf("MaintenanceUseCase.CancelExpiredPendingOrders: %w", err)
	}
	return nil
}

// SyncProductAggregates triggers a sync operation to refresh stale product rating and sales aggregate data.
func (uc *MaintenanceUseCase) SyncProductAggregates(ctx context.Context) error {
	if err := uc.repo.SyncProductAggregates(ctx); err != nil {
		return fmt.Errorf("MaintenanceUseCase.SyncProductAggregates: %w", err)
	}
	return nil
}
