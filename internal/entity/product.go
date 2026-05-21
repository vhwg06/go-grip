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
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	SKU          string         `json:"sku"`
	Description  string         `json:"description"`
	Price        int64          `json:"price"`
	ComparePrice *int64         `json:"compare_price,omitempty"`
	Status       ProductStatus  `json:"status"`
	Brand        string         `json:"brand,omitempty"`
	Attributes   map[string]any `json:"attributes,omitempty"`
	MediaIDs     []string       `json:"media_ids,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type ProductFilter struct {
	Keyword    string
	Brand      string
	MinPrice   *int64
	MaxPrice   *int64
	Sort       string
	Pagination Pagination
}
