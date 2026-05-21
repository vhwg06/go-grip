package entity

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskForbidden      = errors.New("task does not belong to user")
	ErrInvalidTransition  = errors.New("invalid status transition")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrDuplicateSKU       = errors.New("duplicate sku")
	ErrProductUnavailable = errors.New("product unavailable")
	ErrCartBlocked        = errors.New("cart blocked")
)
