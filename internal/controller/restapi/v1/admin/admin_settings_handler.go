package admin

import (
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"

	"context"
	"fmt"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/gofiber/fiber/v2"
)

// AdminGetSetting handles GET /admin/settings/{key}
func (h *Handler) AdminGetSetting(ctx context.Context, request openapi.AdminGetSettingRequestObject) (openapi.AdminGetSettingResponseObject, error) {
	actor := getActor(ctx)

	settings, err := h.adminUC.ListSettings(ctx, actor)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminGetSetting401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminGetSetting403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminGetSetting500JSONResponse{}, nil
		}
	}

	for _, s := range settings {
		if s.Key == request.Key {
			return openapi.AdminGetSetting200JSONResponse(toAdminSettingResponse(s)), nil
		}
	}

	_, errResp := mapAdminError(usermodule.ErrNotFound)
	return openapi.AdminGetSetting404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
}

// AdminUpsertSetting handles PUT /admin/settings/{key}
func (h *Handler) AdminUpsertSetting(ctx context.Context, request openapi.AdminUpsertSettingRequestObject) (openapi.AdminUpsertSettingResponseObject, error) {
	actor := getActor(ctx)

	value := ""
	if request.Body != nil {
		if request.Body.Value != nil {
			value = *request.Body.Value
		}
		if request.Body.ArticleId != nil {
			value = *request.Body.ArticleId
		}
	}

	if err := h.adminUC.UpsertSetting(ctx, actor, request.Key, value); err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminUpsertSetting401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminUpsertSetting403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminUpsertSetting500JSONResponse{}, nil
		}
	}

	now := time.Now()
	key := request.Key
	resp := openapi.AdminSettingResponse{Key: &key, Value: &value, UpdatedAt: &now}
	return openapi.AdminUpsertSetting200JSONResponse(resp), nil
}

// AdminGetStoreSettings handles GET /admin/store-settings
func (h *Handler) AdminGetStoreSettings(ctx context.Context, _ openapi.AdminGetStoreSettingsRequestObject) (openapi.AdminGetStoreSettingsResponseObject, error) {
	actor := getActor(ctx)

	settings, err := h.adminUC.ListSettings(ctx, actor)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminGetStoreSettings401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminGetStoreSettings403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminGetStoreSettings500JSONResponse{}, nil
		}
	}

	m := settingsProjection(settings)
	return openapi.AdminGetStoreSettings200JSONResponse(m), nil
}

// AdminUpdateStoreSettingsBrand handles PUT /admin/store-settings/brand
func (h *Handler) AdminUpdateStoreSettingsBrand(ctx context.Context, request openapi.AdminUpdateStoreSettingsBrandRequestObject) (openapi.AdminUpdateStoreSettingsBrandResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUpdateStoreSettingsBrand401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUpdateStoreSettingsBrand403JSONResponse{}, nil
	}
	if request.Body != nil {
		if err := h.replaceSettingsGroup(ctx, actor, "brand.", map[string]interface{}(*request.Body)); err != nil {
			return openapi.AdminUpdateStoreSettingsBrand500JSONResponse{}, nil
		}
		for k, v := range *request.Body {
			if err := h.adminUC.UpsertSetting(ctx, actor, "brand."+k, anyToString(v)); err != nil {
				return openapi.AdminUpdateStoreSettingsBrand500JSONResponse{}, nil
			}
		}
	}
	return openapi.AdminUpdateStoreSettingsBrand200Response{}, nil
}

// AdminUpdateStoreSettingsContact handles PUT /admin/store-settings/contact
func (h *Handler) AdminUpdateStoreSettingsContact(ctx context.Context, request openapi.AdminUpdateStoreSettingsContactRequestObject) (openapi.AdminUpdateStoreSettingsContactResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUpdateStoreSettingsContact401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUpdateStoreSettingsContact403JSONResponse{}, nil
	}
	if request.Body != nil {
		for k, v := range *request.Body {
			if strings.Contains(strings.ToLower(k), "email") && !validEmail(anyToString(v)) {
				return adminContactBadRequestResponse{}, nil
			}
		}
		if err := h.replaceSettingsGroup(ctx, actor, "contact.", map[string]interface{}(*request.Body)); err != nil {
			return openapi.AdminUpdateStoreSettingsContact500JSONResponse{}, nil
		}
		for k, v := range *request.Body {
			if err := h.adminUC.UpsertSetting(ctx, actor, "contact."+k, anyToString(v)); err != nil {
				return openapi.AdminUpdateStoreSettingsContact500JSONResponse{}, nil
			}
		}
	}
	return openapi.AdminUpdateStoreSettingsContact200Response{}, nil
}

func (h *Handler) replaceSettingsGroup(ctx context.Context, actor usermodule.Actor, prefix string, body map[string]interface{}) error {
	settings, err := h.adminUC.ListSettings(ctx, actor)
	if err != nil {
		return err
	}
	keep := make(map[string]struct{}, len(body))
	for key := range body {
		keep[prefix+key] = struct{}{}
	}
	for _, setting := range settings {
		if strings.HasPrefix(setting.Key, prefix) {
			if _, ok := keep[setting.Key]; !ok {
				if err := h.adminUC.DeleteSetting(ctx, actor, setting.Key); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// AdminUpdateStoreSettingsFooter handles PUT /admin/store-settings/footer
func (h *Handler) AdminUpdateStoreSettingsFooter(ctx context.Context, request openapi.AdminUpdateStoreSettingsFooterRequestObject) (openapi.AdminUpdateStoreSettingsFooterResponseObject, error) {
	actor := getActor(ctx)
	if request.Body != nil {
		body := *request.Body
		if columns, ok := body["columns"].([]any); ok {
			seenTitles := make(map[string]bool)
			for _, col := range columns {
				if colMap, ok := col.(map[string]any); ok {
					if title, ok := colMap["title"].(string); ok {
						if seenTitles[title] {
							var errResp openapi.ErrorResponse
							_ = errResp.Error.FromErrorResponseError0("duplicate column title")
							return openapi.AdminUpdateStoreSettingsFooter400JSONResponse{
								BadRequestResponseJSONResponse: openapi.BadRequestResponseJSONResponse(errResp),
							}, nil
						}
						seenTitles[title] = true
					}
				}
			}
		}
		for k, v := range body {
			_ = h.adminUC.UpsertSetting(ctx, actor, "footer."+k, anyToString(v))
		}
	}
	return openapi.AdminUpdateStoreSettingsFooter200Response{}, nil
}

// AdminUpdateStoreSettingsHomepage handles PUT /admin/store-settings/homepage
func (h *Handler) AdminUpdateStoreSettingsHomepage(ctx context.Context, request openapi.AdminUpdateStoreSettingsHomepageRequestObject) (openapi.AdminUpdateStoreSettingsHomepageResponseObject, error) {
	actor := getActor(ctx)
	if request.Body != nil {
		body := *request.Body
		if newsCount, ok := body["newsCount"].(float64); ok && newsCount < 0 {
			var errResp openapi.ErrorResponse
			_ = errResp.Error.FromErrorResponseError0("invalid newsCount")
			return openapi.AdminUpdateStoreSettingsHomepage400JSONResponse{
				BadRequestResponseJSONResponse: openapi.BadRequestResponseJSONResponse(errResp),
			}, nil
		}
		if newsCount, ok := body["newsCount"].(int); ok && newsCount < 0 {
			var errResp openapi.ErrorResponse
			_ = errResp.Error.FromErrorResponseError0("invalid newsCount")
			return openapi.AdminUpdateStoreSettingsHomepage400JSONResponse{
				BadRequestResponseJSONResponse: openapi.BadRequestResponseJSONResponse(errResp),
			}, nil
		}
		if blocks, ok := body["blocks"].([]any); ok {
			seenKeys := make(map[string]bool)
			seenOrders := make(map[int]bool)
			for _, b := range blocks {
				if blockMap, ok := b.(map[string]any); ok {
					if k, ok := blockMap["key"].(string); ok {
						if seenKeys[k] {
							var errResp openapi.ErrorResponse
							_ = errResp.Error.FromErrorResponseError0("duplicate block key")
							return openapi.AdminUpdateStoreSettingsHomepage400JSONResponse{
								BadRequestResponseJSONResponse: openapi.BadRequestResponseJSONResponse(errResp),
							}, nil
						}
						seenKeys[k] = true
					}
					if ord, ok := blockMap["order"].(float64); ok {
						if seenOrders[int(ord)] {
							var errResp openapi.ErrorResponse
							_ = errResp.Error.FromErrorResponseError0("duplicate block order")
							return openapi.AdminUpdateStoreSettingsHomepage400JSONResponse{
								BadRequestResponseJSONResponse: openapi.BadRequestResponseJSONResponse(errResp),
							}, nil
						}
						seenOrders[int(ord)] = true
					}
					if ord, ok := blockMap["priority"].(float64); ok {
						if seenOrders[int(ord)] {
							var errResp openapi.ErrorResponse
							_ = errResp.Error.FromErrorResponseError0("duplicate block priority")
							return openapi.AdminUpdateStoreSettingsHomepage400JSONResponse{
								BadRequestResponseJSONResponse: openapi.BadRequestResponseJSONResponse(errResp),
							}, nil
						}
						seenOrders[int(ord)] = true
					}
				}
			}
		}
		for k, v := range body {
			_ = h.adminUC.UpsertSetting(ctx, actor, "homepage."+k, anyToString(v))
		}
	}
	return openapi.AdminUpdateStoreSettingsHomepage200Response{}, nil
}

// AdminUpdateStoreSettingsFloatingSupport handles PUT /admin/store-settings/floating-support
func (h *Handler) AdminUpdateStoreSettingsFloatingSupport(ctx context.Context, request openapi.AdminUpdateStoreSettingsFloatingSupportRequestObject) (openapi.AdminUpdateStoreSettingsFloatingSupportResponseObject, error) {
	actor := getActor(ctx)
	if request.Body != nil {
		body := *request.Body
		if actions, ok := body["actions"].([]any); ok {
			for _, raw := range actions {
				action, ok := raw.(map[string]any)
				if !ok {
					return adminFloatingSupportBadRequestResponse{}, nil
				}
				target, present := action["target"]
				if !present || target == nil || strings.TrimSpace(fmt.Sprint(target)) == "" {
					continue
				}
				parsed, err := url.Parse(strings.TrimSpace(fmt.Sprint(target)))
				if err != nil || parsed.Scheme == "" || parsed.Host == "" {
					return adminFloatingSupportBadRequestResponse{}, nil
				}
			}
		}
		for k, v := range *request.Body {
			if err := h.adminUC.UpsertSetting(ctx, actor, "floating_support."+k, anyToString(v)); err != nil {
				return openapi.AdminUpdateStoreSettingsFloatingSupport500JSONResponse{}, nil
			}
		}
	}
	return openapi.AdminUpdateStoreSettingsFloatingSupport200Response{}, nil
}

func validEmail(value string) bool {
	address, err := mail.ParseAddress(strings.TrimSpace(value))
	return err == nil && address.Address == strings.TrimSpace(value) && strings.Contains(address.Address, "@")
}

func settingsProjection(settings []catalogmodule.Setting) openapi.AdminStoreSettingsResponse {
	result := make(openapi.AdminStoreSettingsResponse, len(settings)+3)
	brand := map[string]any{}
	contact := map[string]any{}
	for _, setting := range settings {
		result[setting.Key] = setting.Value
		if strings.HasPrefix(setting.Key, "brand.") {
			brand[strings.TrimPrefix(setting.Key, "brand.")] = setting.Value
		}
		if strings.HasPrefix(setting.Key, "contact.") {
			contact[strings.TrimPrefix(setting.Key, "contact.")] = setting.Value
		}
	}
	result["config"] = map[string]any{"brand": brand, "contact": contact}
	result["stats"] = map[string]any{"settingsCount": len(settings)}
	result["visitorCount"] = 0
	return result
}

type adminContactBadRequestResponse struct{}

func (adminContactBadRequestResponse) VisitAdminUpdateStoreSettingsContactResponse(ctx *fiber.Ctx) error {
	ctx.Status(400)
	return nil
}

type adminFloatingSupportBadRequestResponse struct{}

func (adminFloatingSupportBadRequestResponse) VisitAdminUpdateStoreSettingsFloatingSupportResponse(ctx *fiber.Ctx) error {
	ctx.Status(400)
	return nil
}

// toAdminSettingResponse maps catalogmodule.Setting to openapi.AdminSettingResponse.
func toAdminSettingResponse(s catalogmodule.Setting) openapi.AdminSettingResponse {
	key := s.Key
	value := s.Value
	now := time.Now()
	return openapi.AdminSettingResponse{Key: &key, Value: &value, UpdatedAt: &now}
}

// anyToString converts an interface{} value to its string representation.
// Only string and numeric types are handled; all others produce an empty string.
func anyToString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case int:
		return strconv.Itoa(t)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	default:
		return ""
	}
}

// AdminUpdateCollect handles PUT /admin/collect/setup
func (h *Handler) AdminUpdateCollect(ctx context.Context, request openapi.AdminUpdateCollectRequestObject) (openapi.AdminUpdateCollectResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUpdateCollect401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUpdateCollect403JSONResponse{}, nil
	}
	return openapi.AdminUpdateCollect200Response{}, nil
}
