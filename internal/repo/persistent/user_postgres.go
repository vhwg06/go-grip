package persistent

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserRepo -.
type UserRepo struct {
	*postgres.Postgres
}

// NewUserRepo -.
func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

// Store -.
func (r *UserRepo) Store(ctx context.Context, user *entity.User) error {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("id, username, email, password_hash, created_at, updated_at").
		Values(user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.ErrUserAlreadyExists
		}

		return fmt.Errorf("UserRepo - Store - r.Pool.Exec: %w", err)
	}

	// Also insert into login_users to keep them in sync
	loginSql, loginArgs, err := r.Builder.
		Insert("login_users").
		Columns("id, username, email, password_hash, created_at, updated_at, role_id, role, status, trust_level, is_admin, desktop_notifications_enabled").
		Values(user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt, "00000000-0000-0000-0000-000000000005", "Subscriber", "active", 0, false, false).
		ToSql()
	if err == nil {
		_, _ = r.Pool.Exec(ctx, loginSql, loginArgs...)
	}

	return nil
}

// GetByID -.
func (r *UserRepo) GetByID(ctx context.Context, id string) (entity.User, error) {
	return r.getUser(ctx, "id", id)
}

// GetByEmail -.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	return r.getUser(ctx, "email", email)
}

// List -.
func (r *UserRepo) List(ctx context.Context, filter repo.UserFilter) ([]entity.User, int, error) {
	countSQL, countArgs, err := r.Builder.Select("COUNT(*)").From("users").ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("UserRepo - List - count builder: %w", err)
	}

	var total int
	if err = r.Pool.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("UserRepo - List - count query: %w", err)
	}

	sql, args, err := r.Builder.
		Select("id, username, email, password_hash, created_at, updated_at").
		From("users").
		OrderBy("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("UserRepo - List - data builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("UserRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	users := make([]entity.User, 0, filter.Limit)
	for rows.Next() {
		var user entity.User
		if err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("UserRepo - List - rows.Scan: %w", err)
		}
		users = append(users, user)
	}

	return users, total, nil
}

// Update -.
func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	sql, args, err := r.Builder.
		Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Update - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - Update - r.Pool.Exec: %w", err)
	}
	if result.RowsAffected() == 0 {
		return entity.ErrUserNotFound
	}

	// Also update login_users
	loginSql, loginArgs, err := r.Builder.
		Update("login_users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.ID}).
		ToSql()
	if err == nil {
		_, _ = r.Pool.Exec(ctx, loginSql, loginArgs...)
	}

	return nil
}

// SetStatus -.
func (r *UserRepo) SetStatus(ctx context.Context, id string, status entity.UserStatus) error {
	// Existing databases may not yet have a status column until feature migrations run.
	sql, args, err := r.Builder.Update("users").Set("status", status).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - SetStatus - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - SetStatus - r.Pool.Exec: %w", err)
	}
	if result.RowsAffected() == 0 {
		return entity.ErrUserNotFound
	}

	// Also update login_users status
	loginSql, loginArgs, err := r.Builder.Update("login_users").Set("status", status).Where(sq.Eq{"id": id}).ToSql()
	if err == nil {
		_, _ = r.Pool.Exec(ctx, loginSql, loginArgs...)
	}

	return nil
}

func (r *UserRepo) getUser(ctx context.Context, column, value string) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, username, email, password_hash, created_at, updated_at").
		From("users").
		Where(sq.Eq{column: value}).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - getUser - r.Builder: %w", err)
	}

	var user entity.User

	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, fmt.Errorf("UserRepo - getUser - r.Pool.QueryRow: %w", err)
	}

	return user, nil
}
