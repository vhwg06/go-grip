package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
)

type envelope struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func apiSuccessEnvelope(data any) envelope {
	return envelope{Data: data}
}

func mapDomainError(err error) (int, envelope) {
	if err == nil {
		return http.StatusOK, envelope{}
	}

	switch err {
	case entity.ErrUnauthorized:
		return http.StatusUnauthorized, envelope{Error: "unauthorized"}
	case entity.ErrForbidden:
		return http.StatusForbidden, envelope{Error: "forbidden"}
	case entity.ErrNotFound, entity.ErrOrderNotFound, entity.ErrUserNotFound:
		return http.StatusNotFound, envelope{Error: "not_found"}
	case entity.ErrRateLimited:
		return http.StatusTooManyRequests, envelope{Error: "rate_limited"}
	case entity.ErrInvalidInput, entity.ErrInvalidCredentials:
		return http.StatusBadRequest, envelope{Error: "invalid_input"}
	default:
		return http.StatusInternalServerError, envelope{Error: "internal_error"}
	}
}
