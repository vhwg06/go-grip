package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPublicProductModelListResponse(t *testing.T) {
	t.Parallel()

	t.Run("empty items map returns empty slice and zero total", func(t *testing.T) {
		t.Parallel()
		resMap := map[string]any{
			"items": []map[string]any{},
			"total": 0,
		}
		resp := toPublicProductModelListResponse(resMap)
		assert.NotNil(t, resp.Items)
		assert.Equal(t, 0, len(*resp.Items))
		assert.NotNil(t, resp.Total)
		assert.Equal(t, 0, *resp.Total)
	})

	t.Run("populated items map returns mapped items and total", func(t *testing.T) {
		t.Parallel()
		rawItems := []map[string]any{
			{"id": "model-1", "name": "Model 1"},
			{"id": "model-2", "name": "Model 2"},
		}
		resMap := map[string]any{
			"items": rawItems,
			"total": 2,
		}
		resp := toPublicProductModelListResponse(resMap)
		assert.NotNil(t, resp.Items)
		assert.Equal(t, 2, len(*resp.Items))
		assert.NotNil(t, resp.Total)
		assert.Equal(t, 2, *resp.Total)
		assert.Equal(t, "model-1", (*resp.Items)[0]["id"])
	})
}
