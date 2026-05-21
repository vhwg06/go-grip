package lead

import (
	"context"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/google/uuid"
)

type UseCase struct {
	repo repo.LeadRepo
}

func New(r repo.LeadRepo) *UseCase { return &UseCase{repo: r} }

func (uc *UseCase) Submit(ctx context.Context, lead entity.LeadSubmission) (entity.LeadSubmission, error) {
	lead.ID = uuid.New().String()
	lead.Status = entity.WorkflowStatusNew
	lead.CreatedAt = time.Now().UTC()
	if err := uc.repo.Store(ctx, &lead); err != nil {
		return entity.LeadSubmission{}, err
	}
	return lead, nil
}

func (uc *UseCase) Get(ctx context.Context, id string) (entity.LeadSubmission, error) {
	return uc.repo.Get(ctx, id)
}
