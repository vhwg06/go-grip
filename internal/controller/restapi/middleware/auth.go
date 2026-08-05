package middleware

import (
	"context"
	"net/http"
	"strings"

	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

const _bearerParts = 2

type errorResponse struct {
	Error string `json:"error"`
}

// Auth returns a JWT authentication middleware for Fiber that extracts actor information from Bearer tokens when present.
func Auth(jwtManager *jwt.Manager) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		header := ctx.Get("Authorization")
		if header == "" {
			emptyActor := usermodule.Actor{}
			userCtx := context.WithValue(ctx.UserContext(), "actor", emptyActor)
			ctx.SetUserContext(userCtx)
			ctx.Locals("actor", emptyActor)
			return ctx.Next()
		}

		parts := strings.SplitN(header, " ", _bearerParts)
		if len(parts) != _bearerParts || parts[0] != "Bearer" {
			return ctx.Status(http.StatusUnauthorized).JSON(errorResponse{Error: "invalid authorization header format"})
		}

		userID, isAdmin, username, err := jwtManager.ParseTokenActor(parts[1])
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(errorResponse{Error: "invalid or expired token"})
		}

		actor := usermodule.Actor{
			UserID:   userID,
			Username: username,
			IsAdmin:  isAdmin,
		}

		userCtx := context.WithValue(ctx.UserContext(), "actor", actor)
		userCtx = context.WithValue(userCtx, "userID", userID)
		ctx.SetUserContext(userCtx)

		ctx.Locals("userID", userID)
		ctx.Locals("actor", actor)

		return ctx.Next()
	}
}
