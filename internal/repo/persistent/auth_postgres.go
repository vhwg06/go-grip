package persistent

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepo struct {
	*postgres.Postgres
}

func NewAuthRepo(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{Postgres: pg}
}

var _ repo.AuthRepository = (*AuthRepo)(nil)

func (r *AuthRepo) GetUserByID(ctx context.Context, userID string) (entity.User, error) {
	var row models.User
	if err := r.Gorm.WithContext(ctx).Where("id = ?", userID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("AuthRepo.GetUserByID: %w", err)
	}
	return models.UserToEntity(row), nil
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var row models.User
	if err := r.Gorm.WithContext(ctx).Where("LOWER(email) = LOWER(?)", email).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("AuthRepo.GetUserByEmail: %w", err)
	}
	return models.UserToEntity(row), nil
}

func (r *AuthRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	var row models.User
	if err := r.Gorm.WithContext(ctx).Where("LOWER(username) = LOWER(?)", username).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("AuthRepo.GetUserByUsername: %w", err)
	}
	return models.UserToEntity(row), nil
}

func (r *AuthRepo) UpsertUser(ctx context.Context, user entity.User) (entity.User, error) {
	now := time.Now().UTC()

	var resolved models.User
	err := withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		lookup := tx.Model(&models.User{})

		if user.Provider != "" && user.ProviderID != "" {
			if err := lookup.Where("provider = ? AND provider_id = ?", user.Provider, user.ProviderID).First(&resolved).Error; err == nil {
				return r.updateExistingUser(tx, &resolved, user, now)
			}
		}

		if strings.TrimSpace(user.Email) != "" {
			if err := lookup.Where("LOWER(email) = LOWER(?)", user.Email).First(&resolved).Error; err == nil {
				return r.updateExistingUser(tx, &resolved, user, now)
			}
		}

		if strings.TrimSpace(user.Username) != "" {
			if err := lookup.Where("LOWER(username) = LOWER(?)", user.Username).First(&resolved).Error; err == nil {
				return r.updateExistingUser(tx, &resolved, user, now)
			}
		}

		if user.ID == "" {
			user.ID = uuid.NewString()
		}
		if user.Status == "" {
			user.Status = entity.UserStatusActive
		}
		model := models.EntityToUser(user)
		model.Status = string(user.Status)
		model.CreatedAt = now
		model.UpdatedAt = now
		model.LastLoginAt = now

		if model.Email == "" {
			model.Email = strings.ToLower(model.Username) + "@placeholder.local"
		}

		if err := tx.Create(&model).Error; err != nil {
			return fmt.Errorf("AuthRepo.UpsertUser: create user: %w", err)
		}
		resolved = model
		return nil
	})
	if err != nil {
		return entity.User{}, err
	}

	return models.UserToEntity(resolved), nil
}

func (r *AuthRepo) updateExistingUser(tx *gorm.DB, existing *models.User, incoming entity.User, now time.Time) error {
	if incoming.Provider != "" {
		existing.Provider = incoming.Provider
	}
	if incoming.ProviderID != "" {
		existing.ProviderID = incoming.ProviderID
	}
	if incoming.Username != "" {
		existing.Username = incoming.Username
	}
	if incoming.Email != "" {
		existing.Email = incoming.Email
	}
	existing.LastLoginAt = now
	existing.UpdatedAt = now

	if err := tx.Save(existing).Error; err != nil {
		return fmt.Errorf("AuthRepo.updateExistingUser: %w", err)
	}
	return nil
}

func (r *AuthRepo) StoreRefreshSession(ctx context.Context, session entity.RefreshSession) error {
	model := models.RefreshSession{
		ID:        session.ID,
		UserID:    session.UserID,
		TokenID:   session.TokenID,
		ExpiresAt: session.ExpiresAt,
		RevokedAt: session.RevokedAt,
		CreatedAt: session.CreatedAt,
	}
	if model.ID == "" {
		model.ID = session.TokenID
	}
	if err := r.Gorm.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("AuthRepo.StoreRefreshSession: %w", err)
	}
	return nil
}

func (r *AuthRepo) RevokeRefreshSession(ctx context.Context, tokenID string) error {
	err := r.Gorm.WithContext(ctx).
		Model(&models.RefreshSession{}).
		Where("token_id = ?", tokenID).
		Update("revoked_at", time.Now().UTC()).
		Error
	if err != nil {
		return fmt.Errorf("AuthRepo.RevokeRefreshSession: %w", err)
	}
	return nil
}

func (r *AuthRepo) GetRefreshSession(ctx context.Context, tokenID string) (entity.RefreshSession, error) {
	var model models.RefreshSession
	if err := r.Gorm.WithContext(ctx).Where("token_id = ?", tokenID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.RefreshSession{}, entity.ErrUnauthorized
		}
		return entity.RefreshSession{}, fmt.Errorf("AuthRepo.GetRefreshSession: %w", err)
	}

	return entity.RefreshSession{
		ID:        model.ID,
		UserID:    model.UserID,
		TokenID:   model.TokenID,
		ExpiresAt: model.ExpiresAt,
		RevokedAt: model.RevokedAt,
		CreatedAt: model.CreatedAt,
	}, nil
}
