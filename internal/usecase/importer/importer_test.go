package importer

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestImporterPartialFailures(t *testing.T) {
	t.Parallel()
	catalogRepo := persistent.NewCatalogRepo(nil)
	contentRepo := persistent.NewContentRepo(nil)
	uc := New(persistent.NewImportRepo(nil, catalogRepo, contentRepo), 25)
	result, err := uc.Import(context.Background(), []entity.ImportItem{
		{Type: entity.ImportItemProduct, Data: map[string]any{"title": "Phone", "sku": "sku"}},
		{Type: "unknown", Data: map[string]any{}},
	})
	require.NoError(t, err)
	require.Equal(t, 1, result.Imported)
	require.Len(t, result.Failed, 1)
}
