package order

import "errors"

// Domain errors for Order capability.
var (
	ErrNotFound           = errors.New("order or resource not found")
	ErrInvalidInput       = errors.New("invalid order payload")
	ErrUnauthorized       = errors.New("unauthorized order operation")
	ErrForbidden          = errors.New("forbidden order operation")
	ErrRefundNotAllowed   = errors.New("refund not allowed for order")
	ErrPaymentInvalidSign = errors.New("invalid payment signature")
	ErrOutOfStock         = errors.New("product out of stock")
)
