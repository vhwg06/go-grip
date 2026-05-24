package v1

import (
	"net/http"
	"strconv"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type gripWishlistRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type gripReviewRequest struct {
	ProductID  string `json:"productId"`
	ProductID2 string `json:"product_id"`
	OrderID    string `json:"orderId"`
	OrderID2   string `json:"order_id"`
	Rating     int    `json:"rating"`
	Comment    string `json:"comment"`
	Content    string `json:"content"`
}

// @Summary     List wishlist
// @Description Lists wishlist items
// @ID          grip_wishlist_list
// @Tags        wishlist
// @Produce     json
// @Success     200 {object} gripListResponse
// @Router      /wishlist [get]
func (r *V1) gripWishlistList(ctx *fiber.Ctx) error {
	if r.wishlistUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "wishlist_usecase_not_configured"})
	}

	page := gripPage(ctx)
	items, total, err := r.wishlistUC.List(ctx.UserContext(), page)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	normalized := page.Normalize()
	return ctx.JSON(gripListResponse{Data: items, Meta: entity.Page{Limit: normalized.Limit, Offset: normalized.Offset, Total: total}})
}

// @Summary     Create wishlist item
// @Description Creates a wishlist item for authenticated user
// @ID          grip_wishlist_create
// @Tags        wishlist
// @Accept      json
// @Produce     json
// @Success     201 {object} envelope
// @Failure     400 {object} envelope
// @Failure     401 {object} envelope
// @Security    BearerAuth
// @Router      /wishlist [post]
func (r *V1) gripWishlistCreate(ctx *fiber.Ctx) error {
	var body gripWishlistRequest
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	item, err := r.wishlistUC.Create(ctx.UserContext(), r.gripActor(ctx), body.Title, body.Description)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.Status(http.StatusCreated).JSON(apiSuccessEnvelope(item))
}

func (r *V1) gripWishlistUpdate(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	var body gripWishlistRequest
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	item := entity.WishlistItem{
		ID:          id,
		Title:       body.Title,
		Description: body.Description,
	}
	updated, err := r.wishlistUC.Update(ctx.UserContext(), r.gripActor(ctx), item)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.JSON(apiSuccessEnvelope(updated))
}

func (r *V1) gripWishlistDelete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if err := r.wishlistUC.Delete(ctx.UserContext(), r.gripActor(ctx), id); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) gripWishlistVote(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if err := r.wishlistUC.ToggleVote(ctx.UserContext(), r.gripActor(ctx), id); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Create review
// @Description Creates a review for a delivered order
// @ID          grip_reviews_create
// @Tags        review
// @Accept      json
// @Produce     json
// @Success     201 {object} envelope
// @Failure     400 {object} envelope
// @Failure     401 {object} envelope
// @Security    BearerAuth
// @Router      /reviews [post]
func (r *V1) gripReviewCreate(ctx *fiber.Ctx) error {
	var body gripReviewRequest
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	productId := body.ProductID
	if productId == "" {
		productId = body.ProductID2
	}
	orderId := body.OrderID
	if orderId == "" {
		orderId = body.OrderID2
	}
	comment := body.Comment
	if comment == "" {
		comment = body.Content
	}

	review, err := r.wishlistUC.CreateReview(ctx.UserContext(), r.gripActor(ctx), entity.Review{
		ProductID: productId,
		OrderID:   orderId,
		Rating:    body.Rating,
		Comment:   comment,
	})
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}
	return ctx.Status(http.StatusCreated).JSON(apiSuccessEnvelope(review))
}

func (r *V1) gripReviewList(ctx *fiber.Ctx) error {
	productID := ctx.Query("product_id")
	if productID == "" {
		productID = ctx.Params("id")
	}
	if productID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(envelope{Error: "product_id_required"})
	}

	items, err := r.wishlistUC.ListReviews(ctx.UserContext(), productID)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	res := make([]fiber.Map, 0, len(items))
	for _, item := range items {
		res = append(res, fiber.Map{
			"id":         strconv.FormatInt(item.ID, 10),
			"product_id": item.ProductID,
			"rating":     item.Rating,
			"content":    item.Comment,
			"username":   item.Username,
			"created_at": item.CreatedAt.UnixMilli(),
		})
	}

	return ctx.JSON(apiSuccessEnvelope(res))
}
