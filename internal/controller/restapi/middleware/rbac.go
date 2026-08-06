package middleware

import (
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/gofiber/fiber/v2"
)

// RequireRole allows requests whose role is in allowed.
func RequireRole(allowed ...usermodule.RoleName) fiber.Handler {
	allowedSet := make(map[usermodule.RoleName]struct{}, len(allowed))
	for _, role := range allowed {
		allowedSet[role] = struct{}{}
	}

	return func(ctx *fiber.Ctx) error {
		role, _ := ctx.Locals("role").(usermodule.RoleName)
		if _, ok := allowedSet[role]; ok {
			return ctx.Next()
		}

		return ctx.SendStatus(fiber.StatusForbidden)
	}
}
