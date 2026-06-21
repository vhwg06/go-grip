package v1

import (
	"context"
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type gripProfileUpdateRequest struct {
	Email                       string `json:"email"`
	DisplayName                 string `json:"displayName"`
	DesktopNotificationsEnabled bool   `json:"desktopNotificationsEnabled"`
}

func (r *V1) gripProfileGet(ctx *fiber.Ctx) error {
	if r.profileUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "profile_usecase_not_configured"})
	}

	user, err := r.profileUC.Get(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(user))
}

func (r *V1) gripProfileUpdate(ctx *fiber.Ctx) error {
	if r.profileUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "profile_usecase_not_configured"})
	}

	var body gripProfileUpdateRequest
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if body.Email == "" {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	user, err := r.profileUC.Update(ctx.UserContext(), r.gripActor(ctx), body.Email, body.DisplayName, body.DesktopNotificationsEnabled)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(user))
}

func (r *V1) gripProfileCheckin(ctx *fiber.Ctx) error {
	if r.profileUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "profile_usecase_not_configured"})
	}

	checkin, err := r.profileUC.Checkin(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(checkin))
}

func (r *V1) gripProfileGetSecurity(ctx *fiber.Ctx) error {
	ext, ok := r.profileUC.(interface {
		GetSecurityPosture(ctx context.Context, actor entity.Actor) (any, error)
	})
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "profile_security_not_available"})
	}

	data, err := ext.GetSecurityPosture(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(data))
}

func (r *V1) gripProfileGetSessions(ctx *fiber.Ctx) error {
	ext, ok := r.profileUC.(interface {
		GetRecentSessions(ctx context.Context, actor entity.Actor) (any, error)
	})
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "profile_sessions_not_available"})
	}

	data, err := ext.GetRecentSessions(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(data))
}

