package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

	slides := []AdminBannerSlide{}
	if slidesRaw, ok := bannerBlock.Config["slides"]; ok {
		jsBytes, _ := json.Marshal(slidesRaw)
		_ = json.Unmarshal(jsBytes, &slides)
	}

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

	slides := []AdminBannerSlide{}
	if slidesRaw, ok := bannerBlock.Config["slides"]; ok {
		jsBytes, _ := json.Marshal(slidesRaw)
		_ = json.Unmarshal(jsBytes, &slides)
	}

	idStr := ctx.FormValue("id")
	title := ctx.FormValue("title")
	subtitle := ctx.FormValue("subtitle")
	image := ctx.FormValue("image")
	mobileImage := ctx.FormValue("mobileImage")
	ctaText := ctx.FormValue("ctaText")
	ctaLink := ctx.FormValue("ctaLink")
	sortOrderStr := ctx.FormValue("sortOrder")
	isActiveStr := ctx.FormValue("isActive")

	id, _ := strconv.Atoi(idStr)
	sortOrder, _ := strconv.Atoi(sortOrderStr)
	isActive := true
	if isActiveStr == "false" {
		isActive = false
	}

	log.Printf("[banners-debug] saveAdminBanner inputs: id=%d (%s), title=%s, sortOrder=%d, isActive=%t", id, idStr, title, sortOrder, isActive)

	if id > 0 {
		// Update existing
		updated := false
		for i := range slides {
			if slides[i].ID == id {
				slides[i].Title = title
				slides[i].Subtitle = subtitle
				slides[i].Image = image
				slides[i].MobileImage = mobileImage
				slides[i].CtaText = ctaText
				slides[i].CtaLink = ctaLink
				slides[i].SortOrder = sortOrder
				slides[i].IsActive = isActive
				updated = true
				break
			}
		}
		if !updated {
			slides = append(slides, AdminBannerSlide{
				ID:          id,
				Title:       title,
				Subtitle:    subtitle,
				Image:       image,
				MobileImage: mobileImage,
				CtaText:     ctaText,
				CtaLink:     ctaLink,
				SortOrder:   sortOrder,
				IsActive:    isActive,
			})
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
		slides = append(slides, AdminBannerSlide{
			ID:          newID,
			Title:       title,
			Subtitle:    subtitle,
			Image:       image,
			MobileImage: mobileImage,
			CtaText:     ctaText,
			CtaLink:     ctaLink,
			SortOrder:   sortOrder,
			IsActive:    isActive,
		})
		log.Printf("[banners-debug] saveAdminBanner created new slide ID=%d", newID)
	}

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
	if slidesRaw, ok := bannerBlock.Config["slides"]; ok {
		jsBytes, _ := json.Marshal(slidesRaw)
		_ = json.Unmarshal(jsBytes, &slides)
	}

	log.Printf("[banners-debug] deleteAdminBanner: current slides: %+v", slides)

	nextSlides := []AdminBannerSlide{}
	for _, s := range slides {
		if s.ID != id {
			nextSlides = append(nextSlides, s)
		}
	}

	log.Printf("[banners-debug] deleteAdminBanner: slides after deletion: %+v", nextSlides)
	bannerBlock.Config["slides"] = nextSlides

	_, err = r.homepage.UpdateBlock(ctx.UserContext(), bannerBlock)
	if err != nil {
		log.Printf("[banners-debug] deleteAdminBanner UpdateBlock error: %v", err)
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	log.Println("[banners-debug] deleteAdminBanner completed successfully")
	return ctx.JSON(fiber.Map{"success": true})
}
