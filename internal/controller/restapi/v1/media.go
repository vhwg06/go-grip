package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func (r *V1) createMedia(ctx *fiber.Ctx) error {
	var media entity.MediaAsset
	if err := ctx.BodyParser(&media); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	media, err := r.media.Store(ctx.UserContext(), media)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusCreated).JSON(media)
}

func (r *V1) listMedia(ctx *fiber.Ctx) error {
	items, total, err := r.media.List(ctx.UserContext(), queryPage(ctx))
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(listResponse{Data: items, Meta: entity.Page{Limit: queryPage(ctx).Normalize().Limit, Offset: queryPage(ctx).Normalize().Offset, Total: total}})
}

func (r *V1) deleteMedia(ctx *fiber.Ctx) error {
	if err := r.media.Delete(ctx.UserContext(), ctx.Params("id")); err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.SendStatus(http.StatusNoContent)
}
