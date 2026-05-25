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

func (r *V1) getPresignedURL(ctx *fiber.Ctx) error {
	fileName := ctx.Query("fileName")
	contentType := ctx.Query("contentType")
	if fileName == "" || contentType == "" {
		return errorResponse(ctx, http.StatusBadRequest, "fileName and contentType query parameters are required")
	}

	uploadURL, publicURL, fileID, err := r.media.GeneratePresignedURL(ctx.UserContext(), fileName, contentType)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"upload_url": uploadURL,
		"public_url": publicURL,
		"id":         fileID,
	})
}

func (r *V1) simulateUpload(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusOK)
}
