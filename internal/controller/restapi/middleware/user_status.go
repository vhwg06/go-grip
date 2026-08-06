package middleware

import (
	"net/http"

	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/gofiber/fiber/v2"
)

func RejectBlockedMutations() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		switch ctx.Method() {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		default:
			return ctx.Next()
		}

		actor, ok := ctx.Locals("actor").(usermodule.Actor)
		if !ok {
			return ctx.Next()
		}

		if actor.IsBlocked {
			return ctx.Status(http.StatusForbidden).JSON(errorResponse{Error: "blocked user"})
		}

		return ctx.Next()
	}
}
