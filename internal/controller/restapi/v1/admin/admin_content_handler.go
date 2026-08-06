package admin

import (
	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"

	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// AdminListMedia handles GET /admin/media
func (h *Handler) AdminListMedia(ctx context.Context, request openapi.AdminListMediaRequestObject) (openapi.AdminListMediaResponseObject, error) {
	actor := getActor(ctx)

	page := 1
	pageSize := 24
	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}
	offset := (page - 1) * pageSize
	_, total, err := h.adminUC.ListProducts(ctx, actor, pagination.New(pageSize, offset))
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminListMedia401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminListMedia403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminListMedia500JSONResponse{}, nil
		}
	}

	items := []openapi.AdminMediaResponse{}
	resp := openapi.AdminMediaListResponse{Items: &items, Total: &total}
	return openapi.AdminListMedia200JSONResponse(resp), nil
}

// AdminCreateMedia handles POST /admin/media
func (h *Handler) AdminCreateMedia(ctx context.Context, _ openapi.AdminCreateMediaRequestObject) (openapi.AdminCreateMediaResponseObject, error) {
	actor := getActor(ctx)
	_, _, err := h.adminUC.ListUsers(ctx, actor, pagination.New(1, 0))
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminCreateMedia401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminCreateMedia403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminCreateMedia500JSONResponse{}, nil
		}
	}
	resp := openapi.AdminMediaResponse{}
	return openapi.AdminCreateMedia201JSONResponse(resp), nil
}

// AdminGetPresignedUrl handles GET /admin/media/presigned
func (h *Handler) AdminGetPresignedUrl(ctx context.Context, _ openapi.AdminGetPresignedUrlRequestObject) (openapi.AdminGetPresignedUrlResponseObject, error) {
	actor := getActor(ctx)
	_, _, err := h.adminUC.ListUsers(ctx, actor, pagination.New(1, 0))
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminGetPresignedUrl401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminGetPresignedUrl403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminGetPresignedUrl500JSONResponse{}, nil
		}
	}
	resp := openapi.AdminPresignedUrlResponse{}
	return openapi.AdminGetPresignedUrl200JSONResponse(resp), nil
}

// AdminListBanners handles GET /admin/banners
func (h *Handler) AdminListBanners(ctx context.Context, _ openapi.AdminListBannersRequestObject) (openapi.AdminListBannersResponseObject, error) {
	actor := getActor(ctx)

	settings, err := h.adminUC.ListSettings(ctx, actor)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminListBanners401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminListBanners403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminListBanners500JSONResponse{}, nil
		}
	}

	// Banners are stored as settings with "banner.*" prefix keys.
	items := make([]openapi.AdminBannerResponse, 0)
	for _, s := range settings {
		if len(s.Key) > 7 && s.Key[:7] == "banner." {
			title := s.Key[7:]
			val := s.Value
			items = append(items, openapi.AdminBannerResponse{Title: &title, ImageUrl: &val})
		}
	}
	resp := openapi.AdminBannerListResponse{Items: &items}
	return openapi.AdminListBanners200JSONResponse(resp), nil
}

// AdminSaveBanner handles POST /admin/banners
func (h *Handler) AdminSaveBanner(ctx context.Context, request openapi.AdminSaveBannerRequestObject) (openapi.AdminSaveBannerResponseObject, error) {
	actor := getActor(ctx)

	title := ""
	imageURL := ""
	if request.Body != nil {
		if request.Body.Title != nil {
			title = *request.Body.Title
		}
		if request.Body.ImageUrl != nil {
			imageURL = *request.Body.ImageUrl
		}
	}

	if err := h.adminUC.UpsertSetting(ctx, actor, "banner."+title, imageURL); err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminSaveBanner401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminSaveBanner403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminSaveBanner500JSONResponse{}, nil
		}
	}

	resp := openapi.AdminBannerResponse{Title: &title, ImageUrl: &imageURL}
	return openapi.AdminSaveBanner201JSONResponse(resp), nil
}

// AdminListFaqs handles GET /admin/faqs
func (h *Handler) AdminListFaqs(ctx context.Context, _ openapi.AdminListFaqsRequestObject) (openapi.AdminListFaqsResponseObject, error) {
	actor := getActor(ctx)

	settings, err := h.adminUC.ListSettings(ctx, actor)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminListFaqs401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminListFaqs403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminListFaqs500JSONResponse{}, nil
		}
	}

	// FAQs are stored as settings with "faq.*" prefix keys.
	items := make([]openapi.AdminFaqResponse, 0)
	for _, s := range settings {
		if len(s.Key) > 4 && s.Key[:4] == "faq." {
			question := s.Key[4:]
			answer := s.Value
			items = append(items, openapi.AdminFaqResponse{Question: &question, Answer: &answer})
		}
	}
	resp := openapi.AdminFaqListResponse{Items: &items}
	return openapi.AdminListFaqs200JSONResponse(resp), nil
}

// AdminSaveFaq handles POST /admin/faqs
func (h *Handler) AdminSaveFaq(ctx context.Context, request openapi.AdminSaveFaqRequestObject) (openapi.AdminSaveFaqResponseObject, error) {
	actor := getActor(ctx)

	question := ""
	answer := ""
	if request.Body != nil {
		if request.Body.Question != nil {
			question = *request.Body.Question
		}
		if request.Body.Answer != nil {
			answer = *request.Body.Answer
		}
	}

	if err := h.adminUC.UpsertSetting(ctx, actor, "faq."+question, answer); err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminSaveFaq401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminSaveFaq403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminSaveFaq500JSONResponse{}, nil
		}
	}

	resp := openapi.AdminFaqResponse{Question: &question, Answer: &answer}
	return openapi.AdminSaveFaq201JSONResponse(resp), nil
}

// AdminListMessages handles GET /admin/messages
func (h *Handler) AdminListMessages(ctx context.Context, _ openapi.AdminListMessagesRequestObject) (openapi.AdminListMessagesResponseObject, error) {
	actor := getActor(ctx)

	msgs, err := h.adminUC.ListMessages(ctx, actor)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminListMessages401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminListMessages403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminListMessages500JSONResponse{}, nil
		}
	}

	items := make([]openapi.AdminMessageResponse, 0, len(msgs))
	for _, m := range msgs {
		items = append(items, toAdminMessageResponse(m))
	}
	resp := openapi.AdminMessageListResponse{Items: &items}
	return openapi.AdminListMessages200JSONResponse(resp), nil
}

// AdminBroadcastMessage handles POST /admin/messages/broadcast
func (h *Handler) AdminBroadcastMessage(ctx context.Context, request openapi.AdminBroadcastMessageRequestObject) (openapi.AdminBroadcastMessageResponseObject, error) {
	actor := getActor(ctx)

	title := ""
	body := ""
	imageURL := ""
	if request.Body != nil {
		title = request.Body.Title
		body = request.Body.Body
		if request.Body.ImageUrl != nil {
			imageURL = *request.Body.ImageUrl
		}
	}

	msg, err := h.adminUC.BroadcastMessage(ctx, actor, title, body, imageURL)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminBroadcastMessage401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminBroadcastMessage403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminBroadcastMessage500JSONResponse{}, nil
		}
	}

	resp := toAdminMessageResponse(msg)
	return openapi.AdminBroadcastMessage200JSONResponse(resp), nil
}

// AdminGetNotifications handles GET /admin/notifications
func (h *Handler) AdminGetNotifications(ctx context.Context, _ openapi.AdminGetNotificationsRequestObject) (openapi.AdminGetNotificationsResponseObject, error) {
	actor := getActor(ctx)
	_, _, err := h.adminUC.ListUsers(ctx, actor, pagination.New(1, 0))
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminGetNotifications401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminGetNotifications403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminGetNotifications500JSONResponse{}, nil
		}
	}
	resp := openapi.AdminNotificationStatusResponse{"status": "ready"}
	return openapi.AdminGetNotifications200JSONResponse(resp), nil
}

// AdminListProducts handles GET /admin/products
func (h *Handler) AdminListProducts(ctx context.Context, request openapi.AdminListProductsRequestObject) (openapi.AdminListProductsResponseObject, error) {
	actor := getActor(ctx)

	limit := 20
	offset := 0
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}
	if request.Params.Offset != nil {
		offset = *request.Params.Offset
	}
	pag := pagination.New(limit, offset)

	products, total, err := h.adminUC.ListProducts(ctx, actor, pag)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminListProducts401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminListProducts403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminListProducts500JSONResponse{}, nil
		}
	}

	items := make([]openapi.ProductResponse, 0, len(products))
	for _, p := range products {
		pid := p.ID
		ptitle := p.Title
		pprice := int(p.Price)
		items = append(items, openapi.ProductResponse{Id: pid, Title: ptitle, Price: pprice})
	}
	resp := openapi.AdminProductListResponse{Items: &items, Total: &total}
	return openapi.AdminListProducts200JSONResponse(resp), nil
}

// AdminListCategories handles GET /admin/categories
func (h *Handler) AdminListCategories(ctx context.Context, _ openapi.AdminListCategoriesRequestObject) (openapi.AdminListCategoriesResponseObject, error) {
	actor := getActor(ctx)

	categories, err := h.adminUC.ListAdminCategories(ctx, actor)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminListCategories401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminListCategories403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminListCategories500JSONResponse{}, nil
		}
	}

	items := make([]openapi.CategoryResponse, 0, len(categories))
	for _, c := range categories {
		items = append(items, openapi.CategoryResponse{Id: c.ID, Name: c.Name})
	}
	resp := openapi.AdminListCategories200JSONResponse{Items: &items}
	return resp, nil
}

// AdminUpdateProductEditorial handles PATCH /admin/products/{id}
func (h *Handler) AdminUpdateProductEditorial(ctx context.Context, request openapi.AdminUpdateProductEditorialRequestObject) (openapi.AdminUpdateProductEditorialResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUpdateProductEditorial401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUpdateProductEditorial403JSONResponse{}, nil
	}

	return openapi.AdminUpdateProductEditorial200Response{}, nil
}


// toAdminMessageResponse maps notificationmodule.AdminMessage to openapi.AdminMessageResponse.
func toAdminMessageResponse(m notificationmodule.AdminMessage) openapi.AdminMessageResponse {
	id := m.ID
	msgType := m.TargetType
	title := m.Title
	body := m.Body
	targetUserID := m.TargetValue
	return openapi.AdminMessageResponse{
		Id:           &id,
		Type:         &msgType,
		Title:        &title,
		Body:         &body,
		TargetUserID: &targetUserID,
		CreatedAt:    &m.CreatedAt,
	}
}
