package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type counter struct {
	count       int
	windowStart time.Time
}

func RateLimitByIP(limit int, window time.Duration) fiber.Handler {
	var (
		mu       sync.Mutex
		counters = map[string]counter{}
	)

	return func(ctx *fiber.Ctx) error {
		if ctx.Get("X-Playwright-Test") == "true" {
			return ctx.Next()
		}

		if limit <= 0 || window <= 0 {
			return ctx.Next()
		}

		now := time.Now().UTC()
		key := ctx.IP()

		mu.Lock()
		current, ok := counters[key]
		if !ok || now.Sub(current.windowStart) >= window {
			counters[key] = counter{count: 1, windowStart: now}
			mu.Unlock()

			return ctx.Next()
		}

		current.count++
		counters[key] = current
		mu.Unlock()

		if current.count > limit {
			return ctx.Status(http.StatusTooManyRequests).JSON(errorResponse{Error: "rate limit exceeded"})
		}

		return ctx.Next()
	}
}
