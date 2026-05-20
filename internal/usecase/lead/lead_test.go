package lead

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestLeadUseCaseSubmit(t *testing.T) {
	t.Parallel()
	uc := New(persistent.NewLeadRepo(nil))
	lead, err := uc.Submit(context.Background(), entity.LeadSubmission{Source: "contact", CustomerName: "A", CustomerPhone: "1"})
	require.NoError(t, err)
	require.Equal(t, entity.WorkflowStatusNew, lead.Status)
}
