package entity

import "time"

type WorkflowStatus string

const (
	WorkflowStatusNew        WorkflowStatus = "new"
	WorkflowStatusInProgress WorkflowStatus = "in_progress"
	WorkflowStatusDone       WorkflowStatus = "done"
)

type OrderRequest struct {
	ID            string         `json:"id"`
	CartID        string         `json:"cart_id"`
	CustomerName  string         `json:"customer_name"`
	CustomerPhone string         `json:"customer_phone"`
	CustomerEmail string         `json:"customer_email,omitempty"`
	Address       string         `json:"address,omitempty"`
	Note          string         `json:"note,omitempty"`
	Status        WorkflowStatus `json:"status"`
	CartSnapshot  map[string]any `json:"cart_snapshot,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
}
