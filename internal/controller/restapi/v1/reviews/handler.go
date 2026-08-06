package reviews

import (
	"context"
	"errors"
	"strconv"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Reviews capability.
type Handler struct {
	logger     logger.Interface
	wishlistUC wishlistmodule.WishlistUseCase
}

// NewHandler constructs a new Reviews vertical handler instance.
func NewHandler(wishlistUC wishlistmodule.WishlistUseCase, l logger.Interface) *Handler {
	return &Handler{
		wishlistUC: wishlistUC,
		logger:     l,
	}
}

func getActor(ctx context.Context) usermodule.Actor {
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(usermodule.Actor); ok {
			return a
		}
	}
	return usermodule.Actor{}
}

// GetProductReviews handles GET /catalog/products/{id}/reviews
func (h *Handler) GetProductReviews(ctx context.Context, request openapi.GetProductReviewsRequestObject) (openapi.GetProductReviewsResponseObject, error) {
	if h.wishlistUC == nil {
		return openapi.GetProductReviews200JSONResponse([]openapi.ReviewResponse{}), nil
	}
	reviews, err := h.wishlistUC.ListReviews(ctx, request.Id)
	if err != nil {
		return openapi.GetProductReviews200JSONResponse([]openapi.ReviewResponse{}), nil
	}
	items := make([]openapi.ReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		items = append(items, reviewResponse(review))
	}
	return openapi.GetProductReviews200JSONResponse(items), nil
}

// CreateReview handles POST /reviews
func (h *Handler) CreateReview(ctx context.Context, request openapi.CreateReviewRequestObject) (openapi.CreateReviewResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.CreateReview401JSONResponse{}, nil
	}
	if request.Body == nil || request.Body.ProductId == "" {
		return openapi.CreateReview400JSONResponse{}, nil
	}

	if h.wishlistUC == nil {
		return openapi.CreateReview500JSONResponse{}, nil
	}
	created, err := h.wishlistUC.CreateReview(ctx, wishlistmodule.Actor{UserID: actor.UserID, Username: actor.Username, IsAdmin: actor.IsAdmin}, wishlistmodule.Review{ProductID: request.Body.ProductId, Rating: request.Body.Rating, Comment: request.Body.Content})
	if err != nil {
		if errors.Is(err, wishlistmodule.ErrInvalidInput) {
			return openapi.CreateReview400JSONResponse{}, nil
		}
		if errors.Is(err, wishlistmodule.ErrUnauthorized) {
			return openapi.CreateReview401JSONResponse{}, nil
		}
		return openapi.CreateReview500JSONResponse{}, nil
	}
	return openapi.CreateReview201JSONResponse(reviewResponse(created)), nil
}

// DeleteReview handles DELETE /reviews/{id}
func (h *Handler) DeleteReview(ctx context.Context, request openapi.DeleteReviewRequestObject) (openapi.DeleteReviewResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.DeleteReview401JSONResponse{}, nil
	}
	reviewID, err := strconv.ParseInt(request.Id, 10, 64)
	if err != nil || h.wishlistUC == nil {
		return openapi.DeleteReview404JSONResponse{}, nil
	}
	err = h.wishlistUC.DeleteReview(ctx, wishlistmodule.Actor{UserID: actor.UserID, Username: actor.Username, IsAdmin: actor.IsAdmin}, reviewID)
	if errors.Is(err, wishlistmodule.ErrNotFound) || errors.Is(err, wishlistmodule.ErrForbidden) {
		return openapi.DeleteReview404JSONResponse{}, nil
	}
	if err != nil {
		return openapi.DeleteReview500JSONResponse{}, nil
	}
	return openapi.DeleteReview204Response{}, nil
}

func reviewResponse(review wishlistmodule.Review) openapi.ReviewResponse {
	userID := review.UserID
	return openapi.ReviewResponse{Id: strconv.FormatInt(review.ID, 10), ProductId: review.ProductID, UserId: &userID, Rating: review.Rating, Content: review.Comment, CreatedAt: &review.CreatedAt}
}
