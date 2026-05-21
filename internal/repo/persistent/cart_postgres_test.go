package persistent

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestCartRepo(t *testing.T) {
	t.Parallel()
	repo := NewCartRepo(nil)
	cart := entity.Cart{ID: "c1", SessionID: "session", Status: entity.CartStatusActive}
	require.NoError(t, repo.Store(context.Background(), &cart))
	require.NoError(t, repo.AddItem(context.Background(), "session", &entity.CartItem{ID: "i1", ProductID: "p1", Quantity: 1}))
	got, err := repo.GetBySession(context.Background(), "session")
	require.NoError(t, err)
	require.Len(t, got.Items, 1)
}
