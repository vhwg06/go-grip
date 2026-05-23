package checkout

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/webapi"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/google/uuid"
)

type UseCase struct {
	checkoutRepo repo.CheckoutRepository
	orderRepo    repo.OrderRepository
	verifier     webapi.PaymentVerifier
}

func New(checkoutRepo repo.CheckoutRepository, orderRepo repo.OrderRepository) *UseCase {
	return &UseCase{checkoutRepo: checkoutRepo, orderRepo: orderRepo}
}

func (uc *UseCase) SetPaymentVerifier(verifier webapi.PaymentVerifier) {
	uc.verifier = verifier
}

var _ usecase.Checkout = (*UseCase)(nil)

func (uc *UseCase) Preview(_ context.Context, _ entity.Actor, _ string, quantity int, usePoints bool) (usecase.AmountBreakdown, error) {
	subtotal := entity.Amount(quantity * 10000)
	points := 0
	if usePoints {
		points = quantity * 100
	}
	final := subtotal - entity.Amount(points)
	if final < 0 {
		final = 0
	}

	return usecase.AmountBreakdown{
		Subtotal:    subtotal,
		PointsToUse: points,
		FinalPrice:  final,
	}, nil
}

func (uc *UseCase) CreateOrder(ctx context.Context, actor entity.Actor, productID string, quantity int, email string, usePoints bool) (entity.Order, error) {
	preview, err := uc.Preview(ctx, actor, productID, quantity, usePoints)
	if err != nil {
		return entity.Order{}, err
	}

	order := entity.Order{
		ID:         uuid.NewString(),
		ProductID:  productID,
		Amount:     preview.FinalPrice,
		Quantity:   quantity,
		Email:      email,
		UserID:     actor.UserID,
		Username:   actor.Username,
		Status:     entity.OrderStatusPending,
		PointsUsed: preview.PointsToUse,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	created, err := uc.checkoutRepo.CreateOrderWithReservation(ctx, actor, order)
	if err != nil {
		return entity.Order{}, fmt.Errorf("CheckoutUseCase - CreateOrder - checkoutRepo.CreateOrderWithReservation: %w", err)
	}

	if preview.FinalPrice == 0 {
		if err := uc.checkoutRepo.UpdateOrderStatus(ctx, created.ID, entity.OrderStatusDelivered); err != nil {
			return entity.Order{}, fmt.Errorf("CheckoutUseCase - CreateOrder - checkoutRepo.UpdateOrderStatus: %w", err)
		}
		created.Status = entity.OrderStatusDelivered
	}

	return created, nil
}

func (uc *UseCase) PaymentParams(_ context.Context, _ entity.Actor, orderID string) (usecase.PaymentParams, error) {
	return usecase.PaymentParams{
		URL: "https://payment.example/checkout",
		Fields: map[string]string{
			"order_id": orderID,
		},
	}, nil
}

func (uc *UseCase) PaymentNotify(ctx context.Context, payload map[string]string) error {
	orderID := payload["order_id"]
	if orderID == "" {
		return entity.ErrInvalidInput
	}

	if uc.verifier != nil {
		ok, err := uc.verifier.Verify(ctx, payload)
		if err != nil {
			return fmt.Errorf("CheckoutUseCase - PaymentNotify - verifier.Verify: %w", err)
		}
		if !ok {
			return entity.ErrPaymentInvalidSign
		}
	}

	order, err := uc.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("CheckoutUseCase - PaymentNotify - orderRepo.GetOrderByID: %w", err)
	}
	if order.Status == entity.OrderStatusDelivered {
		return nil
	}

	payment := entity.Payment{
		ID:                uuid.NewString(),
		OrderID:           orderID,
		Provider:          "epay",
		ProviderPaymentID: payload["payment_id"],
		Amount:            order.Amount,
		Status:            "success",
		IsSignatureValid:  true,
		ProcessedAt:       ptrTime(time.Now().UTC()),
		CreatedAt:         time.Now().UTC(),
	}
	if err := uc.checkoutRepo.AttachPayment(ctx, payment); err != nil {
		return fmt.Errorf("CheckoutUseCase - PaymentNotify - checkoutRepo.AttachPayment: %w", err)
	}

	if err := uc.checkoutRepo.UpdateOrderStatus(ctx, orderID, entity.OrderStatusPaid); err != nil {
		return fmt.Errorf("CheckoutUseCase - PaymentNotify - checkoutRepo.UpdateOrderStatus(paid): %w", err)
	}
	if err := uc.checkoutRepo.UpdateOrderStatus(ctx, orderID, entity.OrderStatusDelivered); err != nil {
		return fmt.Errorf("CheckoutUseCase - PaymentNotify - checkoutRepo.UpdateOrderStatus(delivered): %w", err)
	}
	return nil
}

func (uc *UseCase) PaymentStatus(ctx context.Context, orderID string) (entity.Order, error) {
	order, err := uc.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return entity.Order{}, fmt.Errorf("CheckoutUseCase - PaymentStatus - orderRepo.GetOrderByID: %w", err)
	}
	return order, nil
}

func (uc *UseCase) Cancel(ctx context.Context, actor entity.Actor, orderID string) error {
	if err := uc.orderRepo.CancelPendingOrder(ctx, actor, orderID); err != nil {
		return fmt.Errorf("CheckoutUseCase - Cancel - orderRepo.CancelPendingOrder: %w", err)
	}
	if err := uc.checkoutRepo.ReleaseReservation(ctx, orderID); err != nil {
		return fmt.Errorf("CheckoutUseCase - Cancel - checkoutRepo.ReleaseReservation: %w", err)
	}
	return nil
}

func ptrTime(v time.Time) *time.Time {
	return &v
}
