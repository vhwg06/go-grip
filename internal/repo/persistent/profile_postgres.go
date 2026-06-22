package persistent

import (
	"context"
	"errors"
	"fmt"

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
	}

	if err := r.Gorm.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(model).Error; err != nil {
		return entity.User{}, fmt.Errorf("ProfileRepo.UpdateProfile: %w", err)
	}

	return r.GetProfile(ctx, user.ID)
}

func (r *ProfileRepo) GetGorm() *gorm.DB {
	return r.Gorm
}
