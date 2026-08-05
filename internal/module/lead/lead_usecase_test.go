package lead_test

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/module/lead"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestLeadUseCaseSubmit(t *testing.T) {
	t.Parallel()
	uc := lead.NewLeadUseCase(persistent.NewLeadRepo(nil))
	ls, err := uc.Submit(context.Background(), lead.LeadSubmission{
		Source:        "contact",
		CustomerName:  "A",
		CustomerPhone: "1",
	})
	require.NoError(t, err)
	require.Equal(t, lead.WorkflowStatusNew, ls.Status)
}
