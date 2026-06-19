package v1

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

var storeSettingsPhonePattern = regexp.MustCompile(`^[+0-9][0-9 .()-]{5,}$`)

type storeSettingsAdminResponse struct {
	Config       storeSettingsConfig `json:"config"`
	Stats        storeSettingsStats  `json:"stats"`
	VisitorCount int                 `json:"visitorCount"`
}

type storeSettingsStats struct {
	Today storeSettingsStatBucket `json:"today"`
	Week  storeSettingsStatBucket `json:"week"`
	Month storeSettingsStatBucket `json:"month"`
	Total storeSettingsStatBucket `json:"total"`
}

type storeSettingsStatBucket struct {
	Count   int   `json:"count"`
	Revenue int64 `json:"revenue"`
}

type storeSettingsConfig struct {
	Brand        storeSettingsBrand            `json:"brand"`
	Contact      storeSettingsContact          `json:"contact"`
	Homepage     storeSettingsHomepage         `json:"homepage"`
	Footer       storeSettingsFooter           `json:"footer"`
	FloatSupport []storeSettingsFloatingAction `json:"floatingSupport"`
	Visibility   storeSettingsVisibility       `json:"visibility"`
	Registry     storeSettingsRegistry         `json:"registry"`
}

type storeSettingsBrand struct {
	ShopName        string `json:"shopName"`
	ShopDescription string `json:"shopDescription"`
	ShopLogo        string `json:"shopLogo"`
	ThemeColor      string `json:"themeColor"`
}

type storeSettingsContact struct {
	StickyBarAddress string `json:"stickyBarAddress"`
	StickyBarHotline string `json:"stickyBarHotline"`
	ContactEmail     string `json:"contactEmail"`
}

type storeSettingsHomepage struct {
	Blocks    []storeSettingsHomepageBlock `json:"blocks"`
	NewsCount int                          `json:"newsCount"`
}

type storeSettingsHomepageBlock struct {
	Key     string `json:"key"`
	Enabled bool   `json:"enabled"`
	Order   int    `json:"order"`
}

type storeSettingsFooter struct {
	Columns     []storeSettingsFooterColumn `json:"columns"`
	Copyright   string                      `json:"copyright"`
	SocialLinks map[string]string           `json:"socialLinks"`
}

type storeSettingsFooterColumn struct {
	ID    string                    `json:"id"`
	Title string                    `json:"title"`
	Links []storeSettingsFooterLink `json:"links"`
}

type storeSettingsFooterLink struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type storeSettingsFloatingSupportRequest struct {
	Actions []storeSettingsFloatingAction `json:"actions"`
}

type storeSettingsFloatingAction struct {
	Key     string  `json:"key"`
	Enabled bool    `json:"enabled"`
	Target  *string `json:"target"`
}

type storeSettingsVisibility struct {
	NoIndexEnabled    bool `json:"noIndexEnabled"`
	WishlistEnabled   bool `json:"wishlistEnabled"`
	CheckinEnabled    bool `json:"checkinEnabled"`
	CheckinReward     int  `json:"checkinReward"`
	RefundReclaimCard bool `json:"refundReclaimCards"`
}

type storeSettingsRegistry struct {
	Enabled bool `json:"enabled"`
	Joined  bool `json:"joined"`
	HideNav bool `json:"hideNav"`
}

type storeSettingsSiteConfigResponse struct {
	Brand           storeSettingsBrand            `json:"brand"`
	Contact         storeSettingsContact          `json:"contact"`
	Homepage        storeSettingsHomepage         `json:"homepage"`
	Footer          storeSettingsFooter           `json:"footer"`
	FloatingSupport []storeSettingsFloatingAction `json:"floatingSupport"`
	Visibility      storeSettingsVisibility       `json:"visibility"`
	Registry        storeSettingsRegistry         `json:"registry"`
}

func (r *V1) gripAdminGetStoreSettings(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_settings_not_available"})
	}

	settings, err := ext.ListSettings(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(buildStoreSettingsAdminResponse(settings)))
}

func (r *V1) gripAdminPutStoreSettingsBrand(ctx *fiber.Ctx) error {
	var body storeSettingsBrand
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if err := validateStoreSettingsBrand(body); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return r.persistStoreSettings(ctx, map[string]string{
		"shopName":        strings.TrimSpace(body.ShopName),
		"shopDescription": strings.TrimSpace(body.ShopDescription),
		"shopLogo":        strings.TrimSpace(body.ShopLogo),
		"themeColor":      strings.TrimSpace(body.ThemeColor),
	})
}

func (r *V1) gripAdminPutStoreSettingsContact(ctx *fiber.Ctx) error {
	var body storeSettingsContact
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if err := validateStoreSettingsContact(body); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return r.persistStoreSettings(ctx, map[string]string{
		"stickyBarAddress": strings.TrimSpace(body.StickyBarAddress),
		"stickyBarHotline": strings.TrimSpace(body.StickyBarHotline),
		"contactEmail":     strings.TrimSpace(body.ContactEmail),
	})
}

func (r *V1) gripAdminPutStoreSettingsHomepage(ctx *fiber.Ctx) error {
	var body storeSettingsHomepage
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if err := validateStoreSettingsHomepage(body); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	blocks, _ := json.Marshal(body.Blocks)
	return r.persistStoreSettings(ctx, map[string]string{
		"homepageBlocks":    string(blocks),
		"homepageNewsCount": strconv.Itoa(body.NewsCount),
	})
}

func (r *V1) gripAdminPutStoreSettingsFooter(ctx *fiber.Ctx) error {
	var body storeSettingsFooter
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if err := validateStoreSettingsFooter(body); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	columns, _ := json.Marshal(body.Columns)
	socialLinks, _ := json.Marshal(body.SocialLinks)
	return r.persistStoreSettings(ctx, map[string]string{
		"footerColumns":   string(columns),
		"footerCopyright": strings.TrimSpace(body.Copyright),
		"socialLinks":     string(socialLinks),
	})
}

func (r *V1) gripAdminPutStoreSettingsFloatingSupport(ctx *fiber.Ctx) error {
	var body storeSettingsFloatingSupportRequest
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if err := validateStoreSettingsFloatingSupport(body.Actions); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	actions, _ := json.Marshal(body.Actions)
	return r.persistStoreSettings(ctx, map[string]string{
		"floatingSupport": string(actions),
	})
}

func (r *V1) gripAdminPutStoreSettingsVisibility(ctx *fiber.Ctx) error {
	var body storeSettingsVisibility
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if err := validateStoreSettingsVisibility(body); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return r.persistStoreSettings(ctx, map[string]string{
		"noIndexEnabled":     strconv.FormatBool(body.NoIndexEnabled),
		"wishlistEnabled":    strconv.FormatBool(body.WishlistEnabled),
		"checkinEnabled":     strconv.FormatBool(body.CheckinEnabled),
		"checkinReward":      strconv.Itoa(body.CheckinReward),
		"refundReclaimCards": strconv.FormatBool(body.RefundReclaimCard),
	})
}

func (r *V1) gripAdminPutStoreSettingsRegistry(ctx *fiber.Ctx) error {
	var body storeSettingsRegistry
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	return r.persistStoreSettings(ctx, map[string]string{
		"registryOptIn":   strconv.FormatBool(body.Joined),
		"registryHideNav": strconv.FormatBool(body.HideNav),
	})
}

func (r *V1) gripSiteConfig(ctx *fiber.Ctx) error {
	uc, ok := r.gripCatalogUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "catalog_usecase_not_configured"})
	}

	settings, err := uc.ListPublicSettings(ctx.UserContext())
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(buildStoreSettingsSiteConfig(settings)))
}

func (r *V1) persistStoreSettings(ctx *fiber.Ctx, values map[string]string) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_settings_not_available"})
	}

	for key, value := range values {
		if err := ext.SetSetting(ctx.UserContext(), r.gripActor(ctx), key, value); err != nil {
			status, payload := mapDomainError(err)
			return ctx.Status(status).JSON(payload)
		}
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{"updated": true}))
}

func buildStoreSettingsAdminResponse(settings []entity.Setting) storeSettingsAdminResponse {
	return storeSettingsAdminResponse{
		Config:       buildStoreSettingsConfig(settings),
		Stats:        storeSettingsStats{},
		VisitorCount: 0,
	}
}

func buildStoreSettingsSiteConfig(settings []entity.Setting) storeSettingsSiteConfigResponse {
	config := buildStoreSettingsConfig(settings)
	return storeSettingsSiteConfigResponse{
		Brand:           config.Brand,
		Contact:         config.Contact,
		Homepage:        config.Homepage,
		Footer:          config.Footer,
		FloatingSupport: config.FloatSupport,
		Visibility:      config.Visibility,
		Registry:        config.Registry,
	}
}

func buildCatalogSettingsProjection(settings []entity.Setting) fiber.Map {
	config := buildStoreSettingsConfig(settings)
	return fiber.Map{
		"shopName":          config.Brand.ShopName,
		"shopDescription":   config.Brand.ShopDescription,
		"shopLogo":          emptyStringAsNil(config.Brand.ShopLogo),
		"shopFooter":        config.Footer.Copyright,
		"themeColor":        config.Brand.ThemeColor,
		"noindexEnabled":    config.Visibility.NoIndexEnabled,
		"wishlistEnabled":   config.Visibility.WishlistEnabled,
		"checkinEnabled":    config.Visibility.CheckinEnabled,
		"checkinReward":     config.Visibility.CheckinReward,
		"lowStockThreshold": 3,
		"site_name":         config.Brand.ShopName,
		"site_description":  config.Brand.ShopDescription,
		"currency":          "VND",
	}
}

func buildStoreSettingsConfig(settings []entity.Setting) storeSettingsConfig {
	values := map[string]string{}
	for _, setting := range settings {
		values[setting.Key] = setting.Value
	}

	config := storeSettingsConfig{
		Brand: storeSettingsBrand{
			ShopName:        firstNonEmpty(values["shopName"], values["site_name"], "Grip Store"),
			ShopDescription: firstNonEmpty(values["shopDescription"], values["site_description"], "High-quality virtual goods, instant delivery"),
			ShopLogo:        values["shopLogo"],
			ThemeColor:      firstNonEmpty(values["themeColor"], "purple"),
		},
		Contact: storeSettingsContact{
			StickyBarAddress: values["stickyBarAddress"],
			StickyBarHotline: values["stickyBarHotline"],
			ContactEmail:     firstNonEmpty(values["contactEmail"], values["test.support.email"]),
		},
		Homepage: storeSettingsHomepage{
			Blocks: parseJSONSetting(values["homepageBlocks"], []storeSettingsHomepageBlock{
				{Key: "hero", Enabled: true, Order: 1},
				{Key: "categories", Enabled: true, Order: 2},
				{Key: "latest_news", Enabled: true, Order: 3},
			}),
			NewsCount: parseIntSetting(values["homepageNewsCount"], 6),
		},
		Footer: storeSettingsFooter{
			Columns:     parseJSONSetting(values["footerColumns"], []storeSettingsFooterColumn{}),
			Copyright:   firstNonEmpty(values["footerCopyright"], values["shopFooter"], "Copyright © 2026 Grip Store"),
			SocialLinks: parseJSONSetting(values["socialLinks"], map[string]string{}),
		},
		FloatSupport: parseJSONSetting(values["floatingSupport"], []storeSettingsFloatingAction{
			{Key: "zalo", Enabled: false, Target: nil},
			{Key: "messenger", Enabled: false, Target: nil},
			{Key: "hotline", Enabled: false, Target: nil},
			{Key: "scroll_to_top", Enabled: true, Target: nil},
		}),
		Visibility: storeSettingsVisibility{
			NoIndexEnabled:    parseBoolSetting(values["noIndexEnabled"], false),
			WishlistEnabled:   parseBoolSetting(values["wishlistEnabled"], true),
			CheckinEnabled:    parseBoolSetting(values["checkinEnabled"], true),
			CheckinReward:     parseIntSetting(values["checkinReward"], 1),
			RefundReclaimCard: parseBoolSetting(values["refundReclaimCards"], true),
		},
		Registry: storeSettingsRegistry{
			Enabled: true,
			Joined:  parseBoolSetting(values["registryOptIn"], false),
			HideNav: parseBoolSetting(values["registryHideNav"], false),
		},
	}

	slices.SortFunc(config.Homepage.Blocks, func(a, b storeSettingsHomepageBlock) int {
		return a.Order - b.Order
	})

	return config
}

func validateStoreSettingsBrand(body storeSettingsBrand) error {
	if strings.TrimSpace(body.ShopName) == "" {
		return entity.ErrInvalidInput
	}
	if body.ShopLogo != "" && !isAbsoluteURL(body.ShopLogo) {
		return entity.ErrInvalidInput
	}
	return nil
}

func validateStoreSettingsContact(body storeSettingsContact) error {
	if body.ContactEmail != "" {
		if _, err := mail.ParseAddress(body.ContactEmail); err != nil {
			return entity.ErrInvalidInput
		}
	}
	if body.StickyBarHotline != "" && !storeSettingsPhonePattern.MatchString(strings.TrimSpace(body.StickyBarHotline)) {
		return entity.ErrInvalidInput
	}
	return nil
}

func validateStoreSettingsHomepage(body storeSettingsHomepage) error {
	if body.NewsCount < 0 {
		return entity.ErrInvalidInput
	}

	allowed := map[string]struct{}{
		"hero":              {},
		"categories":        {},
		"latest_news":       {},
		"featured_products": {},
		"testimonials":      {},
	}
	seen := map[string]struct{}{}
	for _, block := range body.Blocks {
		if _, ok := allowed[block.Key]; !ok {
			return entity.ErrInvalidInput
		}
		if _, exists := seen[block.Key]; exists {
			return entity.ErrInvalidInput
		}
		seen[block.Key] = struct{}{}
	}
	return nil
}

func validateStoreSettingsFooter(body storeSettingsFooter) error {
	for _, column := range body.Columns {
		if strings.TrimSpace(column.Title) == "" {
			return entity.ErrInvalidInput
		}
		for _, link := range column.Links {
			if strings.TrimSpace(link.Label) == "" || strings.TrimSpace(link.URL) == "" {
				return entity.ErrInvalidInput
			}
		}
	}
	for _, value := range body.SocialLinks {
		if strings.TrimSpace(value) == "" {
			continue
		}
		if !isAbsoluteURL(value) {
			return entity.ErrInvalidInput
		}
	}
	return nil
}

func validateStoreSettingsFloatingSupport(actions []storeSettingsFloatingAction) error {
	allowed := map[string]struct{}{
		"zalo":          {},
		"messenger":     {},
		"hotline":       {},
		"scroll_to_top": {},
	}
	for _, action := range actions {
		if _, ok := allowed[action.Key]; !ok {
			return entity.ErrInvalidInput
		}
		target := ""
		if action.Target != nil {
			target = strings.TrimSpace(*action.Target)
		}
		switch action.Key {
		case "scroll_to_top":
			if target != "" {
				return entity.ErrInvalidInput
			}
		case "zalo", "messenger":
			if action.Enabled && !isAbsoluteURL(target) {
				return entity.ErrInvalidInput
			}
		case "hotline":
			if action.Enabled && !storeSettingsPhonePattern.MatchString(target) {
				return entity.ErrInvalidInput
			}
		}
	}
	return nil
}

func validateStoreSettingsVisibility(body storeSettingsVisibility) error {
	if body.CheckinReward < 0 {
		return entity.ErrInvalidInput
	}
	return nil
}

func parseJSONSetting[T any](raw string, fallback T) T {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	value := fallback
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return fallback
	}
	return value
}

func parseBoolSetting(raw string, fallback bool) bool {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseIntSetting(raw string, fallback int) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func isAbsoluteURL(raw string) bool {
	parsed, err := url.Parse(raw)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}

func emptyStringAsNil(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
