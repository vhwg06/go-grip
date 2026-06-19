package app

import (
	"strings"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestRun_UsesRESTServerOnly_TDD(t *testing.T) {
	t.Setenv("APP_NAME", "go-grip")
	t.Setenv("APP_VERSION", "test")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("LOG_LEVEL", "error")
	t.Setenv("PG_POOL_MAX", "1")
	t.Setenv("PG_URL", "postgres://example.invalid/test")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("ADMIN_USERS", "admin")
	t.Setenv("ECOMMERCE_SCHEDULER_INTERVAL", "2m")

	cfg, err := config.NewConfig()
	require.NoError(t, err)

	s := initServers(cfg, useCases{}, jwt.New(cfg.JWT.Secret, time.Hour), logger.New("error"))
	t.Cleanup(func() {
		s.http.App.Shutdown()
		if s.maintenanceTicker != nil {
			s.maintenanceTicker.Stop()
		}
		if s.maintenanceDone != nil {
			close(s.maintenanceDone)
		}
	})

	require.NotNil(t, s.http)
	require.NotNil(t, s.maintenanceTicker)
	require.NotNil(t, s.maintenanceDone)

	var routes []string
	for _, group := range s.http.App.Stack() {
		for _, route := range group {
			if route == nil {
				continue
			}
			routes = append(routes, route.Path)
		}
	}

	require.Contains(t, routes, "/v1/catalog/products")
	require.Contains(t, routes, "/v1/admin/settings")
	require.Contains(t, routes, "/v1/checkout/orders")

	for _, path := range routes {
		lower := strings.ToLower(path)
		require.NotContains(t, lower, "grpc")
		require.NotContains(t, lower, "nats")
		require.NotContains(t, lower, "amqp")
		require.NotContains(t, lower, "rabbit")
	}
}
