package media

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestMediaUseCaseValidation(t *testing.T) {
	t.Parallel()
	uc := New(persistent.NewMediaRepo(nil), entity.MaxMediaUploadBytes)
	asset, err := uc.Store(context.Background(), entity.MediaAsset{FileName: "a.jpg", MimeType: "image/jpeg", SizeBytes: 10})
	require.NoError(t, err)
	require.NotEmpty(t, asset.ID)
	_, err = uc.Store(context.Background(), entity.MediaAsset{FileName: "a.gif", MimeType: "image/gif", SizeBytes: 10})
	require.ErrorIs(t, err, entity.ErrInvalidInput)
}
