package entity

import "time"

// User -.
type User struct {
	ID                          string     `json:"id"         example:"550e8400-e29b-41d4-a716-446655440000"`
	Username                    string     `json:"username"    example:"johndoe"`
	DisplayName                 string     `json:"display_name,omitempty" example:"John Doe"`
	Email                       string     `json:"email"       example:"john@example.com"`
	PasswordHash                string     `json:"-"`
	RoleID                      string     `json:"role_id,omitempty"`
	Role                        RoleName   `json:"role,omitempty" example:"Administrator"`
	Status                      UserStatus `json:"status,omitempty" example:"active"`
	Provider                    string     `json:"provider,omitempty"`
	ProviderID                  string     `json:"provider_id,omitempty"`
	Points                      int        `json:"points"`
	TrustLevel                  int        `json:"trust_level"`
	IsAdmin                     bool       `json:"is_admin"`
	DesktopNotificationsEnabled bool       `json:"desktop_notifications_enabled"`
	LastLoginAt                 *time.Time `json:"last_login_at,omitempty"`
	LastCheckinAt               *time.Time `json:"last_checkin_at,omitempty"`
	ConsecutiveDays             int        `json:"consecutive_days"`
	CreatedAt                   time.Time  `json:"created_at"  example:"2026-01-01T00:00:00Z"`
	UpdatedAt                   time.Time  `json:"updated_at"  example:"2026-01-01T00:00:00Z"`
} // @name entity.User

// UserStatus represents administrative account state.
type UserStatus string

const (
	UserStatusActive UserStatus = "active"
	UserStatusLocked UserStatus = "locked"
)

type RefreshSession struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenID   string    `json:"token_id"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at"`
	CreatedAt time.Time `json:"created_at"`
}

type DailyCheckin struct {
	ID           int64     `json:"id"`
	UserID       string    `json:"user_id"`
	CheckinDate  time.Time `json:"checkin_date"`
	RewardAmount int       `json:"reward_amount"`
	StreakAfter  int       `json:"streak_after"`
	CreatedAt    time.Time `json:"created_at"`
}
