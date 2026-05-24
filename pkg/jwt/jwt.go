package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// ErrUnexpectedSigningMethod is returned when the JWT signing method is not expected.
var ErrUnexpectedSigningMethod = errors.New("unexpected signing method")

// Manager handles JWT token generation and parsing.
type Manager struct {
	secret   string
	duration time.Duration
}

type accessTokenClaims struct {
	IsAdmin  bool   `json:"is_admin"`
	Username string `json:"username,omitempty"`
	jwtlib.RegisteredClaims
}

// New -.
func New(secret string, duration time.Duration) *Manager {
	return &Manager{
		secret:   secret,
		duration: duration,
	}
}

// GenerateToken creates a new JWT token for the given user ID.
func (m *Manager) GenerateToken(userID string) (string, error) {
	return m.GenerateTokenWithProfile(userID, "", false)
}

// GenerateTokenWithProfile creates a JWT token that carries actor profile claims.
func (m *Manager) GenerateTokenWithProfile(userID, username string, isAdmin bool) (string, error) {
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, accessTokenClaims{
		IsAdmin:  isAdmin,
		Username: username,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(m.duration)),
		},
	})

	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("jwt - GenerateTokenWithProfile - token.SignedString: %w", err)
	}

	return tokenString, nil
}

// ParseToken validates a JWT token and returns the user ID.
func (m *Manager) ParseToken(tokenString string) (string, error) {
	userID, _, _, err := m.ParseTokenActor(tokenString)
	return userID, err
}

// ParseTokenActor validates a JWT token and returns subject plus actor flags.
func (m *Manager) ParseTokenActor(tokenString string) (string, bool, string, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &accessTokenClaims{}, func(token *jwtlib.Token) (any, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
		}

		return []byte(m.secret), nil
	})
	if err != nil {
		return "", false, "", fmt.Errorf("jwt - ParseTokenActor - jwtlib.ParseWithClaims: %w", err)
	}

	claims, ok := token.Claims.(*accessTokenClaims)
	if !ok || claims == nil {
		return "", false, "", fmt.Errorf("jwt - ParseTokenActor - invalid claims")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return "", false, "", fmt.Errorf("jwt - ParseTokenActor - GetSubject: %w", err)
	}

	return sub, claims.IsAdmin, claims.Username, nil
}
