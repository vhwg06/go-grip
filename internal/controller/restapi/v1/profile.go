package v1

import (
	"context"
	"net/http"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type gripProfileUpdateRequest struct {
	Email                       *string `json:"email"`
	DisplayName                 *string `json:"displayName"`
	DisplayNameLegacy           *string `json:"display_name"`
	DesktopNotificationsEnabled *bool   `json:"desktopNotificationsEnabled"`
	DesktopNotificationsLegacy  *bool   `json:"desktop_notifications_enabled"`
}

// profileEnvelope keeps the canonical data envelope while exposing the same
// identity fields at the boundary used by the legacy profile clients.
func profileEnvelope(user entity.User) fiber.Map {
	return fiber.Map{
		"data":                          user,
		"id":                            user.ID,
		"username":                      user.Username,
		"display_name":                  user.DisplayName,
		"email":                         user.Email,
		"role_id":                       user.RoleID,
		"role":                          user.Role,
		"is_admin":                      user.IsAdmin,
		"desktop_notifications_enabled": user.DesktopNotificationsEnabled,
	}
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

	return ctx.JSON(profileEnvelope(user))
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
	current, err := r.profileUC.Get(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	email := current.Email
	if body.Email != nil {
		email = strings.TrimSpace(*body.Email)
	}
	displayName := current.DisplayName
	if body.DisplayName != nil {
		displayName = *body.DisplayName
	} else if body.DisplayNameLegacy != nil {
		displayName = *body.DisplayNameLegacy
	}
	desktopNotificationsEnabled := current.DesktopNotificationsEnabled
	if body.DesktopNotificationsEnabled != nil {
		desktopNotificationsEnabled = *body.DesktopNotificationsEnabled
	} else if body.DesktopNotificationsLegacy != nil {
		desktopNotificationsEnabled = *body.DesktopNotificationsLegacy
	}
	if email == "" {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	user, err := r.profileUC.Update(ctx.UserContext(), r.gripActor(ctx), email, displayName, desktopNotificationsEnabled)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(profileEnvelope(user))
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
