package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthUseCase defines the application service for authentication and session management.
type AuthUseCase interface {
	Login(ctx context.Context, email, password string) (User, string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, actor Actor, refreshToken string) error
	Me(ctx context.Context, actor Actor) (User, error)
}

type authUseCase struct {
	repo       AuthRepo
	jwtManager *jwt.Manager
	refreshTTL time.Duration
	adminUsers map[string]struct{}
}

// NewAuthUseCase constructs a new AuthUseCase instance.
func NewAuthUseCase(authRepo AuthRepo, jwtManager *jwt.Manager, refreshTTL time.Duration, adminUsersCSV string) AuthUseCase {
	adminUsers := make(map[string]struct{})
	for user := range strings.SplitSeq(adminUsersCSV, ",") {
		trimmed := strings.ToLower(strings.TrimSpace(user))
		if trimmed == "" {
			continue
		}
		adminUsers[trimmed] = struct{}{}
	}

	return &authUseCase{
		repo:       authRepo,
		jwtManager: jwtManager,
		refreshTTL: refreshTTL,
		adminUsers: adminUsers,
	}
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (User, string, string, error) {
	u, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, "", "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return User{}, "", "", ErrInvalidCredentials
	}

	u.IsAdmin = u.IsAdmin || uc.isAdminUsername(u.Username)
	accessToken, refreshToken, err := uc.issueTokenPair(ctx, u.ID)
	if err != nil {
		return User{}, "", "", err
	}

	return u, accessToken, refreshToken, nil
}

func (uc *authUseCase) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	if refreshToken == "" {
		return "", "", ErrInvalidInput
	}

	session, err := uc.repo.GetRefreshSession(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("AuthUseCase.Refresh - repo.GetRefreshSession: %w", err)
	}
	if !session.RevokedAt.IsZero() || session.ExpiresAt.Before(time.Now().UTC()) {
		return "", "", ErrUnauthorized
	}

	if err := uc.repo.RevokeRefreshSession(ctx, refreshToken); err != nil {
		return "", "", fmt.Errorf("AuthUseCase.Refresh - repo.RevokeRefreshSession: %w", err)
	}

	access, rotatedRefresh, err := uc.issueTokenPair(ctx, session.UserID)
	if err != nil {
		return "", "", err
	}
	return access, rotatedRefresh, nil
}

func (uc *authUseCase) Logout(ctx context.Context, actor Actor, refreshToken string) error {
	if actor.UserID == "" {
		return ErrUnauthorized
	}
	if refreshToken == "" {
		return ErrInvalidInput
	}
	if err := uc.repo.RevokeRefreshSession(ctx, refreshToken); err != nil {
		return fmt.Errorf("AuthUseCase.Logout - repo.RevokeRefreshSession: %w", err)
	}
	return nil
}

func (uc *authUseCase) Me(ctx context.Context, actor Actor) (User, error) {
	if actor.UserID == "" {
		return User{}, ErrUnauthorized
	}

	u, err := uc.repo.GetUserByID(ctx, actor.UserID)
	if err != nil {
		return User{}, fmt.Errorf("AuthUseCase.Me - repo.GetUserByID: %w", err)
	}
	u.IsAdmin = u.IsAdmin || uc.isAdminUsername(u.Username)
	return u, nil
}

func (uc *authUseCase) issueTokenPair(ctx context.Context, userID string) (string, string, error) {
	u, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("AuthUseCase.issueTokenPair - repo.GetUserByID: %w", err)
	}
	isAdmin := u.IsAdmin || uc.isAdminUsername(u.Username)

	accessToken, err := uc.jwtManager.GenerateTokenWithProfile(userID, u.Username, isAdmin)
	if err != nil {
		return "", "", fmt.Errorf("AuthUseCase.issueTokenPair - jwtManager.GenerateTokenWithProfile: %w", err)
	}

	now := time.Now().UTC()
	refreshToken := uuid.NewString()
	session := RefreshSession{
		ID:        uuid.NewString(),
		UserID:    userID,
		TokenID:   refreshToken,
		ExpiresAt: now.Add(uc.refreshTTL),
		CreatedAt: now,
	}
	if err := uc.repo.StoreRefreshSession(ctx, session); err != nil {
		return "", "", fmt.Errorf("AuthUseCase.issueTokenPair - repo.StoreRefreshSession: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (uc *authUseCase) isAdminUsername(username string) bool {
	_, ok := uc.adminUsers[strings.ToLower(strings.TrimSpace(username))]
	return ok
}
