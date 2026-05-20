package entity

import "time"

// User -.
type User struct {
	ID           string     `json:"id"         example:"550e8400-e29b-41d4-a716-446655440000"`
	Username     string     `json:"username"    example:"johndoe"`
	DisplayName  string     `json:"display_name,omitempty" example:"John Doe"`
	Email        string     `json:"email"       example:"john@example.com"`
	PasswordHash string     `json:"-"`
	RoleID       string     `json:"role_id,omitempty"`
	Role         RoleName   `json:"role,omitempty" example:"Administrator"`
	Status       UserStatus `json:"status,omitempty" example:"active"`
	CreatedAt    time.Time  `json:"created_at"  example:"2026-01-01T00:00:00Z"`
	UpdatedAt    time.Time  `json:"updated_at"  example:"2026-01-01T00:00:00Z"`
} // @name entity.User

// UserStatus represents administrative account state.
type UserStatus string

const (
	UserStatusActive UserStatus = "active"
	UserStatusLocked UserStatus = "locked"
)
