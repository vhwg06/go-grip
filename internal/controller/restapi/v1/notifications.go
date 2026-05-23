package v1

import (
	"net/http"
	"strconv"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func (r *V1) gripNotificationsList(ctx *fiber.Ctx) error {
	if r.notifyUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "notification_usecase_not_configured"})
	}

	page := gripPage(ctx)
	items, total, err := r.notifyUC.Inbox(ctx.UserContext(), r.gripActor(ctx), page)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	normalized := page.Normalize()
	return ctx.JSON(gripListResponse{Data: items, Meta: entity.Page{Limit: normalized.Limit, Offset: normalized.Offset, Total: total}})
}

func (r *V1) gripNotificationsUnread(ctx *fiber.Ctx) error {
	if r.notifyUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "notification_usecase_not_configured"})
	}

	count, err := r.notifyUC.UnreadCount(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.JSON(apiSuccessEnvelope(fiber.Map{"count": count}))
}

func (r *V1) gripNotificationsMarkRead(ctx *fiber.Ctx) error {
	if r.notifyUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "notification_usecase_not_configured"})
	}
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if err := r.notifyUC.MarkRead(ctx.UserContext(), r.gripActor(ctx), id); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) gripNotificationsReadAll(ctx *fiber.Ctx) error {
	if r.notifyUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "notification_usecase_not_configured"})
	}
	if err := r.notifyUC.MarkAllRead(ctx.UserContext(), r.gripActor(ctx)); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) gripNotificationsClear(ctx *fiber.Ctx) error {
	if r.notifyUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "notification_usecase_not_configured"})
	}
	if err := r.notifyUC.Clear(ctx.UserContext(), r.gripActor(ctx)); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.SendStatus(http.StatusNoContent)
}
