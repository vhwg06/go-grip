package entity

import "time"

type Card struct {
	ID              int64      `json:"id"`
	ProductID       string     `json:"product_id"`
	CardKey         string     `json:"card_key"`
	IsUsed          bool       `json:"is_used"`
	ReservedOrderID string     `json:"reserved_order_id"`
	ReservedAt      *time.Time `json:"reserved_at,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	UsedAt          *time.Time `json:"used_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}
