package lead

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// LeadUseCase defines the application service interface for Lead management.
type LeadUseCase interface {
	Submit(ctx context.Context, lead LeadSubmission) (LeadSubmission, error)
	Get(ctx context.Context, id string) (LeadSubmission, error)
}

type leadUseCase struct {
	repo LeadRepo
}

// NewLeadUseCase constructs a new LeadUseCase application service.
func NewLeadUseCase(r LeadRepo) LeadUseCase {
	return &leadUseCase{repo: r}
}

// Submit creates and stores a new LeadSubmission.
func (uc *leadUseCase) Submit(ctx context.Context, lead LeadSubmission) (LeadSubmission, error) {
	lead.ID = uuid.New().String()
	lead.Status = WorkflowStatusNew
	lead.CreatedAt = time.Now().UTC()
	if err := uc.repo.Store(ctx, &lead); err != nil {
		return LeadSubmission{}, err
	}
	return lead, nil
}

// Get retrieves a LeadSubmission by ID.
func (uc *leadUseCase) Get(ctx context.Context, id string) (LeadSubmission, error) {
	return uc.repo.Get(ctx, id)
}
