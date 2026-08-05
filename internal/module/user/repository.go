package user

import "context"

// UserFilter defines options for listing users.
type UserFilter struct {
	Limit  uint64
	Offset uint64
}

// UserRepo defines the persistence port owned by User module for User aggregate.
type UserRepo interface {
	Store(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	List(ctx context.Context, filter UserFilter) ([]User, int, error)
	Update(ctx context.Context, user *User) error
	SetStatus(ctx context.Context, userID string, status UserStatus) error
}

// AuthRepo defines the persistence port for refresh sessions and user auth lookup.
type AuthRepo interface {
	GetUserByID(ctx context.Context, userID string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	UpsertUser(ctx context.Context, user User) (User, error)
	StoreRefreshSession(ctx context.Context, session RefreshSession) error
	GetRefreshSession(ctx context.Context, tokenID string) (RefreshSession, error)
	RevokeRefreshSession(ctx context.Context, tokenID string) error
}

// ProfileRepo defines the persistence port for user profile operations.
type ProfileRepo interface {
	GetProfile(ctx context.Context, userID string) (User, error)
	UpdateProfile(ctx context.Context, user User) (User, error)
	GetRecentSessions(ctx context.Context, userID string) ([]RefreshSession, error)
}
