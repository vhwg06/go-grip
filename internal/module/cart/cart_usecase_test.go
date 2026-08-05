package cart_test

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/module/cart"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

type dummyNotifier struct{}

func (d *dummyNotifier) Dispatch(_ context.Context, _, _, _ string) error {
	return nil
}

func TestCartUseCaseOrderRequest(t *testing.T) {
	t.Parallel()
	cartRepo := persistent.NewCartRepo(nil)
	uc := cart.NewCartUseCase(cartRepo, persistent.NewOrderRequestRepo(nil), &dummyNotifier{})
	_, err := uc.AddItem(context.Background(), "session", "product", 1)
	require.NoError(t, err)
	order, err := uc.SubmitOrder(context.Background(), cart.OrderRequest{CartID: "session", CustomerName: "A", CustomerPhone: "1"})
	require.NoError(t, err)
	require.NotEmpty(t, order.ID)
}
