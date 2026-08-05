package wishlist

import (
	"strconv"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
)

// toWishlistResponse maps a list of wishlistmodule.WishlistItem to openapi.WishlistResponse DTO.
func toWishlistResponse(userID string, items []wishlistmodule.WishlistItem) openapi.WishlistResponse {
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
