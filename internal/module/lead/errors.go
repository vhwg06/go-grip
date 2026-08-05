package lead

import "errors"

// Domain errors specific to Lead capability.
var (
	ErrNotFound = errors.New("lead submission not found")
)
