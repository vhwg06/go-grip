package reviews

import (
	"strconv"
	"context"
	"time"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Reviews capability.
type Handler struct {
	logger logger.Interface
}

// NewHandler constructs a new Reviews vertical handler instance.
func NewHandler(l logger.Interface) *Handler {
	return &Handler{
		logger: l,
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
	_ = request
	return openapi.GetProductReviews200JSONResponse([]openapi.ReviewResponse{}), nil
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

	now := time.Now().UTC()
	userID := actor.UserID
	rev := openapi.ReviewResponse{
		Id:        strconv.FormatInt(time.Now().UnixMilli(), 10),
		ProductId: request.Body.ProductId,
		UserId:    &userID,
		Rating:    request.Body.Rating,
		Content:   request.Body.Content,
		CreatedAt: &now,
	}

	return openapi.CreateReview201JSONResponse(rev), nil
}

// DeleteReview handles DELETE /reviews/{id}
func (h *Handler) DeleteReview(ctx context.Context, request openapi.DeleteReviewRequestObject) (openapi.DeleteReviewResponseObject, error) {
	_ = request
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.DeleteReview401JSONResponse{}, nil
	}

	return openapi.DeleteReview204Response{}, nil
}
