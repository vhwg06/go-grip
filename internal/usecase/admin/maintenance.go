package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/usecase"
)

type MaintenanceUseCase struct {
	repo              repo.MaintenanceRepository
	pendingWindow     time.Duration
	expiredCardsAfter time.Duration
}

func NewMaintenance(maintenanceRepo repo.MaintenanceRepository, pendingWindow time.Duration) *MaintenanceUseCase {
	if pendingWindow <= 0 {
		pendingWindow = 5 * time.Minute
	}
	return &MaintenanceUseCase{
		repo:          maintenanceRepo,
		pendingWindow: pendingWindow,
	}
}

var _ usecase.Maintenance = (*MaintenanceUseCase)(nil)

func (uc *MaintenanceUseCase) CancelExpiredPendingOrders(ctx context.Context) error {
	olderThan := time.Now().UTC().Add(-uc.pendingWindow)
	if err := uc.repo.CancelExpiredPendingOrders(ctx, olderThan); err != nil {
		return fmt.Errorf("MaintenanceUseCase.CancelExpiredPendingOrders: %w", err)
	}
	return nil
}

func (uc *MaintenanceUseCase) CleanupExpiredCards(ctx context.Context) error {
	if err := uc.repo.CleanupExpiredCards(ctx, time.Now().UTC()); err != nil {
		return fmt.Errorf("MaintenanceUseCase.CleanupExpiredCards: %w", err)
	}
	return nil
}

func (uc *MaintenanceUseCase) SyncProductAggregates(ctx context.Context) error {
	if err := uc.repo.SyncProductAggregates(ctx); err != nil {
		return fmt.Errorf("MaintenanceUseCase.SyncProductAggregates: %w", err)
	}
	return nil
}
