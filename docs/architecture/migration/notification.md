# Notification Module Migration Manifest

This document records the inventory of all symbols, contracts, adapters, delivery endpoints, consumers, and migration risks for the `Notification` business module.

---

## 1. Owned Symbols
- **Entities**: `Notification`, `UserNotification`, `BroadcastMessage`, `BroadcastRead`, `AdminMessage`
- **Errors**: `ErrInvalidInput`, `ErrUnauthorized`
- **Use Cases**: `NotificationUseCase`, `NotificationCenterUseCase`
- **Repository Ports**: `NotificationRepo`

---

## 2. Ports & Interfaces
```go
package notification

import (
	"context"
	"time"

	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

type Notification struct {
	Channel string `json:"channel"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type UserNotification struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	Type       string    `json:"type"`
	TitleKey   string    `json:"title_key"`
	ContentKey string    `json:"content_key"`
	Data       string    `json:"data"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

type NotificationRepo interface {
	ListUserNotifications(ctx context.Context, userID string, page pagination.Pagination) ([]UserNotification, int, error)
	MarkAllAsRead(ctx context.Context, userID string) error
	StoreNotification(ctx context.Context, item UserNotification) (UserNotification, error)
}
```

---

## 3. Infrastructure & Delivery Consumers
- `internal/repo/persistent/notification_postgres.go`
- `internal/controller/restapi/v1/notification/`
- `internal/app/app.go`
