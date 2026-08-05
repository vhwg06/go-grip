package persistent

import (
	"context"
	"testing"

	mediamodule "github.com/evrone/go-clean-template/internal/module/media"
	"github.com/stretchr/testify/require"
)

func TestMediaRepo(t *testing.T) {
	t.Parallel()
	repo := NewMediaRepo(nil)
	require.NoError(t, repo.Store(context.Background(), &mediamodule.MediaAsset{ID: "m1", FileName: "a.jpg"}))
	asset, err := repo.Get(context.Background(), "m1")
	require.NoError(t, err)
	require.Equal(t, "m1", asset.ID)
}
