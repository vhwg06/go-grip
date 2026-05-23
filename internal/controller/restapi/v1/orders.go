package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type gripRequestRefundBody struct {
	Reason string `json:"reason"`
}

func (r *V1) gripListOrders(ctx *fiber.Ctx) error {
	if r.orders == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "orders_usecase_not_configured"})
	}

	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		status, body := mapDomainError(entity.ErrUnauthorized)
		return ctx.Status(status).JSON(body)
	}

	page := gripPage(ctx)
	email := ctx.Query("email")
	items, total, err := r.orders.List(ctx.UserContext(), actor, email, page)
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	for i := range items {
		items[i] = gripDecorateOrderStatus(items[i])
	}

	normalized := page.Normalize()
	return ctx.JSON(gripListResponse{
		Data: items,
		Meta: entity.Page{Limit: normalized.Limit, Offset: normalized.Offset, Total: total},
	})
}

func (r *V1) gripGetOrder(ctx *fiber.Ctx) error {
	if r.orders == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "orders_usecase_not_configured"})
	}

	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		status, body := mapDomainError(entity.ErrUnauthorized)
		return ctx.Status(status).JSON(body)
	}

	order, err := r.orders.Get(ctx.UserContext(), actor, ctx.Params("id"))
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	return ctx.JSON(apiSuccessEnvelope(gripDecorateOrderStatus(order)))
}

func (r *V1) gripRequestRefund(ctx *fiber.Ctx) error {
	if r.orders == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "orders_usecase_not_configured"})
	}

	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		status, body := mapDomainError(entity.ErrUnauthorized)
		return ctx.Status(status).JSON(body)
	}

	var body gripRequestRefundBody
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if body.Reason == "" {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	refund, err := r.orders.RequestRefund(ctx.UserContext(), actor, ctx.Params("id"), body.Reason)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.Status(http.StatusCreated).JSON(apiSuccessEnvelope(refund))
}
