package notification

import (
	"context"
	"fmt"

	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// NotificationCenterUseCase defines user inbox and notification management operations.
type NotificationCenterUseCase interface {
	Inbox(ctx context.Context, actor usermodule.Actor, page pagination.Pagination) ([]UserNotification, int, error)
	UnreadCount(ctx context.Context, actor usermodule.Actor) (int, error)
	MarkRead(ctx context.Context, actor usermodule.Actor, notificationID int64) error
	MarkAllRead(ctx context.Context, actor usermodule.Actor) error
	Clear(ctx context.Context, actor usermodule.Actor) error
}

type centerUseCase struct {
	repo NotificationRepo
}

// NewNotificationCenterUseCase constructs a new NotificationCenterUseCase instance.
func NewNotificationCenterUseCase(repo NotificationRepo) NotificationCenterUseCase {
	return &centerUseCase{repo: repo}
}

func (uc *centerUseCase) Inbox(ctx context.Context, actor usermodule.Actor, page pagination.Pagination) ([]UserNotification, int, error) {
	if actor.UserID == "" {
		return nil, 0, ErrUnauthorized
	}

	personal, personalTotal, err := uc.repo.ListUserNotifications(ctx, actor.UserID, page)
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationCenter.Inbox - repo.ListUserNotifications: %w", err)
	}

	broadcasts, _, err := uc.repo.ListBroadcastMessages(ctx, page)
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationCenter.Inbox - repo.ListBroadcastMessages: %w", err)
	}

	out := make([]UserNotification, 0, len(personal)+len(broadcasts))
	out = append(out, personal...)
	for _, b := range broadcasts {
		out = append(out, UserNotification{
			ID:         b.ID,
			UserID:     actor.UserID,
			Type:       "broadcast",
			TitleKey:   b.TitleKey,
			ContentKey: b.ContentKey,
			Data:       b.Data,
			IsRead:     false,
			CreatedAt:  b.CreatedAt,
		})
	}

	return out, personalTotal + len(broadcasts), nil
}

func (uc *centerUseCase) UnreadCount(ctx context.Context, actor usermodule.Actor) (int, error) {
	if actor.UserID == "" {
		return 0, ErrUnauthorized
	}

	personal, _, err := uc.repo.ListUserNotifications(ctx, actor.UserID, pagination.Pagination{Limit: 1000, Offset: 0})
	if err != nil {
		return 0, fmt.Errorf("NotificationCenter.UnreadCount - repo.ListUserNotifications: %w", err)
	}
	count := 0
	for _, item := range personal {
		if !item.IsRead {
			count++
		}
	}

	broadcasts, _, err := uc.repo.ListBroadcastMessages(ctx, pagination.Pagination{Limit: 1000, Offset: 0})
	if err != nil {
		return 0, fmt.Errorf("NotificationCenter.UnreadCount - repo.ListBroadcastMessages: %w", err)
	}
	return count + len(broadcasts), nil
}

func (uc *centerUseCase) MarkRead(ctx context.Context, actor usermodule.Actor, notificationID int64) error {
	if actor.UserID == "" {
		return ErrUnauthorized
	}
	if err := uc.repo.MarkNotificationRead(ctx, actor.UserID, notificationID); err != nil {
		return fmt.Errorf("NotificationCenter.MarkRead - repo.MarkNotificationRead: %w", err)
	}
	return nil
}

func (uc *centerUseCase) MarkAllRead(ctx context.Context, actor usermodule.Actor) error {
	if actor.UserID == "" {
		return ErrUnauthorized
	}
	if err := uc.repo.MarkAllRead(ctx, actor.UserID); err != nil {
		return fmt.Errorf("NotificationCenter.MarkAllRead - repo.MarkAllRead: %w", err)
	}
	return nil
}

func (uc *centerUseCase) Clear(ctx context.Context, actor usermodule.Actor) error {
	if actor.UserID == "" {
		return ErrUnauthorized
	}
	if err := uc.repo.ClearAll(ctx, actor.UserID); err != nil {
		return fmt.Errorf("NotificationCenter.Clear - repo.ClearAll: %w", err)
	}
	return nil
}
