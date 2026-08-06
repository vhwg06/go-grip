package order

import "time"

// Order represents an e-commerce purchase transaction.
type Order struct {
	ID               string      `json:"id"`
	ProductID        string      `json:"product_id"`
	ProductName      string      `json:"product_name"`
	Amount           Amount      `json:"amount"`
	Quantity         int         `json:"quantity"`
	Email            string      `json:"email"`
	UserID           string      `json:"user_id,omitempty"`
	Username         string      `json:"username,omitempty"`
	Payee            string      `json:"payee,omitempty"`
	Status           OrderStatus `json:"status"`
	StatusText       string      `json:"status_text,omitempty"`
	StatusColor      string      `json:"status_color,omitempty"`
	TradeNo          string      `json:"trade_no,omitempty"`
	CurrentPaymentID string      `json:"current_payment_id,omitempty"`
	PaidAt           *time.Time  `json:"paid_at,omitempty"`
	DeliveredAt      *time.Time  `json:"delivered_at,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

// CanRequestRefund checks if an order is eligible for refund.
func (o Order) CanRequestRefund() bool {
	return o.Status == OrderStatusDelivered
}

// IsTerminal checks if an order is in a final non-modifiable state.
func (o Order) IsTerminal() bool {
	return o.Status == OrderStatusDelivered ||
		o.Status == OrderStatusCancelled ||
		o.Status == OrderStatusFailed ||
		o.Status == OrderStatusRefunded
}

// Actor represents an authenticated user context.
type Actor struct {
	UserID   string
	Username string
	IsAdmin  bool
}
