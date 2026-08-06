package wishlist

import (
	"context"

	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// WishlistRepo defines persistence operations for Wishlist items and reviews.
type WishlistRepo interface {
	ListWishlistItems(ctx context.Context, page pagination.Pagination) ([]WishlistItem, int, error)
	StoreWishlistItem(ctx context.Context, item WishlistItem) (WishlistItem, error)
	UpdateWishlistItem(ctx context.Context, item WishlistItem) (WishlistItem, error)
	DeleteWishlistItem(ctx context.Context, itemID int64) error
	ToggleWishlistVote(ctx context.Context, itemID int64, userID string) (bool, error)
	StoreReview(ctx context.Context, review Review) (Review, error)
	ListReviews(ctx context.Context, productID string) ([]Review, error)
	GetReview(ctx context.Context, reviewID int64) (Review, error)
	DeleteReview(ctx context.Context, reviewID int64) error
}

// OrderReader defines read operations consumed by WishlistUseCase for review verification.
type OrderReader interface {
	GetOrderByID(ctx context.Context, orderID string) (ordermodule.Order, error)
}
