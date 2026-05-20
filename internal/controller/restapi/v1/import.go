package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type importBody struct {
	Items []entity.ImportItem `json:"items"`
}

func (r *V1) importInitialContent(ctx *fiber.Ctx) error {
	var body importBody
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	result, err := r.importer.Import(ctx.UserContext(), body.Items)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(result)
}
