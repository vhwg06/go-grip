package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type gripCheckoutPreviewResponse struct {
	ProductID  string        `json:"productId"`
	Quantity   int           `json:"quantity"`
	Subtotal   entity.Amount `json:"subtotal"`
	FinalPrice entity.Amount `json:"finalPrice"`
}

type gripCreateOrderRequest struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
	Email     string `json:"email"`
}

type gripRefundRequest struct {
	Reason string `json:"reason"`
}

func (r *V1) gripCheckoutUC() (usecase.Checkout, bool) {
	return r.checkout, r.checkout != nil
}

// @Summary     Preview checkout
// @Description Calculates subtotal and final payable amount
// @ID          grip_checkout_preview
// @Tags        checkout
// @Produce     json
// @Param       product_id query string true "Product ID"
// @Param       quantity query int true "Quantity"
// @Success     200 {object} envelope
// @Failure     400 {object} envelope
// @Failure     401 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /checkout/preview [get]
func (r *V1) gripCheckoutPreview(ctx *fiber.Ctx) error {
	uc, ok := r.gripCheckoutUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "checkout_usecase_not_configured"})
	}

	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		status, body := mapDomainError(entity.ErrUnauthorized)
		return ctx.Status(status).JSON(body)
	}

	productID := ctx.Query("product_id")
	quantity := ctx.QueryInt("quantity", 1)
	if productID == "" || quantity <= 0 {
		status, body := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(body)
	}

	breakdown, err := uc.Preview(ctx.UserContext(), actor, productID, quantity)
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	return ctx.JSON(apiSuccessEnvelope(gripCheckoutPreviewResponse{
		ProductID:  productID,
		Quantity:   quantity,
		Subtotal:   breakdown.Subtotal,
		FinalPrice: breakdown.FinalPrice,
	}))
}

// @Summary     Create order
// @Description Creates an order and reserves stock for checkout
// @ID          grip_checkout_create_order
// @Tags        checkout
// @Accept      json
// @Produce     json
// @Param       request body gripCreateOrderRequest true "Order payload"
// @Success     201 {object} envelope
// @Failure     400 {object} envelope
// @Failure     401 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /checkout/orders [post]
func (r *V1) gripCreateOrder(ctx *fiber.Ctx) error {
	uc, ok := r.gripCheckoutUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "checkout_usecase_not_configured"})
	}

	var body gripCreateOrderRequest
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}
	if body.ProductID == "" || body.Quantity <= 0 {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		status, payload := mapDomainError(entity.ErrUnauthorized)
		return ctx.Status(status).JSON(payload)
	}

	order, err := uc.CreateOrder(ctx.UserContext(), actor, body.ProductID, body.Quantity, body.Email)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	order = gripDecorateOrderStatus(order)
	return ctx.Status(http.StatusCreated).JSON(apiSuccessEnvelope(order))
}

// @Summary     Create payment order
// @Description Creates a checkout order for direct payment flow
// @ID          grip_checkout_create_payment_order
// @Tags        checkout
// @Accept      json
// @Produce     json
// @Param       request body gripCreateOrderRequest true "Order payload"
// @Success     201 {object} envelope
// @Failure     400 {object} envelope
// @Failure     401 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /checkout/payment-orders [post]
func (r *V1) gripCreatePaymentOrder(ctx *fiber.Ctx) error {
	return r.gripCreateOrder(ctx)
}

// @Summary     Get payment params
// @Description Recreates payment gateway params for an order
// @ID          grip_checkout_payment_params
// @Tags        checkout
// @Produce     json
// @Param       id path string true "Order ID"
// @Success     200 {object} envelope
// @Failure     401 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /checkout/orders/{id}/payment-params [get]
func (r *V1) gripPaymentParams(ctx *fiber.Ctx) error {
	uc, ok := r.gripCheckoutUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "checkout_usecase_not_configured"})
	}

	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		status, body := mapDomainError(entity.ErrUnauthorized)
		return ctx.Status(status).JSON(body)
	}

	params, err := uc.PaymentParams(ctx.UserContext(), actor, ctx.Params("id"))
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}
	return ctx.JSON(apiSuccessEnvelope(params))
}

// @Summary     Get order payment status
// @Description Polls payment and order lifecycle state
// @ID          grip_checkout_order_status
// @Tags        checkout
// @Produce     json
// @Param       id path string true "Order ID"
// @Success     200 {object} envelope
// @Failure     401 {object} envelope
// @Failure     404 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /checkout/orders/{id}/status [get]
func (r *V1) gripOrderStatus(ctx *fiber.Ctx) error {
	uc, ok := r.gripCheckoutUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "checkout_usecase_not_configured"})
	}

	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		status, body := mapDomainError(entity.ErrUnauthorized)
		return ctx.Status(status).JSON(body)
	}

	order, err := uc.PaymentStatus(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	return ctx.JSON(apiSuccessEnvelope(gripDecorateOrderStatus(order)))
}

// @Summary     Cancel pending order
// @Description Cancels a pending order and releases reserved stock
// @ID          grip_checkout_cancel_order
// @Tags        checkout
// @Produce     json
// @Param       id path string true "Order ID"
// @Success     204 {string} string ""
// @Failure     401 {object} envelope
// @Failure     400 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /checkout/orders/{id}/cancel [post]
func (r *V1) gripCancelOrder(ctx *fiber.Ctx) error {
	uc, ok := r.gripCheckoutUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "checkout_usecase_not_configured"})
	}

	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		status, body := mapDomainError(entity.ErrUnauthorized)
		return ctx.Status(status).JSON(body)
	}

	if err := uc.Cancel(ctx.UserContext(), actor, ctx.Params("id")); err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Payment notify callback
// @Description Processes payment provider notifications
// @ID          grip_checkout_notify
// @Tags        checkout
// @Accept      json
// @Produce     json
// @Success     200 {object} envelope
// @Failure     400 {object} envelope
// @Failure     500 {object} envelope
// @Router      /checkout/notify [post]
func (r *V1) gripPaymentNotify(ctx *fiber.Ctx) error {
	uc, ok := r.gripCheckoutUC()
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "checkout_usecase_not_configured"})
	}

	payload := make(map[string]string)
	if err := ctx.BodyParser(&payload); err != nil {
		status, body := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(body)
	}

	if err := uc.PaymentNotify(ctx.UserContext(), payload); err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}

	return ctx.JSON(apiSuccessEnvelope(fiber.Map{"status": "ok"}))
}

// @Summary     Payment callback redirect
// @Description Handles payment provider redirect callback for storefront
// @ID          grip_checkout_callback
// @Tags        checkout
// @Produce     json
// @Param       id path string true "Order ID"
// @Success     200 {object} envelope
// @Router      /checkout/callback/{id} [get]
func (r *V1) gripPaymentCallback(ctx *fiber.Ctx) error {
	return ctx.JSON(apiSuccessEnvelope(fiber.Map{
		"orderId": ctx.Params("id"),
		"status":  "received",
	}))
}
