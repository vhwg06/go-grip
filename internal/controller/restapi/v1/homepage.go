package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func (r *V1) createHomepageBlock(ctx *fiber.Ctx) error {
	var block entity.HomepageBlock
	if err := ctx.BodyParser(&block); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	block, err := r.homepage.StoreBlock(ctx.UserContext(), block)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusCreated).JSON(block)
}

func (r *V1) listHomepageBlocks(ctx *fiber.Ctx) error {
	items, err := r.homepage.ListBlocks(ctx.UserContext(), false)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(listResponse{Data: items})
}

func (r *V1) listPublicHomepage(ctx *fiber.Ctx) error {
	items, err := r.homepage.ListBlocks(ctx.UserContext(), true)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	for i := range items {
		if items[i].BlockType == "banner" {
			slides := extractBannerSlides(items[i])
			activeSlides := slides[:0]
			for _, slide := range slides {
				if slide.IsActive {
					activeSlides = append(activeSlides, slide)
				}
			}
			items[i].Config["slides"] = activeSlides
		} else if items[i].BlockType == "faq" {
			entries := extractFAQEntries(items[i], true)
			items[i].Config["entries"] = entries
		}
	}
	return ctx.JSON(listResponse{Data: items})
}

func (r *V1) updateHomepageBlock(ctx *fiber.Ctx) error {
	var block entity.HomepageBlock
	if err := ctx.BodyParser(&block); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	block.ID = ctx.Params("id")
	block, err := r.homepage.UpdateBlock(ctx.UserContext(), block)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(block)
}

func (r *V1) deleteHomepageBlock(ctx *fiber.Ctx) error {
	if err := r.homepage.DeleteBlock(ctx.UserContext(), ctx.Params("id")); err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) listSupportChannels(ctx *fiber.Ctx) error {
	items, err := r.homepage.ListSupport(ctx.UserContext(), false)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(listResponse{Data: items})
}

func (r *V1) listPublicSupport(ctx *fiber.Ctx) error {
	items, err := r.homepage.ListSupport(ctx.UserContext(), true)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(listResponse{Data: items})
}

func (r *V1) updateSupportChannel(ctx *fiber.Ctx) error {
	var channel entity.SupportChannel
	if err := ctx.BodyParser(&channel); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	channel.ID = ctx.Params("id")
	channel, err := r.homepage.UpdateSupport(ctx.UserContext(), channel)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(channel)
}
