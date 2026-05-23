package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/webapi"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/google/uuid"
)

type UseCase struct {
	repo       repo.AuthRepository
	linuxDO    webapi.OAuthClient
	gitHub     webapi.OAuthClient
	jwtManager *jwt.Manager
	refreshTTL time.Duration
	adminUsers map[string]struct{}
}

func New(authRepo repo.AuthRepository, linuxDO webapi.OAuthClient, gitHub webapi.OAuthClient, jwtManager *jwt.Manager, refreshTTL time.Duration, adminUsersCSV string) *UseCase {
	adminUsers := make(map[string]struct{})
	for _, user := range strings.Split(adminUsersCSV, ",") {
		trimmed := strings.ToLower(strings.TrimSpace(user))
		if trimmed == "" {
			continue
		}
		adminUsers[trimmed] = struct{}{}
	}

	return &UseCase{
		repo:       authRepo,
		linuxDO:    linuxDO,
		gitHub:     gitHub,
		jwtManager: jwtManager,
		refreshTTL: refreshTTL,
		adminUsers: adminUsers,
	}
}

var _ usecase.Auth = (*UseCase)(nil)

func (uc *UseCase) BeginLinuxDO(ctx context.Context) (string, error) {
	if uc.linuxDO == nil {
		return "", entity.ErrInvalidInput
	}
	u, err := uc.linuxDO.BeginAuth(ctx, uuid.NewString())
	if err != nil {
		return "", fmt.Errorf("AuthUseCase.BeginLinuxDO - linuxDO.BeginAuth: %w", err)
	}
	return u, nil
}

func (uc *UseCase) BeginGitHub(ctx context.Context) (string, error) {
	if uc.gitHub == nil {
		return "", entity.ErrInvalidInput
	}
	u, err := uc.gitHub.BeginAuth(ctx, uuid.NewString())
	if err != nil {
		return "", fmt.Errorf("AuthUseCase.BeginGitHub - gitHub.BeginAuth: %w", err)
	}
	return u, nil
}

func (uc *UseCase) CompleteLinuxDO(ctx context.Context, code string) (entity.User, string, string, error) {
	return uc.completeOAuth(ctx, uc.linuxDO, code)
}

func (uc *UseCase) CompleteGitHub(ctx context.Context, code string) (entity.User, string, string, error) {
	return uc.completeOAuth(ctx, uc.gitHub, code)
}

func (uc *UseCase) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	if refreshToken == "" {
		return "", "", entity.ErrInvalidInput
	}

	session, err := uc.repo.GetRefreshSession(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("AuthUseCase.Refresh - repo.GetRefreshSession: %w", err)
	}
	if !session.RevokedAt.IsZero() || session.ExpiresAt.Before(time.Now().UTC()) {
		return "", "", entity.ErrUnauthorized
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

func (uc *UseCase) Logout(ctx context.Context, actor entity.Actor, refreshToken string) error {
	if actor.UserID == "" {
		return entity.ErrUnauthorized
	}
	if refreshToken == "" {
		return entity.ErrInvalidInput
	}
	if err := uc.repo.RevokeRefreshSession(ctx, refreshToken); err != nil {
		return fmt.Errorf("AuthUseCase.Logout - repo.RevokeRefreshSession: %w", err)
	}
	return nil
}

func (uc *UseCase) Me(ctx context.Context, actor entity.Actor) (entity.User, error) {
	if actor.UserID == "" {
		return entity.User{}, entity.ErrUnauthorized
	}

	user, err := uc.repo.GetUserByID(ctx, actor.UserID)
	if err != nil {
		return entity.User{}, fmt.Errorf("AuthUseCase.Me - repo.GetUserByID: %w", err)
	}
	user.IsAdmin = user.IsAdmin || uc.isAdminUsername(user.Username)
	return user, nil
}

func (uc *UseCase) completeOAuth(ctx context.Context, client webapi.OAuthClient, code string) (entity.User, string, string, error) {
	if client == nil {
		return entity.User{}, "", "", entity.ErrInvalidInput
	}
	if strings.TrimSpace(code) == "" {
		return entity.User{}, "", "", entity.ErrInvalidInput
	}

	identity, err := client.ExchangeCode(ctx, code)
	if err != nil {
		return entity.User{}, "", "", fmt.Errorf("AuthUseCase.completeOAuth - client.ExchangeCode: %w", err)
	}

	user, err := uc.repo.UpsertUser(ctx, entity.User{
		Provider:   identity.Provider,
		ProviderID: identity.ProviderID,
		Username:   identity.Username,
		Email:      identity.Email,
		Status:     entity.UserStatusActive,
	})
	if err != nil {
		return entity.User{}, "", "", fmt.Errorf("AuthUseCase.completeOAuth - repo.UpsertUser: %w", err)
	}

	user.IsAdmin = user.IsAdmin || uc.isAdminUsername(user.Username)
	accessToken, refreshToken, err := uc.issueTokenPair(ctx, user.ID)
	if err != nil {
		return entity.User{}, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (uc *UseCase) issueTokenPair(ctx context.Context, userID string) (string, string, error) {
	accessToken, err := uc.jwtManager.GenerateToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("AuthUseCase.issueTokenPair - jwtManager.GenerateToken: %w", err)
	}

	now := time.Now().UTC()
	refreshToken := uuid.NewString()
	session := entity.RefreshSession{
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

func (uc *UseCase) isAdminUsername(username string) bool {
	_, ok := uc.adminUsers[strings.ToLower(strings.TrimSpace(username))]
	return ok
}
