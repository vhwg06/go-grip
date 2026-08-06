package catalog

import (
	"time"

	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// ProductStatus represents publication status of a catalog item.
// Actor represents an authenticated user context.
type Actor struct {
	UserID   string
	Username string
	IsAdmin  bool
}

type ProductStatus string

const (
	ProductStatusActive    ProductStatus = "active"
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusArchived  ProductStatus = "archived"
	ProductStatusSuspended ProductStatus = "suspended"
)

// ProductSpecItem represents a key-value attribute pair for a product.
type ProductSpecItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Product represents a catalog merchandise entity.
type Product struct {
	ID                      string            `json:"id"`
	Title                   string            `json:"title"`
	SKU                     string            `json:"sku"`
	Description             string            `json:"description"`
	Price                   int64             `json:"price"`
	ComparePrice            *int64            `json:"compare_price,omitempty"`
	Status                  ProductStatus     `json:"status"`
	Brand                   string            `json:"brand,omitempty"`
	CategoryID              string            `json:"category_id"`
	ImageURL                string            `json:"image_url,omitempty"`
	Images                  []string          `json:"images"`
	IsHot                   bool              `json:"is_hot"`
	IsActive                bool              `json:"is_active"`
	SortOrder               int               `json:"sort_order"`
	PurchaseLimit           int               `json:"purchase_limit"`
	PurchaseWarning         string            `json:"purchase_warning,omitempty"`
	VisibilityLevel         int               `json:"visibility_level"`
	StockCount              int               `json:"stock_count"`
	LockedCount             int               `json:"locked_count"`
	SoldCount               int               `json:"sold_count"`
	Rating                  float64           `json:"rating"`
	ReviewCount             int               `json:"review_count"`
	Attributes              map[string]any    `json:"attributes,omitempty"`
	MediaIDs                []string          `json:"media_ids,omitempty"`
	IntroArticleID          string            `json:"intro_article_id,omitempty"`
	MaxPurchaseableQuantity int               `json:"max_purchaseable_quantity,omitempty"`
	Specs                   []ProductSpecItem `json:"specs,omitempty"`
	CreatedAt               time.Time         `json:"created_at"`
	UpdatedAt               time.Time         `json:"updated_at"`
}

// Setting represents a system key-value setting.
type Setting struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductFilter represents search/filtering criteria for product listings.
type ProductFilter struct {
	Keyword    string
	CategoryID string
	Brand      string
	MinPrice   *int64
	MaxPrice   *int64
	Sort       string
	Pagination pagination.Pagination
}
