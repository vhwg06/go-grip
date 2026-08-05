package user

import (
	"context"
	"fmt"
	"time"
)

// ProfileUseCase defines the application service for User profile management.
type ProfileUseCase interface {
	Get(ctx context.Context, actor Actor) (User, error)
	Update(ctx context.Context, actor Actor, email string, displayName string, desktopNotificationsEnabled bool) (User, error)
	GetSecurityPosture(ctx context.Context, actor Actor) (any, error)
	GetRecentSessions(ctx context.Context, actor Actor) (any, error)
}

type profileUseCase struct {
	repo ProfileRepo
}

// NewProfileUseCase constructs a new ProfileUseCase instance.
func NewProfileUseCase(profileRepo ProfileRepo) ProfileUseCase {
	return &profileUseCase{repo: profileRepo}
}

func (uc *profileUseCase) Get(ctx context.Context, actor Actor) (User, error) {
	if actor.UserID == "" {
		return User{}, ErrUnauthorized
	}

	u, err := uc.repo.GetProfile(ctx, actor.UserID)
	if err != nil {
		return User{}, fmt.Errorf("ProfileUseCase.Get - repo.GetProfile: %w", err)
	}
	return u, nil
}

func (uc *profileUseCase) Update(ctx context.Context, actor Actor, email string, displayName string, desktopNotificationsEnabled bool) (User, error) {
	u, err := uc.Get(ctx, actor)
	if err != nil {
		return User{}, err
	}

	u.Email = email
	u.DisplayName = displayName
	u.DesktopNotificationsEnabled = desktopNotificationsEnabled
	updated, err := uc.repo.UpdateProfile(ctx, u)
	if err != nil {
		return User{}, fmt.Errorf("ProfileUseCase.Update - repo.UpdateProfile: %w", err)
	}
	return updated, nil
}

func (uc *profileUseCase) GetSecurityPosture(ctx context.Context, actor Actor) (any, error) {
	u, err := uc.Get(ctx, actor)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"password_last_changed_at": u.CreatedAt.Format(time.RFC3339),
		"two_factor_enabled":      false,
		"backup_email":            u.Email,
	}, nil
}

func (uc *profileUseCase) GetRecentSessions(ctx context.Context, actor Actor) (any, error) {
	sessions, err := uc.repo.GetRecentSessions(ctx, actor.UserID)
	if err != nil {
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
