package persistent

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/module/lead"
	"github.com/stretchr/testify/require"
)

func TestLeadRepo(t *testing.T) {
	t.Parallel()
	repo := NewLeadRepo(nil)
	require.NoError(t, repo.Store(context.Background(), &lead.LeadSubmission{ID: "l1"}))
	l, err := repo.Get(context.Background(), "l1")
	require.NoError(t, err)
	require.Equal(t, "l1", l.ID)
}
