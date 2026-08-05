# Wishlist Module Migration Manifest

This document records the inventory of all symbols, contracts, adapters, delivery endpoints, consumers, and migration risks for the `Wishlist` business module.

---

## 1. Owned Symbols
- **Entities**: `WishlistItem`, `WishlistVote`, `Review`, `ReviewStatus`
- **Errors**: `ErrNotFound`, `ErrInvalidInput`, `ErrUnauthorized`, `ErrForbidden`
- **Use Cases**: `WishlistUseCase` (`Wishlist` interface)
- **Repository Ports**: `WishlistRepo`

---

## 2. Ports & Interfaces
```go
package wishlist

import (
	"context"
	"time"

	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

type WishlistItem struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	VoteCount   int       `json:"vote_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WishlistVote struct {
	ID        int64     `json:"id"`
	ItemID    int64     `json:"item_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "PENDING"
	ReviewStatusApproved ReviewStatus = "APPROVED"
	ReviewStatusHidden   ReviewStatus = "HIDDEN"
	ReviewStatusFeatured ReviewStatus = "FEATURED"
)

type Review struct {
	ID                 int64        `json:"id"`
	ProductID          string       `json:"product_id"`
	ProductName        string       `json:"product_name,omitempty"`
	OrderID            string       `json:"order_id"`
	UserID             string       `json:"user_id"`
	Username           string       `json:"username"`
	Rating             int          `json:"rating"`
	Comment            string       `json:"comment"`
	Status             ReviewStatus `json:"status"`
	Attachments        []string     `json:"attachments,omitempty"`
	IsVerifiedPurchase bool         `json:"is_verified_purchase,omitempty"`
	FlaggedReason      *string      `json:"flagged_reason,omitempty"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
}

type ReviewModerationStats struct {
	PendingCount  int `json:"pending_count"`
	ApprovedCount int `json:"approved_count"`
	FlaggedCount  int `json:"flagged_count"`
}

type WishlistRepo interface {
	ListWishlistItems(ctx context.Context, page pagination.Pagination) ([]WishlistItem, int, error)
	StoreWishlistItem(ctx context.Context, item WishlistItem) (WishlistItem, error)
	UpdateWishlistItem(ctx context.Context, item WishlistItem) (WishlistItem, error)
	DeleteWishlistItem(ctx context.Context, itemID int64) error
	ToggleWishlistVote(ctx context.Context, itemID int64, userID string) (bool, error)
	StoreReview(ctx context.Context, review Review) (Review, error)
	ListReviews(ctx context.Context, productID string) ([]Review, error)
}
```

---

## 3. Infrastructure & Delivery Consumers
- `internal/repo/persistent/wishlist_postgres.go`
- `internal/controller/restapi/v1/wishlist/`
- `internal/app/app.go`
