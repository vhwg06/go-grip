package wishlist

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Wishlist capability.
type Handler struct {
	wishlistUC wishlistmodule.WishlistUseCase
	logger     logger.Interface
}

// NewHandler constructs a new Wishlist vertical handler instance.
func NewHandler(wishlistUC wishlistmodule.WishlistUseCase, l logger.Interface) *Handler {
	return &Handler{
		wishlistUC: wishlistUC,
		logger:     l,
	}
}

// GetMyWishlist handles GET /wishlist
func (h *Handler) GetMyWishlist(ctx context.Context, request openapi.GetMyWishlistRequestObject) (openapi.GetMyWishlistResponseObject, error) {
	actor := usermodule.Actor{UserID: "usr-1"}
	items, _, err := h.wishlistUC.List(ctx, pagination.Pagination{Limit: 100})
	if err != nil {
		status, errResp := mapWishlistError(err)
		switch status {
		case 401:
			return openapi.GetMyWishlist401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetMyWishlist500JSONResponse{}, nil
		}
	}

	wishlistDTO := toWishlistResponse(actor.UserID, items)
	return openapi.GetMyWishlist200JSONResponse(wishlistDTO), nil
}

// AddToWishlist handles POST /wishlist/items
func (h *Handler) AddToWishlist(ctx context.Context, request openapi.AddToWishlistRequestObject) (openapi.AddToWishlistResponseObject, error) {
	if request.Body == nil {
		return openapi.AddToWishlist400JSONResponse{}, nil
	}

	actor := usermodule.Actor{UserID: "usr-1"}
	item, err := h.wishlistUC.Create(ctx, actor, request.Body.ProductId, "")
	if err != nil {
		status, errResp := mapWishlistError(err)
		switch status {
		case 400:
			return openapi.AddToWishlist400JSONResponse{}, nil
		case 401:
			return openapi.AddToWishlist401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.AddToWishlist500JSONResponse{}, nil
		}
	}

	wishlistDTO := toWishlistResponse(actor.UserID, []wishlistmodule.WishlistItem{item})
	return openapi.AddToWishlist200JSONResponse(wishlistDTO), nil
}

// RemoveFromWishlist handles DELETE /wishlist/items/{productId}
func (h *Handler) RemoveFromWishlist(ctx context.Context, request openapi.RemoveFromWishlistRequestObject) (openapi.RemoveFromWishlistResponseObject, error) {
	actor := usermodule.Actor{UserID: "usr-1"}
	items, _, err := h.wishlistUC.List(ctx, pagination.Pagination{Limit: 100})
	if err != nil {
		status, errResp := mapWishlistError(err)
		switch status {
		case 401:
			return openapi.RemoveFromWishlist401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.RemoveFromWishlist404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.RemoveFromWishlist500JSONResponse{}, nil
		}
	}

	wishlistDTO := toWishlistResponse(actor.UserID, items)
	return openapi.RemoveFromWishlist200JSONResponse(wishlistDTO), nil
}
