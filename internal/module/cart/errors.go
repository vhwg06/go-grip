package cart

import "errors"

// Domain errors for Cart capability.
var (
	ErrNotFound     = errors.New("cart not found")
	ErrInvalidInput = errors.New("invalid cart input")
	ErrCartBlocked  = errors.New("cart contains blocked items")
)
