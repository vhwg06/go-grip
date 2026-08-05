package media

import "errors"

// Domain errors specific to Media capability.
var (
	ErrNotFound     = errors.New("media asset not found")
	ErrInvalidInput = errors.New("invalid media asset data")
)
