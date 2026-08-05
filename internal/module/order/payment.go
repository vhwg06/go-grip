package order

import "time"

// Payment represents payment transaction details associated with an Order.
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
