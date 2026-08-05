package notification

import (
	"context"

	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// NotificationRepo defines persistence operations for UserNotifications and Broadcasts.
type NotificationRepo interface {
	ListUserNotifications(ctx context.Context, userID string, page pagination.Pagination) ([]UserNotification, int, error)
	ListBroadcastMessages(ctx context.Context, page pagination.Pagination) ([]BroadcastMessage, int, error)
	MarkNotificationRead(ctx context.Context, userID string, notificationID int64) error
	MarkAllRead(ctx context.Context, userID string) error
	ClearAll(ctx context.Context, userID string) error
}
