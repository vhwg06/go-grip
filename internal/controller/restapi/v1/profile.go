package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type gripProfileUpdateRequest struct {
	Email                       string `json:"email"`
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

	user, err := r.profileUC.Update(ctx.UserContext(), r.gripActor(ctx), body.Email, body.DesktopNotificationsEnabled)
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
