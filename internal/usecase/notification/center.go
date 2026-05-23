package notification

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/usecase"
)

type CenterUseCase struct {
	repo repo.NotificationRepository
}

func NewCenter(notificationRepo repo.NotificationRepository) *CenterUseCase {
	return &CenterUseCase{repo: notificationRepo}
}

var _ usecase.NotificationCenter = (*CenterUseCase)(nil)

func (uc *CenterUseCase) Inbox(ctx context.Context, actor entity.Actor, page entity.Pagination) ([]entity.UserNotification, int, error) {
	if actor.UserID == "" {
		return nil, 0, entity.ErrUnauthorized
	}

	personal, personalTotal, err := uc.repo.ListUserNotifications(ctx, actor.UserID, page)
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationCenter.Inbox - repo.ListUserNotifications: %w", err)
	}

	broadcasts, _, err := uc.repo.ListBroadcastMessages(ctx, page)
	if err != nil {
		return nil, 0, fmt.Errorf("NotificationCenter.Inbox - repo.ListBroadcastMessages: %w", err)
	}

	out := make([]entity.UserNotification, 0, len(personal)+len(broadcasts))
	out = append(out, personal...)
	for _, b := range broadcasts {
		out = append(out, entity.UserNotification{
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

func (uc *CenterUseCase) UnreadCount(ctx context.Context, actor entity.Actor) (int, error) {
	if actor.UserID == "" {
		return 0, entity.ErrUnauthorized
	}

	personal, _, err := uc.repo.ListUserNotifications(ctx, actor.UserID, entity.Pagination{Limit: 1000, Offset: 0})
	if err != nil {
		return 0, fmt.Errorf("NotificationCenter.UnreadCount - repo.ListUserNotifications: %w", err)
	}
	count := 0
	for _, item := range personal {
		if !item.IsRead {
			count++
		}
	}

	broadcasts, _, err := uc.repo.ListBroadcastMessages(ctx, entity.Pagination{Limit: 1000, Offset: 0})
	if err != nil {
		return 0, fmt.Errorf("NotificationCenter.UnreadCount - repo.ListBroadcastMessages: %w", err)
	}
	return count + len(broadcasts), nil
}

func (uc *CenterUseCase) MarkRead(ctx context.Context, actor entity.Actor, notificationID int64) error {
	if actor.UserID == "" {
		return entity.ErrUnauthorized
	}
	if err := uc.repo.MarkNotificationRead(ctx, actor.UserID, notificationID); err != nil {
		return fmt.Errorf("NotificationCenter.MarkRead - repo.MarkNotificationRead: %w", err)
	}
	return nil
}

func (uc *CenterUseCase) MarkAllRead(ctx context.Context, actor entity.Actor) error {
	if actor.UserID == "" {
		return entity.ErrUnauthorized
	}
	if err := uc.repo.MarkAllRead(ctx, actor.UserID); err != nil {
		return fmt.Errorf("NotificationCenter.MarkAllRead - repo.MarkAllRead: %w", err)
	}
	return nil
}

func (uc *CenterUseCase) Clear(ctx context.Context, actor entity.Actor) error {
	if actor.UserID == "" {
		return entity.ErrUnauthorized
	}
	if err := uc.repo.ClearAll(ctx, actor.UserID); err != nil {
		return fmt.Errorf("NotificationCenter.Clear - repo.ClearAll: %w", err)
	}
	return nil
}
