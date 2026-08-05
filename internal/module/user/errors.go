package user

import "errors"

// Domain errors for User & Identity capability.
var (
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("action forbidden")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrNotFound           = errors.New("user not found")
	ErrInvalidInput       = errors.New("invalid input data")
)
