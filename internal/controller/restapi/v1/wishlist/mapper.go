package wishlist

import (
	"strconv"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// toWishlistResponse maps a list of entity.WishlistItem to openapi.WishlistResponse DTO.
func toWishlistResponse(userID string, items []entity.WishlistItem) openapi.WishlistResponse {
	dtoItems := make([]openapi.WishlistItemResponse, len(items))
	for i, item := range items {
		productID := strconv.FormatInt(item.ID, 10)
		dtoItems[i] = openapi.WishlistItemResponse{
			ProductId: productID,
			CreatedAt: &item.CreatedAt,
		}
	}

	return openapi.WishlistResponse{
		Id:     "wish-" + userID,
		UserId: userID,
		Items:  dtoItems,
	}
}
