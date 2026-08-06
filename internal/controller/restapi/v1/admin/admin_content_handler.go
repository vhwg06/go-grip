package admin

import (
	"context"
	"sort"

	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"

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
	if h.mediaUC == nil {
		return openapi.AdminListMedia500JSONResponse{}, nil
	}
	assets, total, err := h.mediaUC.List(ctx, pagination.New(pageSize, offset), "")
	if err != nil {
		return openapi.AdminListMedia500JSONResponse{}, nil
	}

	if actor.UserID == "" {
		return openapi.AdminListMedia401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminListMedia403JSONResponse{}, nil
	}
	items := make([]openapi.AdminMediaResponse, 0, len(assets))
	for _, asset := range assets {
		id, url, fileName, mimeType, size := asset.ID, asset.URL, asset.FileName, asset.MimeType, asset.SizeBytes
		items = append(items, openapi.AdminMediaResponse{Id: &id, Url: &url, FileName: &fileName, MimeType: &mimeType, Size: &size, CreatedAt: &asset.CreatedAt})
	}
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
func (h *Handler) AdminGetPresignedUrl(ctx context.Context, request openapi.AdminGetPresignedUrlRequestObject) (openapi.AdminGetPresignedUrlResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminGetPresignedUrl401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminGetPresignedUrl403JSONResponse{}, nil
	}
	if h.mediaUC == nil {
		return openapi.AdminGetPresignedUrl500JSONResponse{}, nil
	}
	contentType := "image/png"
	if request.Params.ContentType != nil && *request.Params.ContentType != "" {
		contentType = *request.Params.ContentType
	}
	uploadURL, publicURL, fileID, err := h.mediaUC.GeneratePresignedURL(ctx, request.Params.FileName, contentType)
	if err != nil {
		return openapi.AdminGetPresignedUrl500JSONResponse{}, nil
	}
	return openapi.AdminGetPresignedUrl200JSONResponse{
		Id:        &fileID,
		PublicUrl: &publicURL,
		UploadUrl: &uploadURL,
		Url:       &publicURL,
		Key:       &fileID,
	}, nil
}

// AdminListBanners handles GET /admin/banners
func (h *Handler) AdminListBanners(ctx context.Context, _ openapi.AdminListBannersRequestObject) (openapi.AdminListBannersResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminListBanners401JSONResponse{}, nil
	}
	if !actor.IsAdmin || h.homepageUC == nil {
		return openapi.AdminListBanners403JSONResponse{}, nil
	}
	blocks, err := h.homepageUC.ListBlocks(ctx, false)
	if err != nil {
		return openapi.AdminListBanners500JSONResponse{}, nil
	}
	items := make([]openapi.AdminBannerResponse, 0)
	for _, block := range blocks {
		if block.BlockType == "banner" {
			items = append(items, bannerResponse(block))
		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].SortOrder != nil && items[j].SortOrder != nil && *items[i].SortOrder < *items[j].SortOrder
	})
	resp := openapi.AdminBannerListResponse{Items: &items}
	return openapi.AdminListBanners200JSONResponse(resp), nil
}

// AdminSaveBanner handles POST /admin/banners
func (h *Handler) AdminSaveBanner(ctx context.Context, request openapi.AdminSaveBannerRequestObject) (openapi.AdminSaveBannerResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminSaveBanner401JSONResponse{}, nil
	}
	if !actor.IsAdmin || h.homepageUC == nil || request.Body == nil {
		return openapi.AdminSaveBanner403JSONResponse{}, nil
	}
	block := bannerBlock(*request.Body)
	if request.Body.Id != nil && *request.Body.Id != "" {
		block.ID = *request.Body.Id
		if _, err := h.homepageUC.UpdateBlock(ctx, block); err != nil {
			return openapi.AdminSaveBanner500JSONResponse{}, nil
		}
	} else {
		created, err := h.homepageUC.StoreBlock(ctx, block)
		if err != nil {
			return openapi.AdminSaveBanner500JSONResponse{}, nil
		}
		block = created
	}
	resp := bannerResponse(block)
	return openapi.AdminSaveBanner201JSONResponse(resp), nil
}

// AdminListFaqs handles GET /admin/faqs
func (h *Handler) AdminListFaqs(ctx context.Context, _ openapi.AdminListFaqsRequestObject) (openapi.AdminListFaqsResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminListFaqs401JSONResponse{}, nil
	}
	if !actor.IsAdmin || h.homepageUC == nil {
		return openapi.AdminListFaqs403JSONResponse{}, nil
	}
	blocks, err := h.homepageUC.ListBlocks(ctx, false)
	if err != nil {
		return openapi.AdminListFaqs500JSONResponse{}, nil
	}
	items := make([]openapi.AdminFaqResponse, 0)
	for _, block := range blocks {
		if block.BlockType == "faq" {
			items = append(items, faqResponse(block))
		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].SortOrder != nil && items[j].SortOrder != nil && *items[i].SortOrder < *items[j].SortOrder
	})
	resp := openapi.AdminFaqListResponse{Items: &items}
	return openapi.AdminListFaqs200JSONResponse(resp), nil
}

// AdminSaveFaq handles POST /admin/faqs
func (h *Handler) AdminSaveFaq(ctx context.Context, request openapi.AdminSaveFaqRequestObject) (openapi.AdminSaveFaqResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminSaveFaq401JSONResponse{}, nil
	}
	if !actor.IsAdmin || h.homepageUC == nil || request.Body == nil {
		return openapi.AdminSaveFaq403JSONResponse{}, nil
	}
	block := faqBlock(*request.Body)
	if request.Body.Id != nil && *request.Body.Id != "" {
		block.ID = *request.Body.Id
		if _, err := h.homepageUC.UpdateBlock(ctx, block); err != nil {
			return openapi.AdminSaveFaq500JSONResponse{}, nil
		}
	} else {
		created, err := h.homepageUC.StoreBlock(ctx, block)
		if err != nil {
			return openapi.AdminSaveFaq500JSONResponse{}, nil
		}
		block = created
	}
	resp := faqResponse(block)
	return openapi.AdminSaveFaq201JSONResponse(resp), nil
}

func bannerBlock(request openapi.AdminBannerRequest) contentmodule.HomepageBlock {
	config := map[string]any{}
	if request.Title != nil {
		config["title"] = *request.Title
	}
	if request.Subtitle != nil {
		config["subtitle"] = *request.Subtitle
	}
	if request.Image != nil {
		config["image"] = *request.Image
	}
	if request.ImageUrl != nil {
		config["imageUrl"] = *request.ImageUrl
	}
	if request.MobileImage != nil {
		config["mobileImage"] = *request.MobileImage
	}
	if request.CtaText != nil {
		config["ctaText"] = *request.CtaText
	}
	if request.CtaLink != nil {
		config["ctaLink"] = *request.CtaLink
	}
	if request.LinkUrl != nil {
		config["linkUrl"] = *request.LinkUrl
	}
	active := false
	if request.IsActive != nil {
		active = *request.IsActive
	}
	position := 0
	if request.SortOrder != nil {
		position = *request.SortOrder
	}
	return contentmodule.HomepageBlock{ID: stringValue(config, "id"), BlockType: "banner", Config: config, Position: position, IsActive: active}
}

func bannerResponse(block contentmodule.HomepageBlock) openapi.AdminBannerResponse {
	config := block.Config
	id, title, subtitle, image, imageURL, mobileImage := block.ID, mapString(config, "title"), mapString(config, "subtitle"), mapString(config, "image"), mapString(config, "imageUrl"), mapString(config, "mobileImage")
	ctaText, ctaLink, linkURL := mapString(config, "ctaText"), mapString(config, "ctaLink"), mapString(config, "linkUrl")
	active, position := block.IsActive, block.Position
	return openapi.AdminBannerResponse{Id: &id, Title: &title, Subtitle: &subtitle, Image: &image, ImageUrl: &imageURL, MobileImage: &mobileImage, CtaText: &ctaText, CtaLink: &ctaLink, LinkUrl: &linkURL, IsActive: &active, SortOrder: &position}
}

func faqBlock(request openapi.AdminFaqRequest) contentmodule.HomepageBlock {
	config := map[string]any{}
	if request.Question != nil {
		config["question"] = *request.Question
	}
	if request.Answer != nil {
		config["answer"] = *request.Answer
	}
	active := false
	if request.IsActive != nil {
		active = *request.IsActive
	}
	position := 0
	if request.SortOrder != nil {
		position = *request.SortOrder
	}
	return contentmodule.HomepageBlock{BlockType: "faq", Config: config, Position: position, IsActive: active}
}

func faqResponse(block contentmodule.HomepageBlock) openapi.AdminFaqResponse {
	config := block.Config
	id, question, answer := block.ID, mapString(config, "question"), mapString(config, "answer")
	active, position := block.IsActive, block.Position
	return openapi.AdminFaqResponse{Id: &id, Question: &question, Answer: &answer, IsActive: &active, SortOrder: &position}
}

func mapString(values map[string]any, key string) string {
	if value, ok := values[key].(string); ok {
		return value
	}
	return ""
}

func stringValue(values map[string]any, key string) string { return mapString(values, key) }

// AdminDeleteBanner handles DELETE /admin/banners/{id}.
func (h *Handler) AdminDeleteBanner(ctx context.Context, request openapi.AdminDeleteBannerRequestObject) (openapi.AdminDeleteBannerResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminDeleteBanner401JSONResponse{}, nil
	}
	if !actor.IsAdmin || h.homepageUC == nil {
		return openapi.AdminDeleteBanner403JSONResponse{}, nil
	}
	if err := h.homepageUC.DeleteBlock(ctx, request.Id); err != nil {
		return openapi.AdminDeleteBanner500JSONResponse{}, nil
	}
	return openapi.AdminDeleteBanner204Response{}, nil
}

// AdminDeleteFaq handles DELETE /admin/faqs/{id}.
func (h *Handler) AdminDeleteFaq(ctx context.Context, request openapi.AdminDeleteFaqRequestObject) (openapi.AdminDeleteFaqResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminDeleteFaq401JSONResponse{}, nil
	}
	if !actor.IsAdmin || h.homepageUC == nil {
		return openapi.AdminDeleteFaq403JSONResponse{}, nil
	}
	if err := h.homepageUC.DeleteBlock(ctx, request.Id); err != nil {
		return openapi.AdminDeleteFaq500JSONResponse{}, nil
	}
	return openapi.AdminDeleteFaq204Response{}, nil
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
