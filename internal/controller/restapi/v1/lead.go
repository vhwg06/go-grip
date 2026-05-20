package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func (r *V1) submitLead(ctx *fiber.Ctx) error {
	var lead entity.LeadSubmission
	if err := ctx.BodyParser(&lead); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	lead, err := r.lead.Submit(ctx.UserContext(), lead)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusCreated).JSON(lead)
}

func (r *V1) getLead(ctx *fiber.Ctx) error {
	lead, err := r.lead.Get(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, "lead not found")
	}
	return ctx.JSON(lead)
}
