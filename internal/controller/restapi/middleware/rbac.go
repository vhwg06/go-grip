package middleware

import (
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

// RequireRole allows requests whose role is in allowed.
func RequireRole(allowed ...entity.RoleName) fiber.Handler {
	allowedSet := make(map[entity.RoleName]struct{}, len(allowed))
	for _, role := range allowed {
		allowedSet[role] = struct{}{}
	}

	return func(ctx *fiber.Ctx) error {
		role, _ := ctx.Locals("role").(entity.RoleName)
		if _, ok := allowedSet[role]; ok {
			return ctx.Next()
		}

		return ctx.SendStatus(fiber.StatusForbidden)
	}
}
