package persistent

import (
	"context"
	"testing"

	cartmodule "github.com/evrone/go-clean-template/internal/module/cart"
	"github.com/stretchr/testify/require"
)

func TestOrderRequestRepo(t *testing.T) {
	t.Parallel()
	repo := NewOrderRequestRepo(nil)
	require.NoError(t, repo.Store(context.Background(), &cartmodule.OrderRequest{ID: "o1"}))
}
