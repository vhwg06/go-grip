package notification

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
)

type UseCase struct {
	enabled bool
}

func New(enabled bool) *UseCase {
	return &UseCase{enabled: enabled}
}

func (uc *UseCase) Dispatch(ctx context.Context, notification entity.Notification) error {
	_ = ctx
	_ = notification
	if !uc.enabled {
		return nil
	}
	return nil
}
