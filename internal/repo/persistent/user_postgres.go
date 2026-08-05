package persistent

import (
	"context"
	"errors"
	"fmt"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/evrone/go-clean-template/internal/entity"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserRepo -.
type UserRepo struct {
	*postgres.Postgres
	mu    sync.RWMutex
	items map[string]usermodule.User
}

// NewUserRepo -.
func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{Postgres: pg, items: make(map[string]usermodule.User)}
}

// Store -.
func (r *UserRepo) Store(ctx context.Context, u *usermodule.User) error {
	if r.Postgres == nil || r.Pool == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		if r.items == nil {
			r.items = make(map[string]usermodule.User)
		}
		r.items[u.ID] = *u
		return nil
	}
	sql, args, err := r.Builder.
		Insert("users").
		Columns("id, username, email, password_hash, created_at, updated_at").
		Values(u.ID, u.Username, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt).
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

	loginSql, loginArgs, err := r.Builder.
		Insert("login_users").
		Columns("id, username, email, password_hash, created_at, updated_at, role_id, role, status, trust_level, is_admin, desktop_notifications_enabled").
		Values(u.ID, u.Username, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt, "00000000-0000-0000-0000-000000000005", "Subscriber", "active", 0, false, false).
		ToSql()
	if err == nil {
		_, _ = r.Pool.Exec(ctx, loginSql, loginArgs...)
	}

	return nil
}

// GetByID -.
func (r *UserRepo) GetByID(ctx context.Context, id string) (usermodule.User, error) {
	return r.getUser(ctx, "id", id)
}

// GetByEmail -.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (usermodule.User, error) {
	return r.getUser(ctx, "email", email)
}

// GetByUsername -.
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (usermodule.User, error) {
	return r.getUser(ctx, "username", username)
}

// List -.
func (r *UserRepo) List(ctx context.Context, filter usermodule.UserFilter) ([]usermodule.User, int, error) {
	if r.Postgres == nil || r.Pool == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		res := make([]usermodule.User, 0, len(r.items))
		for _, item := range r.items {
			res = append(res, item)
		}
		return res, len(res), nil
	}
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

	users := make([]usermodule.User, 0, filter.Limit)
	for rows.Next() {
		var u usermodule.User
		if err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("UserRepo - List - rows.Scan: %w", err)
		}
		users = append(users, u)
	}

	return users, total, nil
}

// Update -.
func (r *UserRepo) Update(ctx context.Context, u *usermodule.User) error {
	if r.Postgres == nil || r.Pool == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.items[u.ID]; !ok {
			return usermodule.ErrNotFound
		}
		r.items[u.ID] = *u
		return nil
	}
	sql, args, err := r.Builder.
		Update("users").
		Set("username", u.Username).
		Set("email", u.Email).
		Set("updated_at", u.UpdatedAt).
		Where(sq.Eq{"id": u.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Update - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - Update - r.Pool.Exec: %w", err)
	}
	if result.RowsAffected() == 0 {
		return usermodule.ErrNotFound
	}

	loginSql, loginArgs, err := r.Builder.
		Update("login_users").
		Set("username", u.Username).
		Set("email", u.Email).
		Set("updated_at", u.UpdatedAt).
		Where(sq.Eq{"id": u.ID}).
		ToSql()
	if err == nil {
		_, _ = r.Pool.Exec(ctx, loginSql, loginArgs...)
	}

	return nil
}

// SetStatus -.
func (r *UserRepo) SetStatus(ctx context.Context, id string, status usermodule.UserStatus) error {
	if r.Postgres == nil || r.Pool == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		item, ok := r.items[id]
		if !ok {
			return usermodule.ErrNotFound
		}
		item.Status = status
		r.items[id] = item
		return nil
	}
	sql, args, err := r.Builder.Update("users").Set("status", status).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - SetStatus - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - SetStatus - r.Pool.Exec: %w", err)
	}
	if result.RowsAffected() == 0 {
		return usermodule.ErrNotFound
	}

	loginSql, loginArgs, err := r.Builder.Update("login_users").Set("status", status).Where(sq.Eq{"id": id}).ToSql()
	if err == nil {
		_, _ = r.Pool.Exec(ctx, loginSql, loginArgs...)
	}

	return nil
}

func (r *UserRepo) getUser(ctx context.Context, column, value string) (usermodule.User, error) {
	if r.Postgres == nil || r.Pool == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		for _, item := range r.items {
			if (column == "id" && item.ID == value) ||
				(column == "email" && item.Email == value) ||
				(column == "username" && item.Username == value) {
				return item, nil
			}
		}
		return usermodule.User{}, usermodule.ErrNotFound
	}
	sql, args, err := r.Builder.
		Select("id, username, email, password_hash, created_at, updated_at").
		From("users").
		Where(sq.Eq{column: value}).
		ToSql()
	if err != nil {
		return usermodule.User{}, fmt.Errorf("UserRepo - getUser - r.Builder: %w", err)
	}

	var u usermodule.User

	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return usermodule.User{}, usermodule.ErrNotFound
		}

		return usermodule.User{}, fmt.Errorf("UserRepo - getUser - r.Pool.QueryRow: %w", err)
	}

	return u, nil
}
