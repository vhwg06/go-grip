package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type AdminBannerSlide struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Image       string `json:"image"`
	MobileImage string `json:"mobileImage"`
	CtaText     string `json:"ctaText"`
	CtaLink     string `json:"ctaLink"`
	SortOrder   int    `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
}

func (r *V1) listAdminBanners(ctx *fiber.Ctx) error {
	log.Println("[banners-debug] listAdminBanners called")
	blocks, err := r.homepage.ListBlocks(ctx.UserContext(), false)
	if err != nil {
		log.Printf("[banners-debug] listAdminBanners ListBlocks error: %v", err)
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	var bannerBlock *entity.HomepageBlock
	for i := range blocks {
		if blocks[i].BlockType == "banner" {
			bannerBlock = &blocks[i]
			break
		}
	}

	if bannerBlock == nil {
		log.Println("[banners-debug] listAdminBanners: no banner block found")
		return ctx.JSON([]any{})
	}

	slides := extractBannerSlides(*bannerBlock)

	log.Printf("[banners-debug] listAdminBanners returning %d slides: %+v", len(slides), slides)
	return ctx.JSON(slides)
}

func (r *V1) saveAdminBanner(ctx *fiber.Ctx) error {
	log.Println("[banners-debug] saveAdminBanner called")
	blocks, err := r.homepage.ListBlocks(ctx.UserContext(), false)
	if err != nil {
		log.Printf("[banners-debug] saveAdminBanner ListBlocks error: %v", err)
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	var bannerBlock entity.HomepageBlock
	found := false
	for i := range blocks {
		if blocks[i].BlockType == "banner" {
			bannerBlock = blocks[i]
			found = true
			break
		}
	}

	if !found {
		bannerBlock = entity.HomepageBlock{
			BlockType: "banner",
			IsActive:  true,
			Config:    map[string]any{"slides": []any{}},
		}
	}

	slides := extractBannerSlides(bannerBlock)
	payload := parseBannerPayload(ctx)

	log.Printf("[banners-debug] saveAdminBanner inputs: id=%d, title=%s, sortOrder=%d, isActive=%t", payload.ID, payload.Title, payload.SortOrder, payload.IsActive)

	if payload.ID > 0 {
		// Update existing
		updated := false
		for i := range slides {
			if slides[i].ID == payload.ID {
				slides[i] = payload
				updated = true
				break
			}
		}
		if !updated {
			slides = append(slides, payload)
		}
	} else {
		// Create new, find max ID
		maxID := 0
		for _, s := range slides {
			if s.ID > maxID {
				maxID = s.ID
			}
		}
		newID := maxID + 1
		payload.ID = newID
		slides = append(slides, payload)
		log.Printf("[banners-debug] saveAdminBanner created new slide ID=%d", newID)
	}

	sortBannerSlides(slides)
	bannerBlock.Config["slides"] = slides

	if !found {
		savedBlock, storeErr := r.homepage.StoreBlock(ctx.UserContext(), bannerBlock)
		if storeErr == nil {
			bannerBlock = savedBlock
		}
		err = storeErr
	} else {
		_, err = r.homepage.UpdateBlock(ctx.UserContext(), bannerBlock)
	}

	if err != nil {
		log.Printf("[banners-debug] saveAdminBanner save error: %v", err)
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	log.Printf("[banners-debug] saveAdminBanner success, block ID=%s", bannerBlock.ID)
	return ctx.JSON(fiber.Map{"success": true})
}

func (r *V1) deleteAdminBanner(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	log.Printf("[banners-debug] deleteAdminBanner called for idStr=%s (parsed id=%d, err=%v)", idStr, id, err)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	blocks, err := r.homepage.ListBlocks(ctx.UserContext(), false)
	if err != nil {
		log.Printf("[banners-debug] deleteAdminBanner ListBlocks error: %v", err)
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	var bannerBlock entity.HomepageBlock
	found := false
	for i := range blocks {
		if blocks[i].BlockType == "banner" {
			bannerBlock = blocks[i]
			found = true
			break
		}
	}

	if !found {
		log.Println("[banners-debug] deleteAdminBanner: no banner block found to delete from")
		return ctx.JSON(fiber.Map{"success": true})
	}

	slides := []AdminBannerSlide{}
	slides = extractBannerSlides(bannerBlock)

	log.Printf("[banners-debug] deleteAdminBanner: current slides: %+v", slides)

	nextSlides := []AdminBannerSlide{}
	for _, s := range slides {
		if s.ID != id {
			nextSlides = append(nextSlides, s)
		}
	}

	log.Printf("[banners-debug] deleteAdminBanner: slides after deletion: %+v", nextSlides)
	sortBannerSlides(nextSlides)
	bannerBlock.Config["slides"] = nextSlides

	_, err = r.homepage.UpdateBlock(ctx.UserContext(), bannerBlock)
	if err != nil {
		log.Printf("[banners-debug] deleteAdminBanner UpdateBlock error: %v", err)
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	log.Println("[banners-debug] deleteAdminBanner completed successfully")
	return ctx.JSON(fiber.Map{"success": true})
}

func extractBannerSlides(block entity.HomepageBlock) []AdminBannerSlide {
	slides := []AdminBannerSlide{}
	if slidesRaw, ok := block.Config["slides"]; ok {
		jsBytes, _ := json.Marshal(slidesRaw)
		_ = json.Unmarshal(jsBytes, &slides)
	}
	sortBannerSlides(slides)
	return slides
}

func sortBannerSlides(slides []AdminBannerSlide) {
	slices.SortFunc(slides, func(a, b AdminBannerSlide) int {
		if a.SortOrder != b.SortOrder {
			return a.SortOrder - b.SortOrder
		}
		return a.ID - b.ID
	})
}

func parseBannerPayload(ctx *fiber.Ctx) AdminBannerSlide {
	var payload AdminBannerSlide
	if len(ctx.Body()) > 0 {
		_ = ctx.BodyParser(&payload)
	}
	if payload.Title != "" || payload.Image != "" || payload.SortOrder != 0 || payload.ID != 0 || payload.MobileImage != "" || payload.CtaText != "" || payload.CtaLink != "" {
		return payload
	}

	id, _ := strconv.Atoi(ctx.FormValue("id"))
	sortOrder, _ := strconv.Atoi(ctx.FormValue("sortOrder"))
	isActive := true
	if strings.TrimSpace(ctx.FormValue("isActive")) == "false" {
		isActive = false
	}
	return AdminBannerSlide{
		ID:          id,
		Title:       ctx.FormValue("title"),
		Subtitle:    ctx.FormValue("subtitle"),
		Image:       ctx.FormValue("image"),
		MobileImage: ctx.FormValue("mobileImage"),
		CtaText:     ctx.FormValue("ctaText"),
		CtaLink:     ctx.FormValue("ctaLink"),
		SortOrder:   sortOrder,
		IsActive:    isActive,
	}
}
