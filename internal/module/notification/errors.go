package notification

import "errors"

// Domain errors for Notification capability.
var (
	ErrNotFound     = errors.New("notification not found")
	ErrInvalidInput = errors.New("invalid notification payload")
	ErrUnauthorized = errors.New("unauthorized notification operation")
)
