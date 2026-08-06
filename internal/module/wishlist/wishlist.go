package wishlist

import "time"

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

type ReviewModerationStats struct {
	Pending  int `json:"pending"`
	Featured int `json:"featured"`
	Hidden   int `json:"hidden"`
}

// Actor represents an authenticated user context.
type Actor struct {
	UserID   string
	Username string
	IsAdmin  bool
}
