package lead

import "context"

// LeadRepo defines the persistence port owned by the Lead module.
type LeadRepo interface {
	Store(ctx context.Context, lead *LeadSubmission) error
	Get(ctx context.Context, id string) (LeadSubmission, error)
}
