package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"gorm.io/gorm"
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

	if errors.Is(err, entity.ErrUnauthorized) {
		return http.StatusUnauthorized, envelope{Error: "unauthorized"}
	}
	if errors.Is(err, entity.ErrForbidden) {
		return http.StatusForbidden, envelope{Error: "forbidden"}
	}
	if errors.Is(err, entity.ErrNotFound) || errors.Is(err, entity.ErrOrderNotFound) || errors.Is(err, entity.ErrUserNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, envelope{Error: "not_found"}
	}
	if errStr := err.Error(); strings.Contains(errStr, "invalid input syntax for type uuid") || strings.Contains(errStr, "22P02") {
		return http.StatusNotFound, envelope{Error: "not_found"}
	}
	if errors.Is(err, entity.ErrRateLimited) {
		return http.StatusTooManyRequests, envelope{Error: "rate_limited"}
	}
	if errors.Is(err, entity.ErrInvalidInput) || errors.Is(err, entity.ErrInvalidCredentials) {
		return http.StatusBadRequest, envelope{Error: "invalid_input"}
	}
	if errStr := err.Error(); strings.Contains(strings.ToLower(errStr), "duplicate") || strings.Contains(strings.ToLower(errStr), "unique constraint") {
		return http.StatusConflict, envelope{Error: "conflict"}
	}

	return http.StatusInternalServerError, envelope{Error: "internal_error"}
}


