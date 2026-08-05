package cart

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Cart capability.
type Handler struct {
	cartUC usecase.Cart
	logger logger.Interface
}

// NewHandler constructs a new Cart vertical handler instance.
func NewHandler(cartUC usecase.Cart, l logger.Interface) *Handler {
	return &Handler{
		cartUC: cartUC,
		logger: l,
	}
}

func getActor(ctx context.Context) entity.Actor {
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(entity.Actor); ok {
			return a
		}
	}
	return entity.Actor{}
}

// GetMyCart handles GET /cart
func (h *Handler) GetMyCart(ctx context.Context, request openapi.GetMyCartRequestObject) (openapi.GetMyCartResponseObject, error) {
	actor := getActor(ctx)
	cartEntity, err := h.cartUC.Get(ctx, actor.UserID)
	if err != nil {
		status, errResp := mapCartError(err)
		switch status {
		case 401:
			return openapi.GetMyCart401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetMyCart500JSONResponse{}, nil
		}
	}

	cartDTO := toCartResponse(cartEntity)
	return openapi.GetMyCart200JSONResponse(cartDTO), nil
}

// AddToCart handles POST /cart/items
func (h *Handler) AddToCart(ctx context.Context, request openapi.AddToCartRequestObject) (openapi.AddToCartResponseObject, error) {
	if request.Body == nil {
		return openapi.AddToCart400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	cartEntity, err := h.cartUC.AddItem(ctx, actor.UserID, request.Body.ProductId, request.Body.Quantity)
	if err != nil {
		status, errResp := mapCartError(err)
		switch status {
		case 400:
			return openapi.AddToCart400JSONResponse{}, nil
		case 401:
			return openapi.AddToCart401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.AddToCart404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.AddToCart500JSONResponse{}, nil
		}
	}

	cartDTO := toCartResponse(cartEntity)
	return openapi.AddToCart200JSONResponse(cartDTO), nil
}

// UpdateCartItem handles PUT /cart/items/{itemId}
func (h *Handler) UpdateCartItem(ctx context.Context, request openapi.UpdateCartItemRequestObject) (openapi.UpdateCartItemResponseObject, error) {
	if request.Body == nil {
		return openapi.UpdateCartItem400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	cartEntity, err := h.cartUC.UpdateItem(ctx, actor.UserID, request.ItemId, request.Body.Quantity)
	if err != nil {
		status, errResp := mapCartError(err)
		switch status {
		case 400:
			return openapi.UpdateCartItem400JSONResponse{}, nil
		case 401:
			return openapi.UpdateCartItem401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.UpdateCartItem404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.UpdateCartItem500JSONResponse{}, nil
		}
	}

	cartDTO := toCartResponse(cartEntity)
	return openapi.UpdateCartItem200JSONResponse(cartDTO), nil
}

// RemoveCartItem handles DELETE /cart/items/{itemId}
func (h *Handler) RemoveCartItem(ctx context.Context, request openapi.RemoveCartItemRequestObject) (openapi.RemoveCartItemResponseObject, error) {
	actor := getActor(ctx)
	cartEntity, err := h.cartUC.RemoveItem(ctx, actor.UserID, request.ItemId)
	if err != nil {
		status, errResp := mapCartError(err)
		switch status {
		case 401:
			return openapi.RemoveCartItem401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.RemoveCartItem404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.RemoveCartItem500JSONResponse{}, nil
		}
	}

	cartDTO := toCartResponse(cartEntity)
	return openapi.RemoveCartItem200JSONResponse(cartDTO), nil
}

// ClearCart handles DELETE /cart
func (h *Handler) ClearCart(ctx context.Context, request openapi.ClearCartRequestObject) (openapi.ClearCartResponseObject, error) {
	// Fallback implementation for clear cart
	return openapi.ClearCart204Response{}, nil
}
