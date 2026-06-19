package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig_BuildsPostgresURLFromDiscreteEnv(t *testing.T) {
	t.Setenv("APP_NAME", "go-grip")
	t.Setenv("APP_VERSION", "test")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("PG_POOL_MAX", "2")
	t.Setenv("POSTGRES_HOST", "db")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_DB", "grip")
	t.Setenv("POSTGRES_USER", "grip")
	t.Setenv("POSTGRES_PASSWORD", "p@ss:/?#word")
	t.Setenv("POSTGRES_SSL_MODE", "disable")
	t.Setenv("JWT_SECRET", "secret")

	cfg, err := NewConfig()
	require.NoError(t, err)
	require.Equal(t, "postgres://grip:p%40ss%3A%2F%3F%23word@db:5432/grip?sslmode=disable", cfg.PG.URL)
}

