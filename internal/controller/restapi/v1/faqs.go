package v1

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type adminFAQEntry struct {
	ID        int    `json:"id"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	SortOrder int    `json:"sortOrder"`
	IsActive  bool   `json:"isActive"`
}

func (r *V1) listAdminFAQs(ctx *fiber.Ctx) error {
	block, err := r.loadFAQBlock(ctx, false)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	if block == nil {
		return ctx.JSON([]adminFAQEntry{})
	}
	return ctx.JSON(extractFAQEntries(*block, false))
}

func (r *V1) saveAdminFAQ(ctx *fiber.Ctx) error {
	blocks, err := r.homepage.ListBlocks(ctx.UserContext(), false)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	var faqBlock entity.HomepageBlock
	found := false
	for i := range blocks {
		if blocks[i].BlockType == "faq" {
			faqBlock = blocks[i]
			found = true
			break
		}
	}
	if !found {
		faqBlock = entity.HomepageBlock{
			BlockType: "faq",
			IsActive:  true,
			Config:    map[string]any{"entries": []any{}},
		}
	}

	entries := extractFAQEntries(faqBlock, false)
	payload, err := parseFAQPayload(ctx)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	if strings.TrimSpace(payload.Question) == "" || strings.TrimSpace(payload.Answer) == "" {
		return errorResponse(ctx, http.StatusBadRequest, "question and answer are required")
	}

	if payload.ID > 0 {
		updated := false
		for i := range entries {
			if entries[i].ID == payload.ID {
				entries[i] = payload
				updated = true
				break
			}
		}
		if !updated {
			entries = append(entries, payload)
		}
	} else {
		nextID := 1
		for _, item := range entries {
			if item.ID >= nextID {
				nextID = item.ID + 1
			}
		}
		payload.ID = nextID
		entries = append(entries, payload)
	}

	sortFAQEntries(entries)
	if faqBlock.Config == nil {
		faqBlock.Config = map[string]any{}
	}
	faqBlock.Config["entries"] = entries

	if !found {
		_, err = r.homepage.StoreBlock(ctx.UserContext(), faqBlock)
	} else {
		_, err = r.homepage.UpdateBlock(ctx.UserContext(), faqBlock)
	}
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success":   true,
		"id":        payload.ID,
		"question":  payload.Question,
		"answer":    payload.Answer,
		"sortOrder": payload.SortOrder,
		"isActive":  payload.IsActive,
	})
}

func (r *V1) deleteAdminFAQ(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	block, err := r.loadFAQBlock(ctx, false)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	if block == nil {
		return ctx.JSON(fiber.Map{"success": true})
	}

	entries := extractFAQEntries(*block, false)
	nextEntries := entries[:0]
	for _, item := range entries {
		if item.ID != id {
			nextEntries = append(nextEntries, item)
		}
	}
	block.Config["entries"] = nextEntries

	if _, err := r.homepage.UpdateBlock(ctx.UserContext(), *block); err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{"success": true})
}

func (r *V1) listActiveFAQs(ctx *fiber.Ctx) error {
	block, err := r.loadFAQBlock(ctx, true)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	items := []adminFAQEntry{}
	if block != nil {
		items = extractFAQEntries(*block, true)
	}

	type faqItem struct {
		ID        string `json:"id"`
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		SortOrder int    `json:"sortOrder"`
		IsActive  bool   `json:"isActive"`
	}
	out := make([]faqItem, 0, len(items))
	for _, item := range items {
		out = append(out, faqItem{
			ID:        strconv.Itoa(item.ID),
			Question:  item.Question,
			Answer:    item.Answer,
			SortOrder: item.SortOrder,
			IsActive:  item.IsActive,
		})
	}
	return ctx.JSON(fiber.Map{"items": out})
}

func (r *V1) loadFAQBlock(ctx *fiber.Ctx, activeOnly bool) (*entity.HomepageBlock, error) {
	blocks, err := r.homepage.ListBlocks(ctx.UserContext(), activeOnly)
	if err != nil {
		return nil, err
	}
	for i := range blocks {
		if blocks[i].BlockType == "faq" {
			return &blocks[i], nil
		}
	}
	return nil, nil
}

func extractFAQEntries(block entity.HomepageBlock, activeOnly bool) []adminFAQEntry {
	entries := []adminFAQEntry{}
	if raw, ok := block.Config["entries"]; ok {
		jsBytes, _ := json.Marshal(raw)
		_ = json.Unmarshal(jsBytes, &entries)
	}
	if activeOnly {
		filtered := entries[:0]
		for _, item := range entries {
			if item.IsActive {
				filtered = append(filtered, item)
			}
		}
		entries = filtered
	}
	sortFAQEntries(entries)
	return entries
}

func sortFAQEntries(entries []adminFAQEntry) {
	slices.SortFunc(entries, func(a, b adminFAQEntry) int {
		if a.SortOrder != b.SortOrder {
			return a.SortOrder - b.SortOrder
		}
		return a.ID - b.ID
	})
}

func parseFAQPayload(ctx *fiber.Ctx) (adminFAQEntry, error) {
	var payload adminFAQEntry
	var raw map[string]any
	if len(ctx.Body()) > 0 {
		if err := ctx.BodyParser(&payload); err == nil {
			_ = json.Unmarshal(ctx.Body(), &raw)
			if payload.ID == 0 {
				payload.ID = jsonInt(raw, "id")
			}
			if payload.SortOrder == 0 {
				payload.SortOrder = jsonInt(raw, "sortOrder", "sort_order")
			}
			if payload.Question == "" {
				payload.Question = jsonString(raw, "question")
			}
			if payload.Answer == "" {
				payload.Answer = jsonString(raw, "answer")
			}
			if value, ok := jsonBool(raw, "isActive", "is_active", "active"); ok {
				payload.IsActive = value
			}
			if len(raw) > 0 {
				return payload, nil
			}
		}
	}

	id, _ := strconv.Atoi(ctx.FormValue("id"))
	sortOrder, _ := strconv.Atoi(ctx.FormValue("sortOrder"))
	isActive := true
	if value := strings.TrimSpace(ctx.FormValue("isActive")); value != "" {
		isActive = value != "false"
	}

	payload = adminFAQEntry{
		ID:        id,
		Question:  strings.TrimSpace(ctx.FormValue("question")),
		Answer:    strings.TrimSpace(ctx.FormValue("answer")),
		SortOrder: sortOrder,
		IsActive:  isActive,
	}
	return payload, nil
}

func jsonInt(raw map[string]any, names ...string) int {
	for _, name := range names {
		value, ok := raw[name]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case float64:
			return int(typed)
		case string:
			parsed, _ := strconv.Atoi(strings.TrimSpace(typed))
			return parsed
		}
	}
	return 0
}

func jsonString(raw map[string]any, names ...string) string {
	for _, name := range names {
		if value, ok := raw[name].(string); ok {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func jsonBool(raw map[string]any, names ...string) (bool, bool) {
	for _, name := range names {
		value, ok := raw[name]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case bool:
			return typed, true
		case string:
			parsed, err := strconv.ParseBool(strings.TrimSpace(typed))
			return parsed, err == nil
		}
	}
	return false, false
}
