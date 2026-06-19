package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func (r *V1) gripAdminListReviews(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_reviews_not_available"})
	}

	page := gripPage(ctx)
	reviews, stats, total, err := ext.ListReviews(ctx.UserContext(), r.gripActor(ctx), page, strings.TrimSpace(ctx.Query("q")), strings.TrimSpace(ctx.Query("status")))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	items := make([]fiber.Map, 0, len(reviews))
	for _, review := range reviews {
		items = append(items, fiber.Map{
			"id":                 review.ID,
			"productId":          review.ProductID,
			"productName":        review.ProductName,
			"orderId":            review.OrderID,
			"userId":             review.UserID,
			"username":           review.Username,
			"rating":             review.Rating,
			"comment":            review.Comment,
			"attachments":        review.Attachments,
			"status":             review.Status,
			"isVerifiedPurchase": review.IsVerifiedPurchase,
			"flaggedReason":      review.FlaggedReason,
			"createdAt":          review.CreatedAt.UTC().Format(timeRFC3339),
		})
	}

	normalized := page.Normalize()
	return ctx.JSON(apiSuccessEnvelope(fiber.Map{
		"reviews":  items,
		"stats":    stats,
		"total":    total,
		"page":     normalized.Offset/normalized.Limit + 1,
		"pageSize": normalized.Limit,
	}))
}

func (r *V1) gripAdminApproveReview(ctx *fiber.Ctx) error {
	return r.gripAdminUpdateReviewStatus(ctx, entity.ReviewStatusApproved)
}

func (r *V1) gripAdminHideReview(ctx *fiber.Ctx) error {
	return r.gripAdminUpdateReviewStatus(ctx, entity.ReviewStatusHidden)
}

func (r *V1) gripAdminFeatureReview(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_reviews_not_available"})
	}

	reviewID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	var body struct {
		IsFeatured *bool `json:"isFeatured"`
	}
	if err := ctx.BodyParser(&body); err != nil || body.IsFeatured == nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	nextStatus := entity.ReviewStatusApproved
	if *body.IsFeatured {
		nextStatus = entity.ReviewStatusFeatured
	}

	review, err := ext.UpdateReviewStatus(ctx.UserContext(), r.gripActor(ctx), reviewID, nextStatus)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{
		"success": true,
		"id":      review.ID,
		"status":  review.Status,
	}))
}

func (r *V1) gripAdminPublishSelectedReviews(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_reviews_not_available"})
	}

	var body struct {
		IDs []int64 `json:"ids"`
	}
	if err := ctx.BodyParser(&body); err != nil || len(body.IDs) == 0 {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	count, err := ext.BulkPublishReviews(ctx.UserContext(), r.gripActor(ctx), body.IDs)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{
		"success": true,
		"count":   count,
		"status":  entity.ReviewStatusApproved,
	}))
}

func (r *V1) gripAdminDeleteReview(ctx *fiber.Ctx) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_reviews_not_available"})
	}

	reviewID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	if err := ext.DeleteReview(ctx.UserContext(), r.gripActor(ctx), reviewID); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{
		"success": true,
		"id":      reviewID,
	}))
}

func (r *V1) gripAdminUpdateReviewStatus(ctx *fiber.Ctx, status entity.ReviewStatus) error {
	ext, ok := r.adminUC.(adminExtendedUseCase)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "admin_reviews_not_available"})
	}

	reviewID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		httpStatus, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(httpStatus).JSON(payload)
	}

	review, err := ext.UpdateReviewStatus(ctx.UserContext(), r.gripActor(ctx), reviewID, status)
	if err != nil {
		httpStatus, payload := mapDomainError(err)
		return ctx.Status(httpStatus).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{
		"success": true,
		"id":      review.ID,
		"status":  review.Status,
	}))
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"
