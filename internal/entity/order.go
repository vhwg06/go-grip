package entity

import "time"

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
	PointsUsed       int         `json:"points_used"`
	CurrentPaymentID string      `json:"current_payment_id,omitempty"`
	PaidAt           *time.Time  `json:"paid_at,omitempty"`
	DeliveredAt      *time.Time  `json:"delivered_at,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

type Payment struct {
	ID                     string     `json:"id"`
	OrderID                string     `json:"order_id"`
	Provider               string     `json:"provider"`
	ProviderPaymentID      string     `json:"provider_payment_id,omitempty"`
	Amount                 Amount     `json:"amount"`
	Status                 string     `json:"status"`
	RequestPayloadSummary  string     `json:"request_payload_summary,omitempty"`
	CallbackPayloadSummary string     `json:"callback_payload_summary,omitempty"`
	IsSignatureValid       bool       `json:"is_signature_valid"`
	ProcessedAt            *time.Time `json:"processed_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
}

type RefundStatus string

const (
	RefundStatusPending  RefundStatus = "pending"
	RefundStatusApproved RefundStatus = "approved"
	RefundStatusRejected RefundStatus = "rejected"
)

type RefundRequest struct {
	ID            int64        `json:"id"`
	OrderID       string       `json:"order_id"`
	UserID        string       `json:"user_id"`
	Username      string       `json:"username"`
	Reason        string       `json:"reason"`
	Status        RefundStatus `json:"status"`
	AdminUsername string       `json:"admin_username,omitempty"`
	AdminNote     string       `json:"admin_note,omitempty"`
	ProcessedAt   *time.Time   `json:"processed_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (o Order) CanRequestRefund() bool {
	return o.Status == OrderStatusDelivered
}

func (o Order) IsTerminal() bool {
	return o.Status == OrderStatusDelivered ||
		o.Status == OrderStatusCancelled ||
		o.Status == OrderStatusFailed ||
		o.Status == OrderStatusRefunded
}
