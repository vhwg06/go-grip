package catalog

import "errors"

// Domain errors for Catalog capability.
var (
	ErrNotFound     = errors.New("catalog resource not found")
	ErrInvalidInput = errors.New("invalid catalog input")
)
