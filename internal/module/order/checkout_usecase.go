package order

import (
	"context"
	"fmt"
	"time"

	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/google/uuid"
)

// AmountBreakdown represents price details for checkout preview.
type AmountBreakdown struct {
	Subtotal   Amount
	FinalPrice Amount
}

// PaymentParams represents parameters for initiating payment.
type PaymentParams struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

// PaymentVerifier defines contract for verifying payment callbacks.
type PaymentVerifier interface {
	Verify(ctx context.Context, payload map[string]string) (bool, error)
}

// CheckoutUseCase defines application services for order checkout and payment processing.
type CheckoutUseCase interface {
	Preview(ctx context.Context, actor usermodule.Actor, productID string, quantity int) (AmountBreakdown, error)
	CreateOrder(ctx context.Context, actor usermodule.Actor, productID string, quantity int, email string) (Order, error)
	PaymentParams(ctx context.Context, actor usermodule.Actor, orderID string) (PaymentParams, error)
	PaymentNotify(ctx context.Context, payload map[string]string) error
	PaymentStatus(ctx context.Context, orderID string) (Order, error)
	Cancel(ctx context.Context, actor usermodule.Actor, orderID string) error
	SetPaymentVerifier(verifier PaymentVerifier)
}

type checkoutUseCase struct {
	checkoutRepo CheckoutRepo
	orderRepo    OrderRepo
	verifier     PaymentVerifier
}

// NewCheckoutUseCase constructs a new CheckoutUseCase instance.
func NewCheckoutUseCase(checkoutRepo CheckoutRepo, orderRepo OrderRepo) CheckoutUseCase {
	return &checkoutUseCase{checkoutRepo: checkoutRepo, orderRepo: orderRepo}
}

func (uc *checkoutUseCase) SetPaymentVerifier(verifier PaymentVerifier) {
	uc.verifier = verifier
}

func (uc *checkoutUseCase) Preview(_ context.Context, _ usermodule.Actor, _ string, quantity int) (AmountBreakdown, error) {
	subtotal := Amount(quantity * 10000)
	return AmountBreakdown{
		Subtotal:   subtotal,
		FinalPrice: subtotal,
	}, nil
}

func (uc *checkoutUseCase) CreateOrder(ctx context.Context, actor usermodule.Actor, productID string, quantity int, email string) (Order, error) {
	preview, err := uc.Preview(ctx, actor, productID, quantity)
	if err != nil {
		return Order{}, err
	}

	o := Order{
		ID:        uuid.NewString(),
		ProductID: productID,
		Amount:    preview.FinalPrice,
		Quantity:  quantity,
		Email:     email,
		UserID:    actor.UserID,
		Username:  actor.Username,
		Status:    OrderStatusPending,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	created, err := uc.checkoutRepo.CreateOrderWithReservation(ctx, actor, o)
	if err != nil {
		return Order{}, fmt.Errorf("CheckoutUseCase - CreateOrder - checkoutRepo.CreateOrderWithReservation: %w", err)
	}

	return created, nil
}

func (uc *checkoutUseCase) PaymentParams(_ context.Context, _ usermodule.Actor, orderID string) (PaymentParams, error) {
	return PaymentParams{
		URL: "https://payment.example/checkout",
		Fields: map[string]string{
			"order_id": orderID,
		},
	}, nil
}

func (uc *checkoutUseCase) PaymentNotify(ctx context.Context, payload map[string]string) error {
	orderID := payload["order_id"]
	if orderID == "" {
		return ErrInvalidInput
	}

	if uc.verifier != nil {
		ok, err := uc.verifier.Verify(ctx, payload)
		if err != nil {
			return fmt.Errorf("CheckoutUseCase - PaymentNotify - verifier.Verify: %w", err)
		}
		if !ok {
			return ErrPaymentInvalidSign
		}
	}

	o, err := uc.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("CheckoutUseCase - PaymentNotify - orderRepo.GetOrderByID: %w", err)
	}
	if o.Status == OrderStatusDelivered {
		return nil
	}

	now := time.Now().UTC()
	payment := Payment{
		ID:                uuid.NewString(),
		OrderID:           orderID,
		Provider:          "epay",
		ProviderPaymentID: payload["payment_id"],
		Amount:            o.Amount,
		Status:            "success",
		IsSignatureValid:  true,
		ProcessedAt:       &now,
		CreatedAt:         now,
	}
	if err := uc.checkoutRepo.AttachPayment(ctx, payment); err != nil {
		return fmt.Errorf("CheckoutUseCase - PaymentNotify - checkoutRepo.AttachPayment: %w", err)
	}

	if err := uc.checkoutRepo.UpdateOrderStatus(ctx, orderID, OrderStatusPaid); err != nil {
		return fmt.Errorf("CheckoutUseCase - PaymentNotify - checkoutRepo.UpdateOrderStatus(paid): %w", err)
	}
	if err := uc.checkoutRepo.UpdateOrderStatus(ctx, orderID, OrderStatusDelivered); err != nil {
		return fmt.Errorf("CheckoutUseCase - PaymentNotify - checkoutRepo.UpdateOrderStatus(delivered): %w", err)
	}
	return nil
}

func (uc *checkoutUseCase) PaymentStatus(ctx context.Context, orderID string) (Order, error) {
	o, err := uc.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return Order{}, fmt.Errorf("CheckoutUseCase - PaymentStatus - orderRepo.GetOrderByID: %w", err)
	}
	return o, nil
}

func (uc *checkoutUseCase) Cancel(ctx context.Context, actor usermodule.Actor, orderID string) error {
	if err := uc.orderRepo.CancelPendingOrder(ctx, actor, orderID); err != nil {
		return fmt.Errorf("CheckoutUseCase - Cancel - orderRepo.CancelPendingOrder: %w", err)
	}
	if err := uc.checkoutRepo.ReleaseReservation(ctx, orderID); err != nil {
		return fmt.Errorf("CheckoutUseCase - Cancel - checkoutRepo.ReleaseReservation: %w", err)
	}
	return nil
}
