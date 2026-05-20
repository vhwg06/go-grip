package user

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UseCase -.
type UseCase struct {
	repo repo.UserRepo
	jwt  *jwt.Manager
}

// New -.
func New(r repo.UserRepo, j *jwt.Manager) *UseCase {
	return &UseCase{
		repo: r,
		jwt:  j,
	}
}

// Register -.
func (uc *UseCase) Register(ctx context.Context, username, email, password string) (entity.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Register - bcrypt.GenerateFromPassword: %w", err)
	}

	now := time.Now().UTC()

	user := entity.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = uc.repo.Store(ctx, &user)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Register - uc.repo.Store: %w", err)
	}

	return user, nil
}

// Login -.
func (uc *UseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", entity.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", entity.ErrInvalidCredentials
	}

	token, err := uc.jwt.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("UserUseCase - Login - uc.jwt.GenerateToken: %w", err)
	}

	return token, nil
}

// GetUser -.
func (uc *UseCase) GetUser(ctx context.Context, userID string) (entity.User, error) {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - GetUser - uc.repo.GetByID: %w", err)
	}

	return user, nil
}

// List -.
func (uc *UseCase) List(ctx context.Context, limit, offset int) ([]entity.User, int, error) {
	page := entity.Pagination{Limit: limit, Offset: offset}.Normalize()

	users, total, err := uc.repo.List(ctx, repo.UserFilter{Limit: uint64(page.Limit), Offset: uint64(page.Offset)})
	if err != nil {
		return nil, 0, fmt.Errorf("UserUseCase - List - uc.repo.List: %w", err)
	}

	return users, total, nil
}

// CreateAdminUser -.
func (uc *UseCase) CreateAdminUser(ctx context.Context, _ string, username, email, password string, role entity.RoleName) (entity.User, error) {
	user, err := uc.Register(ctx, username, email, password)
	if err != nil {
		return entity.User{}, err
	}
	user.Role = role
	user.Status = entity.UserStatusActive

	return user, nil
}

// UpdateProfile -.
func (uc *UseCase) UpdateProfile(ctx context.Context, userID, displayName string) (entity.User, error) {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - UpdateProfile - uc.repo.GetByID: %w", err)
	}

	user.DisplayName = displayName
	user.Username = displayName
	user.UpdatedAt = time.Now().UTC()

	if err = uc.repo.Update(ctx, &user); err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - UpdateProfile - uc.repo.Update: %w", err)
	}

	return user, nil
}

// Lock -.
func (uc *UseCase) Lock(ctx context.Context, _ string, userID string) error {
	if err := uc.repo.SetStatus(ctx, userID, entity.UserStatusLocked); err != nil {
		return fmt.Errorf("UserUseCase - Lock - uc.repo.SetStatus: %w", err)
	}

	return nil
}

// Unlock -.
func (uc *UseCase) Unlock(ctx context.Context, _ string, userID string) error {
	if err := uc.repo.SetStatus(ctx, userID, entity.UserStatusActive); err != nil {
		return fmt.Errorf("UserUseCase - Unlock - uc.repo.SetStatus: %w", err)
	}

	return nil
}
