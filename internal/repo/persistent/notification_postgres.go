package persistent

import (
	"context"
	"fmt"
	"time"

	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm/clause"
)

type NotificationRepo struct {
	*postgres.Postgres
}

func NewNotificationRepo(pg *postgres.Postgres) *NotificationRepo {
	return &NotificationRepo{Postgres: pg}
}

var _ notificationmodule.NotificationRepo = (*NotificationRepo)(nil)

func (r *NotificationRepo) ListUserNotifications(ctx context.Context, userID string, page pagination.Pagination) ([]notificationmodule.UserNotification, int, error) {
	query := r.Gorm.WithContext(ctx).Model(&models.UserNotification{}).Where("user_id = ?", userID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo.ListUserNotifications(count): %w", err)
	}
	normalized := page.Normalize()
	var rows []models.UserNotification
	if err := query.Order("created_at DESC").Limit(normalized.Limit).Offset(normalized.Offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo.ListUserNotifications(find): %w", err)
	}
	items := make([]notificationmodule.UserNotification, 0, len(rows))
	for _, row := range rows {
		items = append(items, notificationmodule.UserNotification{
			ID:         row.ID,
			UserID:     row.UserID,
			Type:       row.Type,
			TitleKey:   row.TitleKey,
			ContentKey: row.ContentKey,
			Data:       row.Data,
			IsRead:     row.IsRead,
			CreatedAt:  row.CreatedAt,
		})
	}
	return items, int(total), nil
}

func (r *NotificationRepo) ListBroadcastMessages(ctx context.Context, page pagination.Pagination) ([]notificationmodule.BroadcastMessage, int, error) {
	query := r.Gorm.WithContext(ctx).Model(&models.BroadcastMessage{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo.ListBroadcastMessages(count): %w", err)
	}
	normalized := page.Normalize()
	var rows []models.BroadcastMessage
	if err := query.Order("created_at DESC").Limit(normalized.Limit).Offset(normalized.Offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("NotificationRepo.ListBroadcastMessages(find): %w", err)
	}
	items := make([]notificationmodule.BroadcastMessage, 0, len(rows))
	for _, row := range rows {
		items = append(items, notificationmodule.BroadcastMessage{
			ID:         row.ID,
			TitleKey:   row.TitleKey,
			ContentKey: row.ContentKey,
			Data:       row.Data,
			Sender:     row.Sender,
			CreatedAt:  row.CreatedAt,
		})
	}
	return items, int(total), nil
}

func (r *NotificationRepo) MarkNotificationRead(ctx context.Context, userID string, notificationID int64) error {
	if err := r.Gorm.WithContext(ctx).Model(&models.UserNotification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]any{"is_read": true}).Error; err != nil {
		return fmt.Errorf("NotificationRepo.MarkNotificationRead: %w", err)
	}
	return nil
}

func (r *NotificationRepo) MarkAllRead(ctx context.Context, userID string) error {
	if err := r.Gorm.WithContext(ctx).Model(&models.UserNotification{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{"is_read": true}).Error; err != nil {
		return fmt.Errorf("NotificationRepo.MarkAllRead: %w", err)
	}
	return nil
}

func (r *NotificationRepo) ClearAll(ctx context.Context, userID string) error {
	if err := r.Gorm.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.UserNotification{}).Error; err != nil {
		return fmt.Errorf("NotificationRepo.ClearAll: %w", err)
	}

	// Store a marker record in broadcast_reads to avoid resurrecting old broadcasts.
	marker := models.BroadcastRead{
		MessageID: 0,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}
	if err := r.Gorm.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "message_id"}, {Name: "user_id"}},
			DoNothing: true,
		}).
		Create(&marker).Error; err != nil {
		return fmt.Errorf("NotificationRepo.ClearAll(marker): %w", err)
	}
	return nil
}
