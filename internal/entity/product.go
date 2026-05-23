package entity

import "time"

type ProductStatus string

const (
	ProductStatusActive    ProductStatus = "active"
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusArchived  ProductStatus = "archived"
	ProductStatusSuspended ProductStatus = "suspended"
)

type Product struct {
	ID                      string         `json:"id"`
	Title                   string         `json:"title"`
	SKU                     string         `json:"sku"`
	Description             string         `json:"description"`
	Price                   int64          `json:"price"`
	ComparePrice            *int64         `json:"compare_price,omitempty"`
	Status                  ProductStatus  `json:"status"`
	Brand                   string         `json:"brand,omitempty"`
	CategoryID              string         `json:"category_id,omitempty"`
	ImageURL                string         `json:"image_url,omitempty"`
	IsHot                   bool           `json:"is_hot"`
	IsActive                bool           `json:"is_active"`
	IsShared                bool           `json:"is_shared"`
	SortOrder               int            `json:"sort_order"`
	PurchaseLimit           int            `json:"purchase_limit"`
	PurchaseWarning         string         `json:"purchase_warning,omitempty"`
	VisibilityLevel         int            `json:"visibility_level"`
	StockCount              int            `json:"stock_count"`
	LockedCount             int            `json:"locked_count"`
	SoldCount               int            `json:"sold_count"`
	Rating                  float64        `json:"rating"`
	ReviewCount             int            `json:"review_count"`
	Attributes              map[string]any `json:"attributes,omitempty"`
	MediaIDs                []string       `json:"media_ids,omitempty"`
	MaxPurchaseableQuantity int            `json:"max_purchaseable_quantity,omitempty"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
}

type Card struct {
	ID              int64      `json:"id"`
	ProductID       string     `json:"product_id"`
	CardKey         string     `json:"card_key"`
	IsUsed          bool       `json:"is_used"`
	ReservedOrderID string     `json:"reserved_order_id,omitempty"`
	ReservedAt      *time.Time `json:"reserved_at,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	UsedAt          *time.Time `json:"used_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type Setting struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductFilter struct {
	Keyword    string
	Brand      string
	MinPrice   *int64
	MaxPrice   *int64
	Sort       string
	Pagination Pagination
}
