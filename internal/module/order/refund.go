package order

import "time"

// RefundStatus represents approval state of a RefundRequest.
type RefundStatus string

const (
	RefundStatusPending  RefundStatus = "pending"
	RefundStatusApproved RefundStatus = "approved"
	RefundStatusRejected RefundStatus = "rejected"
)

// RefundRequest represents a customer request to refund an Order.
type RefundRequest struct {
	ID            int64        `json:"id"`
	OrderID       string       `json:"order_id"`
	UserID        string       `json:"user_id"`
	Username      string       `json:"username"`
	Reason        string       `json:"reason"`
	Status        RefundStatus `json:"status"`
	AdminUsername string       `json:"admin_username,omitempty"`
	AdminNote     string       `json:"admin_note,omitempty"`
	ProductName   string       `json:"product_name,omitempty"`
	Amount        Amount       `json:"amount,omitempty"`
	TradeNo       string       `json:"trade_no,omitempty"`
	OrderStatus   string       `json:"order_status,omitempty"`
	ProcessedAt   *time.Time   `json:"processed_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
