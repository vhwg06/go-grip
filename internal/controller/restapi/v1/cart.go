package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type createCartBody struct {
	SessionID string `json:"session_id"`
}

type cartItemBody struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (r *V1) createCart(ctx *fiber.Ctx) error {
	var body createCartBody
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	cart, err := r.cart.Create(ctx.UserContext(), actor.UserID)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusCreated).JSON(cart)
}

func (r *V1) getCart(ctx *fiber.Ctx) error {
	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	cart, err := r.cart.Get(ctx.UserContext(), actor.UserID)
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, "cart not found")
	}
	return ctx.JSON(cart)
}

func (r *V1) addCartItem(ctx *fiber.Ctx) error {
	var body cartItemBody
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	cart, err := r.cart.AddItem(ctx.UserContext(), actor.UserID, body.ProductID, body.Quantity)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(cart)
}

func (r *V1) updateCartItem(ctx *fiber.Ctx) error {
	var body cartItemBody
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	cart, err := r.cart.UpdateItem(ctx.UserContext(), actor.UserID, ctx.Params("item_id"), body.Quantity)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(cart)
}

func (r *V1) removeCartItem(ctx *fiber.Ctx) error {
	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	cart, err := r.cart.RemoveItem(ctx.UserContext(), actor.UserID, ctx.Params("item_id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(cart)
}

func (r *V1) submitOrder(ctx *fiber.Ctx) error {
	actor := r.gripActor(ctx)
	if actor.UserID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	var order entity.OrderRequest
	if err := ctx.BodyParser(&order); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	order.CartID = actor.UserID
	order, err := r.cart.SubmitOrder(ctx.UserContext(), order)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusCreated).JSON(order)
}
