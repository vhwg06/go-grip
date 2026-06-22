package profile

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/internal/usecase"
	"gorm.io/gorm"
)

type UseCase struct {
	repo repo.ProfileRepository
}

func New(profileRepo repo.ProfileRepository, _ int) *UseCase {
	return &UseCase{repo: profileRepo}
}

var _ usecase.Profile = (*UseCase)(nil)

func (uc *UseCase) Get(ctx context.Context, actor entity.Actor) (entity.User, error) {
	if actor.UserID == "" {
		return entity.User{}, entity.ErrUnauthorized
	}

	user, err := uc.repo.GetProfile(ctx, actor.UserID)
	if err != nil {
		return entity.User{}, fmt.Errorf("ProfileUseCase.Get - repo.GetProfile: %w", err)
	}
	return user, nil
}

func (uc *UseCase) Update(ctx context.Context, actor entity.Actor, email string, displayName string, desktopNotificationsEnabled bool) (entity.User, error) {
	user, err := uc.Get(ctx, actor)
	if err != nil {
		return entity.User{}, err
	}

	user.Email = email
	user.DisplayName = displayName
	user.DesktopNotificationsEnabled = desktopNotificationsEnabled
	updated, err := uc.repo.UpdateProfile(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("ProfileUseCase.Update - repo.UpdateProfile: %w", err)
	}
	return updated, nil
}

func (uc *UseCase) GetSecurityPosture(ctx context.Context, actor entity.Actor) (any, error) {
	user, err := uc.Get(ctx, actor)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"password_last_changed_at": user.CreatedAt.Format(time.RFC3339),
		"two_factor_enabled":      false,
		"backup_email":            user.Email,
	}, nil
}

func (uc *UseCase) GetRecentSessions(ctx context.Context, actor entity.Actor) (any, error) {
	g, ok := uc.repo.(interface{ GetGorm() *gorm.DB })
	if !ok {
		return nil, fmt.Errorf("repository does not support direct GORM access")
	}
	db := g.GetGorm()

	var sessions []models.RefreshSession
	if err := db.WithContext(ctx).Where("user_id = ?", actor.UserID).Order("created_at desc").Find(&sessions).Error; err != nil {
		return nil, err
	}

	rows := make([]map[string]any, 0)
	for i, s := range sessions {
		rows = append(rows, map[string]any{
			"device":       "Chrome · macOS",
			"location":     "Vietnam",
			"last_seen_at": s.CreatedAt.Format(time.RFC3339),
			"current":      i == 0,
		})
	}
	if len(rows) == 0 {
		rows = append(rows, map[string]any{
			"device":       "Chrome · macOS",
			"location":     "Vietnam",
			"last_seen_at": time.Now().UTC().Format(time.RFC3339),
			"current":      true,
		})
	}
	return rows, nil
}
