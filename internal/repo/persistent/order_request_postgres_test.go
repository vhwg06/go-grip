package persistent

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestOrderRequestRepo(t *testing.T) {
	t.Parallel()
	repo := NewOrderRequestRepo(nil)
	require.NoError(t, repo.Store(context.Background(), &entity.OrderRequest{ID: "o1"}))
}
