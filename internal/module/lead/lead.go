package lead

import "time"

// WorkflowStatus represents the lifecycle state of a lead.
type WorkflowStatus string

const (
	WorkflowStatusNew        WorkflowStatus = "new"
	WorkflowStatusInProgress WorkflowStatus = "in_progress"
	WorkflowStatusDone       WorkflowStatus = "done"
)

// LeadSubmission represents a customer sales lead contact submission.
type LeadSubmission struct {
	ID            string         `json:"id"`
	Source        string         `json:"source"`
	CustomerName  string         `json:"customer_name"`
	CustomerPhone string         `json:"customer_phone"`
	CustomerEmail string         `json:"customer_email,omitempty"`
	Message       string         `json:"message,omitempty"`
	Status        WorkflowStatus `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
}
