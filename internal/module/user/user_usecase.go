package user

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase defines the application service for User registration and account management.
type UserUseCase interface {
	Register(ctx context.Context, username, email, password string) (User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetUser(ctx context.Context, userID string) (User, error)
	List(ctx context.Context, limit, offset int) ([]User, int, error)
	CreateAdminUser(ctx context.Context, actorID, username, email, password string, role RoleName) (User, error)
	UpdateProfile(ctx context.Context, userID, displayName string) (User, error)
	Lock(ctx context.Context, actorID, userID string) error
	Unlock(ctx context.Context, actorID, userID string) error
}

type userUseCase struct {
	repo UserRepo
	jwt  *jwt.Manager
}

// NewUserUseCase constructs a new UserUseCase instance.
func NewUserUseCase(r UserRepo, j *jwt.Manager) UserUseCase {
	return &userUseCase{
		repo: r,
		jwt:  j,
	}
}

func (uc *userUseCase) Register(ctx context.Context, username, email, password string) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("UserUseCase - Register - bcrypt.GenerateFromPassword: %w", err)
	}

	now := time.Now().UTC()

	u := User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = uc.repo.Store(ctx, &u)
	if err != nil {
		return User{}, fmt.Errorf("UserUseCase - Register - uc.repo.Store: %w", err)
	}

	return u, nil
}

func (uc *userUseCase) Login(ctx context.Context, email, password string) (string, error) {
	u, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := uc.jwt.GenerateToken(u.ID)
	if err != nil {
		return "", fmt.Errorf("UserUseCase - Login - uc.jwt.GenerateToken: %w", err)
	}

	return token, nil
}

func (uc *userUseCase) GetUser(ctx context.Context, userID string) (User, error) {
	u, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("UserUseCase - GetUser - uc.repo.GetByID: %w", err)
	}

	return u, nil
}

func (uc *userUseCase) List(ctx context.Context, limit, offset int) ([]User, int, error) {
	page := pagination.Pagination{Limit: limit, Offset: offset}.Normalize()

	users, total, err := uc.repo.List(ctx, UserFilter{Limit: uint64(page.Limit), Offset: uint64(page.Offset)})
	if err != nil {
		return nil, 0, fmt.Errorf("UserUseCase - List - uc.repo.List: %w", err)
	}

	return users, total, nil
}

func (uc *userUseCase) CreateAdminUser(ctx context.Context, _ string, username, email, password string, role RoleName) (User, error) {
	u, err := uc.Register(ctx, username, email, password)
	if err != nil {
		return User{}, err
	}
	u.Role = role
	u.Status = UserStatusActive

	return u, nil
}

func (uc *userUseCase) UpdateProfile(ctx context.Context, userID, displayName string) (User, error) {
	u, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("UserUseCase - UpdateProfile - uc.repo.GetByID: %w", err)
	}

	u.DisplayName = displayName
	u.Username = displayName
	u.UpdatedAt = time.Now().UTC()

	if err = uc.repo.Update(ctx, &u); err != nil {
		return User{}, fmt.Errorf("UserUseCase - UpdateProfile - uc.repo.Update: %w", err)
	}

	return u, nil
}

func (uc *userUseCase) Lock(ctx context.Context, _ string, userID string) error {
	if err := uc.repo.SetStatus(ctx, userID, UserStatusLocked); err != nil {
		return fmt.Errorf("UserUseCase - Lock - uc.repo.SetStatus: %w", err)
	}

	return nil
}

func (uc *userUseCase) Unlock(ctx context.Context, _ string, userID string) error {
	if err := uc.repo.SetStatus(ctx, userID, UserStatusActive); err != nil {
		return fmt.Errorf("UserUseCase - Unlock - uc.repo.SetStatus: %w", err)
	}

	return nil
}
