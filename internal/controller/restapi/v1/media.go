package v1

import (
	"net/http"
	"os"
	"path/filepath"

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
	// Accept both limit/offset (legacy) and pageSize/page (frontend convention)
	pageSize := ctx.QueryInt("pageSize", 0)
	if pageSize <= 0 {
		pageSize = ctx.QueryInt("limit", 20)
	}
	page := ctx.QueryInt("page", 1)
	offset := ctx.QueryInt("offset", (page-1)*pageSize)
	q := ctx.Query("q")

	pagination := entity.Pagination{Limit: pageSize, Offset: offset}
	items, total, err := r.media.List(ctx.UserContext(), pagination, q)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	norm := pagination.Normalize()
	return ctx.JSON(listResponse{Data: items, Meta: entity.Page{Limit: norm.Limit, Offset: norm.Offset, Total: total}})
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
	filename := ctx.Params("filename")
	if filename == "" {
		return ctx.SendStatus(http.StatusOK)
	}
	uploadsDir := "/tmp/uploads"
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		return ctx.SendStatus(http.StatusOK)
	}
	destPath := uploadsDir + "/" + filepath.Base(filename)
	body := ctx.Body()
	if len(body) > 0 {
		_ = os.WriteFile(destPath, body, 0o644)
	}
	return ctx.SendStatus(http.StatusOK)
}
