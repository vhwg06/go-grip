package cart

import "time"

// CartStatus represents active state of a shopping cart.
type CartStatus string

const (
	CartStatusActive    CartStatus = "active"
	CartStatusAbandoned CartStatus = "abandoned"
	CartStatusConverted CartStatus = "converted"
)

// CartItem represents an individual line item inside a Cart.
type CartItem struct {
	ID              string         `json:"id"`
	CartID          string         `json:"cart_id"`
	ProductID       string         `json:"product_id"`
	Quantity        int            `json:"quantity"`
	UnitPrice       int64          `json:"unit_price"`
	ProductSnapshot map[string]any `json:"product_snapshot,omitempty"`
	Blocked         bool           `json:"blocked"`
}

// Cart represents a user shopping cart.
type Cart struct {
	ID        string     `json:"id"`
	SessionID string     `json:"session_id"`
	Status    CartStatus `json:"status"`
	Items     []CartItem `json:"items,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
