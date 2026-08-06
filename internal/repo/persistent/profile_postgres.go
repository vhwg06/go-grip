package persistent

import (
	"context"
	"errors"
	"fmt"

	usermodule "github.com/evrone/go-clean-template/internal/module/user"
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

func (r *ProfileRepo) GetProfile(ctx context.Context, userID string) (usermodule.User, error) {
	if r.Postgres == nil || r.Gorm == nil {
		return usermodule.User{}, usermodule.ErrNotFound
	}
	var row models.User
	if err := r.Gorm.WithContext(ctx).Where("id = ?", userID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return usermodule.User{}, usermodule.ErrNotFound
		}
		return usermodule.User{}, fmt.Errorf("ProfileRepo.GetProfile: %w", err)
	}
	return models.UserToModule(row), nil
}

func (r *ProfileRepo) UpdateProfile(ctx context.Context, u usermodule.User) (usermodule.User, error) {
	if r.Postgres == nil || r.Gorm == nil {
		return u, nil
	}
	model := map[string]any{
		"desktop_notifications_enabled": u.DesktopNotificationsEnabled,
	}
	if u.Email != "" {
		model["email"] = u.Email
	}
	if u.DisplayName != "" {
		model["display_name"] = u.DisplayName
	}

	if err := r.Gorm.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", u.ID).
		Updates(model).Error; err != nil {
		return usermodule.User{}, fmt.Errorf("ProfileRepo.UpdateProfile: %w", err)
	}

	return r.GetProfile(ctx, u.ID)
}

func (r *ProfileRepo) GetRecentSessions(ctx context.Context, userID string) ([]usermodule.RefreshSession, error) {
	if r.Postgres == nil || r.Gorm == nil {
		return nil, nil
	}
	var sessions []models.RefreshSession
	if err := r.Gorm.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&sessions).Error; err != nil {
		return nil, err
	}
	res := make([]usermodule.RefreshSession, len(sessions))
	for i, s := range sessions {
		res[i] = usermodule.RefreshSession{
			ID:        s.ID,
			UserID:    s.UserID,
			TokenID:   s.TokenID,
			ExpiresAt: s.ExpiresAt,
			RevokedAt: s.RevokedAt,
			CreatedAt: s.CreatedAt,
		}
	}
	return res, nil
}

func (r *ProfileRepo) GetGorm() *gorm.DB {
	if r.Postgres == nil {
		return nil
	}
	return r.Gorm
}
