package cart

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/evrone/go-clean-template/internal/usecase/notification"
	"github.com/stretchr/testify/require"
)

func TestCartUseCaseOrderRequest(t *testing.T) {
	t.Parallel()
	cartRepo := persistent.NewCartRepo(nil)
	uc := New(cartRepo, persistent.NewOrderRequestRepo(nil), notification.New(false))
	_, err := uc.AddItem(context.Background(), "session", "product", 1)
	require.NoError(t, err)
	order, err := uc.SubmitOrder(context.Background(), entity.OrderRequest{CartID: "session", CustomerName: "A", CustomerPhone: "1"})
	require.NoError(t, err)
	require.NotEmpty(t, order.ID)
}
