package admin

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// AdminListReviews handles GET /admin/reviews
func (h *Handler) AdminListReviews(ctx context.Context, request openapi.AdminListReviewsRequestObject) (openapi.AdminListReviewsResponseObject, error) {
	actor := getActor(ctx)

	page := 1
	pageSize := 20
	query := ""
	status := ""

	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}
	if request.Params.Q != nil {
		query = *request.Params.Q
	}
	if request.Params.Status != nil {
		status = *request.Params.Status
	}

	offset := (page - 1) * pageSize
	pag := pagination.New(pageSize, offset)

	reviews, total, err := h.adminUC.ListReviews(ctx, actor, pag, query, status)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminListReviews401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminListReviews403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminListReviews500JSONResponse{}, nil
		}
	}

	items := make([]openapi.AdminReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		items = append(items, toAdminReviewResponse(r))
	}
	totalInt := total
	resp := openapi.AdminReviewListResponse{
		Items: &items,
		Total: &totalInt,
	}
	return openapi.AdminListReviews200JSONResponse(resp), nil
}

// AdminPublishSelectedReviews handles POST /admin/reviews/publish-selected
func (h *Handler) AdminPublishSelectedReviews(ctx context.Context, request openapi.AdminPublishSelectedReviewsRequestObject) (openapi.AdminPublishSelectedReviewsResponseObject, error) {
	actor := getActor(ctx)

	var ids []int64
	if request.Body != nil && request.Body.Ids != nil {
		ids = request.Body.Ids
	}

	_, err := h.adminUC.BulkPublishReviews(ctx, actor, ids)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminPublishSelectedReviews401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminPublishSelectedReviews403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminPublishSelectedReviews500JSONResponse{}, nil
		}
	}
	return openapi.AdminPublishSelectedReviews200Response{}, nil
}

// AdminDeleteReview handles DELETE /admin/reviews/{reviewId}
func (h *Handler) AdminDeleteReview(ctx context.Context, request openapi.AdminDeleteReviewRequestObject) (openapi.AdminDeleteReviewResponseObject, error) {
	actor := getActor(ctx)

	if err := h.adminUC.DeleteAdminReview(ctx, actor, request.ReviewId); err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminDeleteReview401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminDeleteReview403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		case 404:
			return openapi.AdminDeleteReview404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminDeleteReview500JSONResponse{}, nil
		}
	}
	return openapi.AdminDeleteReview204Response{}, nil
}

// AdminApproveReview handles PUT /admin/reviews/{reviewId}/approve
func (h *Handler) AdminApproveReview(ctx context.Context, request openapi.AdminApproveReviewRequestObject) (openapi.AdminApproveReviewResponseObject, error) {
	actor := getActor(ctx)

	if err := h.adminUC.UpdateReviewStatus(ctx, actor, request.ReviewId, entity.ReviewStatusApproved); err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminApproveReview401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminApproveReview403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		case 404:
			return openapi.AdminApproveReview404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminApproveReview500JSONResponse{}, nil
		}
	}
	return openapi.AdminApproveReview200Response{}, nil
}

// AdminHideReview handles PUT /admin/reviews/{reviewId}/hide
func (h *Handler) AdminHideReview(ctx context.Context, request openapi.AdminHideReviewRequestObject) (openapi.AdminHideReviewResponseObject, error) {
	actor := getActor(ctx)

	if err := h.adminUC.UpdateReviewStatus(ctx, actor, request.ReviewId, entity.ReviewStatusHidden); err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminHideReview401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminHideReview403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		case 404:
			return openapi.AdminHideReview404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminHideReview500JSONResponse{}, nil
		}
	}
	return openapi.AdminHideReview200Response{}, nil
}

// AdminFeatureReview handles PUT /admin/reviews/{reviewId}/feature
func (h *Handler) AdminFeatureReview(ctx context.Context, request openapi.AdminFeatureReviewRequestObject) (openapi.AdminFeatureReviewResponseObject, error) {
	actor := getActor(ctx)

	if err := h.adminUC.UpdateReviewStatus(ctx, actor, request.ReviewId, entity.ReviewStatusFeatured); err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminFeatureReview401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminFeatureReview403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		case 404:
			return openapi.AdminFeatureReview404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminFeatureReview500JSONResponse{}, nil
		}
	}
	return openapi.AdminFeatureReview200Response{}, nil
}

// toAdminReviewResponse maps entity.Review to openapi.AdminReviewResponse.
func toAdminReviewResponse(r entity.Review) openapi.AdminReviewResponse {
	id := r.ID
	productID := r.ProductID
	productName := r.ProductName
	userID := r.UserID
	username := r.Username
	rating := r.Rating
	comment := r.Comment
	status := string(r.Status)
	verified := r.IsVerifiedPurchase

	return openapi.AdminReviewResponse{
		Id:                 &id,
		ProductId:          &productID,
		ProductName:        &productName,
		UserId:             &userID,
		Username:           &username,
		Rating:             &rating,
		Comment:            &comment,
		Status:             &status,
		IsVerifiedPurchase: &verified,
		CreatedAt:          &r.CreatedAt,
	}
}


// AdminGetReview handles GET /admin/reviews/{reviewId}
func (h *Handler) AdminGetReview(ctx context.Context, request openapi.AdminGetReviewRequestObject) (openapi.AdminGetReviewResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminGetReview401Response{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminGetReview403Response{}, nil
	}
	return openapi.AdminGetReview200Response{}, nil
}
