package importer

import "errors"

// Domain errors specific to Importer capability.
var (
	ErrInvalidInput = errors.New("invalid import item payload or threshold exceeded")
)
