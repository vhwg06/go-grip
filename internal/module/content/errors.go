package content

import "errors"

// Domain errors for Content capability.
var (
	ErrNotFound     = errors.New("content or page not found")
	ErrInvalidInput = errors.New("invalid content input")
	ErrUnauthorized = errors.New("unauthorized content operation")
)
