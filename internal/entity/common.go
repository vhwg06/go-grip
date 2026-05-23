package entity

import "time"

// Actor represents the authenticated principal attached by middleware.
type Actor struct {
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	Email      string `json:"email,omitempty"`
	IsAdmin    bool   `json:"is_admin"`
	IsBlocked  bool   `json:"is_blocked"`
	TrustLevel int    `json:"trust_level"`
}

// Amount stores money in minor units.
type Amount int64

// TimestampRange is used by list filters with created/updated windows.
type TimestampRange struct {
	From time.Time
	To   time.Time
}

type OrderStatus string

const (
	OrderStatusPending       OrderStatus = "pending"
	OrderStatusPaid          OrderStatus = "paid"
	OrderStatusDelivered     OrderStatus = "delivered"
	OrderStatusCancelled     OrderStatus = "cancelled"
	OrderStatusFailed        OrderStatus = "failed"
	OrderStatusRefundPending OrderStatus = "refund_pending"
	OrderStatusRefunded      OrderStatus = "refunded"
)
