package profile

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/usecase"
)

type UseCase struct {
	repo          repo.ProfileRepository
	checkinReward int
}

func New(profileRepo repo.ProfileRepository, checkinReward int) *UseCase {
	if checkinReward <= 0 {
		checkinReward = 10
	}
	return &UseCase{repo: profileRepo, checkinReward: checkinReward}
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

func (uc *UseCase) Checkin(ctx context.Context, actor entity.Actor) (entity.DailyCheckin, error) {
	user, err := uc.Get(ctx, actor)
	if err != nil {
		return entity.DailyCheckin{}, err
	}

	now := time.Now().UTC()
	if user.LastCheckinAt != nil && sameDate(*user.LastCheckinAt, now) {
		return entity.DailyCheckin{}, entity.ErrInvalidInput
	}

	streak := 1
	if user.LastCheckinAt != nil && sameDate(user.LastCheckinAt.AddDate(0, 0, 1).UTC(), now) {
		streak = user.ConsecutiveDays + 1
	}

	user.ConsecutiveDays = streak
	user.Points += uc.checkinReward
	user.LastCheckinAt = &now

	if _, err := uc.repo.UpdateProfile(ctx, user); err != nil {
		return entity.DailyCheckin{}, fmt.Errorf("ProfileUseCase.Checkin - repo.UpdateProfile: %w", err)
	}

	checkin := entity.DailyCheckin{
		UserID:       user.ID,
		CheckinDate:  now,
		RewardAmount: uc.checkinReward,
		StreakAfter:  streak,
		CreatedAt:    now,
	}
	if err := uc.repo.RecordDailyCheckin(ctx, checkin); err != nil {
		return entity.DailyCheckin{}, fmt.Errorf("ProfileUseCase.Checkin - repo.RecordDailyCheckin: %w", err)
	}

	return checkin, nil
}

func sameDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
