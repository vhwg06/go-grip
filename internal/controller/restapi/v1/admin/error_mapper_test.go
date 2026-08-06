package admin

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

func TestMapAdminError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "nil error returns 200",
			err:          nil,
			expectedCode: http.StatusOK,
			expectedMsg:  "",
		},
		{
			name:         "ErrOrderNotFound returns 404",
			err:          ordermodule.ErrNotFound,
			expectedCode: http.StatusNotFound,
			expectedMsg:  "Admin resource not found",
		},
		{
			name:         "wrapped ErrOrderNotFound returns 404",
			err:          fmt.Errorf("repository query failed: %w", ordermodule.ErrNotFound),
			expectedCode: http.StatusNotFound,
			expectedMsg:  "Admin resource not found",
		},
		{
			name:         "ErrNotFound returns 404",
			err:          usermodule.ErrNotFound,
			expectedCode: http.StatusNotFound,
			expectedMsg:  "Admin resource not found",
		},
		{
			name:         "usermodule ErrNotFound returns 404",
			err:          usermodule.ErrNotFound,
			expectedCode: http.StatusNotFound,
			expectedMsg:  "Admin resource not found",
		},
		{
			name:         "ordermodule ErrNotFound returns 404",
			err:          ordermodule.ErrNotFound,
			expectedCode: http.StatusNotFound,
			expectedMsg:  "Admin resource not found",
		},
		{
			name:         "ErrUnauthorized returns 401",
			err:          usermodule.ErrUnauthorized,
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Authentication required",
		},
		{
			name:         "usermodule ErrUnauthorized returns 401",
			err:          usermodule.ErrUnauthorized,
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Authentication required",
		},
		{
			name:         "ordermodule ErrUnauthorized returns 401",
			err:          ordermodule.ErrUnauthorized,
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Authentication required",
		},
		{
			name:         "ErrForbidden returns 403",
			err:          usermodule.ErrForbidden,
			expectedCode: http.StatusForbidden,
			expectedMsg:  "Administrative access denied",
		},
		{
			name:         "usermodule ErrForbidden returns 403",
			err:          usermodule.ErrForbidden,
			expectedCode: http.StatusForbidden,
			expectedMsg:  "Administrative access denied",
		},
		{
			name:         "ordermodule ErrForbidden returns 403",
			err:          ordermodule.ErrForbidden,
			expectedCode: http.StatusForbidden,
			expectedMsg:  "Administrative access denied",
		},
		{
			name:         "unmapped generic error returns 500",
			err:          errors.New("unexpected database connection error"),
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "An internal server error occurred",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, resp := mapAdminError(tt.err)
			assert.Equal(t, tt.expectedCode, code)

			if tt.expectedCode == http.StatusOK {
				return
			}

			payload, err := resp.Error.AsErrorPayload()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, payload.Message)
		})
	}
}
