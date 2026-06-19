package persistent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm"
)

type ProfileRepo struct {
	*postgres.Postgres
}

func NewProfileRepo(pg *postgres.Postgres) *ProfileRepo {
	return &ProfileRepo{Postgres: pg}
}

var _ repo.ProfileRepository = (*ProfileRepo)(nil)

func (r *ProfileRepo) GetProfile(ctx context.Context, userID string) (entity.User, error) {
	var row models.User
	if err := r.Gorm.WithContext(ctx).Where("id = ?", userID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("ProfileRepo.GetProfile: %w", err)
	}
	return models.UserToEntity(row), nil
}

func (r *ProfileRepo) UpdateProfile(ctx context.Context, user entity.User) (entity.User, error) {
	model := map[string]any{
		"email":                         user.Email,
		"display_name":                  user.DisplayName,
		"desktop_notifications_enabled": user.DesktopNotificationsEnabled,
		"points":                        user.Points,
		"last_checkin_at":               profileTimeOrZero(user.LastCheckinAt),
		"consecutive_days":              user.ConsecutiveDays,
		"updated_at":                    time.Now().UTC(),
	}

	if err := r.Gorm.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(model).Error; err != nil {
		return entity.User{}, fmt.Errorf("ProfileRepo.UpdateProfile: %w", err)
	}

	return r.GetProfile(ctx, user.ID)
}

func (r *ProfileRepo) RecordDailyCheckin(ctx context.Context, checkin entity.DailyCheckin) error {
	model := models.DailyCheckin{
		UserID:      checkin.UserID,
		CheckinDate: checkin.CheckinDate,
		Reward:      checkin.RewardAmount,
		StreakAfter: checkin.StreakAfter,
		CreatedAt:   checkin.CreatedAt,
	}
	if err := r.Gorm.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("ProfileRepo.RecordDailyCheckin: %w", err)
	}
	return nil
}

func profileTimeOrZero(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
