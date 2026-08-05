package persistent

import (
	"context"
	"testing"

	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	"github.com/stretchr/testify/require"
)

func TestSupportChannelRepo(t *testing.T) {
	t.Parallel()
	repo := NewSupportChannelRepo(nil)
	require.NoError(t, repo.Update(context.Background(), &contentmodule.SupportChannel{ID: "s1", ChannelType: "call", IsEnabled: true}))
	items, err := repo.List(context.Background(), true)
	require.NoError(t, err)
	require.Len(t, items, 1)
}
