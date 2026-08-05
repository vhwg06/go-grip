package wishlist

import "errors"

// Domain errors for Wishlist capability.
var (
	ErrNotFound     = errors.New("wishlist item or review not found")
	ErrInvalidInput = errors.New("invalid wishlist input")
	ErrUnauthorized = errors.New("unauthorized wishlist operation")
	ErrForbidden    = errors.New("forbidden wishlist operation")
)
