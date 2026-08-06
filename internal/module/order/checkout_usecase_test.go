package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type checkoutRepoStub struct {
	created int
}

func (s *checkoutRepoStub) CreateOrderWithReservation(_ context.Context, _ Actor, order Order) (Order, error) {
	s.created++
	return order, nil
}

func (s *checkoutRepoStub) AttachPayment(context.Context, Payment) error { return nil }

func (s *checkoutRepoStub) UpdateOrderStatus(context.Context, string, OrderStatus) error { return nil }

func (s *checkoutRepoStub) ReleaseReservation(context.Context, string) error { return nil }

func TestCheckoutCreateOrderRejectsInvalidPublicInputBeforePersistence(t *testing.T) {
	ctx := context.Background()
	repo := &checkoutRepoStub{}
	uc := NewCheckoutUseCase(repo, nil)
	actor := Actor{UserID: "user-1", Username: "buyer"}

	for _, test := range []struct {
		name      string
		productID string
		quantity  int
		email     string
	}{
		{name: "invalid product id", productID: "not-a-uuid", quantity: 1, email: "buyer@example.com"},
		{name: "invalid quantity", productID: "11111111-1111-1111-1111-111111111111", quantity: 0, email: "buyer@example.com"},
		{name: "missing email", productID: "11111111-1111-1111-1111-111111111111", quantity: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := uc.CreateOrder(ctx, actor, test.productID, test.quantity, test.email)
			require.ErrorIs(t, err, ErrInvalidInput)
			require.Zero(t, repo.created)
		})
	}
}
