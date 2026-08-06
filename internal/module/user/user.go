package user

import "time"

// UserStatus represents administrative account state.
type UserStatus string

const (
	UserStatusActive UserStatus = "active"
	UserStatusLocked UserStatus = "locked"
)

// User represents an account holder.
type User struct {
	ID                          string     `json:"id"`
	Username                    string     `json:"username"`
	DisplayName                 string     `json:"display_name,omitempty"`
	Email                       string     `json:"email"`
	PasswordHash                string     `json:"-"`
	RoleID                      string     `json:"role_id,omitempty"`
	Role                        RoleName   `json:"role,omitempty"`
	Status                      UserStatus `json:"status,omitempty"`
	Provider                    string     `json:"provider,omitempty"`
	ProviderID                  string     `json:"provider_id,omitempty"`
	TrustLevel                  int        `json:"trust_level"`
	IsAdmin                     bool       `json:"is_admin"`
	DesktopNotificationsEnabled bool       `json:"desktop_notifications_enabled"`
	LastLoginAt                 *time.Time `json:"last_login_at,omitempty"`
	CustomerID                  *string    `json:"customerId,omitempty"`
	OrderCount                  *int       `json:"orderCount,omitempty"`
	RefundCount                 *int       `json:"refundCount,omitempty"`
	ReviewCount                 *int       `json:"reviewCount,omitempty"`
	IsBlocked                   bool       `json:"is_blocked"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
}

// RefreshSession represents a stored refresh token session.
type RefreshSession struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenID   string    `json:"token_id"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Actor represents an authenticated user context passed across application boundaries.
type Actor struct {
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	Email      string `json:"email,omitempty"`
	IsAdmin    bool   `json:"is_admin"`
	IsBlocked  bool   `json:"is_blocked"`
	TrustLevel int    `json:"trust_level"`
}
