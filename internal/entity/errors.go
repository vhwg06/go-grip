package entity

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskForbidden      = errors.New("task does not belong to user")
	ErrInvalidTransition  = errors.New("invalid status transition")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrDuplicateSKU       = errors.New("duplicate sku")
	ErrProductUnavailable = errors.New("product unavailable")
	ErrCartBlocked        = errors.New("cart blocked")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrPaymentFailed      = errors.New("payment failed")
	ErrPaymentInvalidSign = errors.New("payment signature invalid")
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderStateConflict = errors.New("order state conflict")
	ErrOutOfStock         = errors.New("out of stock")
	ErrPointsInsufficient = errors.New("insufficient points")
	ErrRefundNotAllowed   = errors.New("refund not allowed")
	ErrRateLimited        = errors.New("rate limited")
)
