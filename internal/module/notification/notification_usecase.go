package notification

import "context"

// NotificationUseCase handles outbound notification dispatch.
type NotificationUseCase interface {
	Dispatch(ctx context.Context, notification Notification) error
}

type notificationUseCase struct {
	enabled bool
}

// NewNotificationUseCase constructs a new NotificationUseCase instance.
func NewNotificationUseCase(enabled bool) NotificationUseCase {
	return &notificationUseCase{enabled: enabled}
}

func (uc *notificationUseCase) Dispatch(ctx context.Context, notification Notification) error {
	_ = ctx
	_ = notification
	if !uc.enabled {
		return nil
	}
	return nil
}
