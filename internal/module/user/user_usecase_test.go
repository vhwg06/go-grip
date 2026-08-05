package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/stretchr/testify/require"
)

func TestUserUseCaseRegisterAndLogin(t *testing.T) {
	t.Parallel()
	jwtMgr := jwt.New("secret-key", time.Hour)
	repo := persistent.NewUserRepo(nil)
	uc := user.NewUserUseCase(repo, jwtMgr)

	registered, err := uc.Register(context.Background(), "johndoe", "john@example.com", "password123")
	require.NoError(t, err)
	require.Equal(t, "johndoe", registered.Username)

	token, err := uc.Login(context.Background(), "john@example.com", "password123")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	_, err = uc.Login(context.Background(), "john@example.com", "wrongpassword")
	require.ErrorIs(t, err, user.ErrInvalidCredentials)
}
