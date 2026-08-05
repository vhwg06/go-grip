package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/module/user"
	"github.com/stretchr/testify/require"
)

type maintenanceRepoStub struct {
	cancelExpiredFunc func(ctx context.Context, olderThan time.Time) error
	syncAggregatesFunc func(ctx context.Context) error
}

func (s *maintenanceRepoStub) CancelExpiredPendingOrders(ctx context.Context, olderThan time.Time) error {
	if s.cancelExpiredFunc != nil {
		return s.cancelExpiredFunc(ctx, olderThan)
	}
	return nil
}

func (s *maintenanceRepoStub) SyncProductAggregates(ctx context.Context) error {
	if s.syncAggregatesFunc != nil {
		return s.syncAggregatesFunc(ctx)
	}
	return nil
}

func TestMaintenanceUseCase_CancelExpiredPendingOrders(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		called := false
		stub := &maintenanceRepoStub{
			cancelExpiredFunc: func(ctx context.Context, olderThan time.Time) error {
				called = true
				require.WithinDuration(t, time.Now().UTC().Add(-10*time.Minute), olderThan, 2*time.Second)
				return nil
			},
		}
		uc := user.NewMaintenanceUseCase(stub, 10*time.Minute)
		err := uc.CancelExpiredPendingOrders(context.Background())
		require.NoError(t, err)
		require.True(t, called)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("db failure")
		stub := &maintenanceRepoStub{
			cancelExpiredFunc: func(ctx context.Context, olderThan time.Time) error {
				return expectedErr
			},
		}
		uc := user.NewMaintenanceUseCase(stub, 5*time.Minute)
		err := uc.CancelExpiredPendingOrders(context.Background())
		require.ErrorIs(t, err, expectedErr)
	})
}

func TestMaintenanceUseCase_SyncProductAggregates(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		called := false
		stub := &maintenanceRepoStub{
			syncAggregatesFunc: func(ctx context.Context) error {
				called = true
				return nil
			},
		}
		uc := user.NewMaintenanceUseCase(stub, 0)
		err := uc.SyncProductAggregates(context.Background())
		require.NoError(t, err)
		require.True(t, called)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("sync error")
		stub := &maintenanceRepoStub{
			syncAggregatesFunc: func(ctx context.Context) error {
				return expectedErr
			},
		}
		uc := user.NewMaintenanceUseCase(stub, 0)
		err := uc.SyncProductAggregates(context.Background())
		require.ErrorIs(t, err, expectedErr)
	})
}
