package persistent

import (
	"context"
	"testing"

	cartmodule "github.com/evrone/go-clean-template/internal/module/cart"
	"github.com/stretchr/testify/require"
)

func TestCartRepo(t *testing.T) {
	t.Parallel()
	repo := NewCartRepo(nil)
	cart := cartmodule.Cart{ID: "c1", SessionID: "session", Status: cartmodule.CartStatusActive}
	require.NoError(t, repo.Store(context.Background(), &cart))
	require.NoError(t, repo.AddItem(context.Background(), "session", &cartmodule.CartItem{ID: "i1", ProductID: "p1", Quantity: 1}))
	got, err := repo.GetBySession(context.Background(), "session")
	require.NoError(t, err)
	require.Len(t, got.Items, 1)
}
