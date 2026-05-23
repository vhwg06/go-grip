package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func RequireAdminUsernames(adminUsersCSV string) fiber.Handler {
	rawUsers := strings.Split(adminUsersCSV, ",")
	adminUsers := make([]string, 0, len(rawUsers))
	for _, username := range rawUsers {
		trimmed := strings.TrimSpace(strings.ToLower(username))
		if trimmed == "" {
			continue
		}
		adminUsers = append(adminUsers, trimmed)
	}

	return func(ctx *fiber.Ctx) error {
		actor, ok := ctx.Locals("actor").(entity.Actor)
		if !ok {
			return ctx.Status(http.StatusUnauthorized).JSON(errorResponse{Error: "missing actor context"})
		}

		if actor.IsAdmin || slices.Contains(adminUsers, strings.ToLower(actor.Username)) {
			return ctx.Next()
		}

		return ctx.Status(http.StatusForbidden).JSON(errorResponse{Error: "admin access required"})
	}
}
