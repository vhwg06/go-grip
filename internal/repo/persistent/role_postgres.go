package persistent

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type RoleRepo struct {
	*postgres.Postgres
}

func NewRoleRepo(pg *postgres.Postgres) *RoleRepo {
	return &RoleRepo{Postgres: pg}
}

func (r *RoleRepo) List(ctx context.Context) ([]entity.Role, error) {
	_ = ctx
	return []entity.Role{
		{ID: "administrator", Name: entity.RoleAdministrator},
		{ID: "editor", Name: entity.RoleEditor},
		{ID: "author", Name: entity.RoleAuthor},
		{ID: "contributor", Name: entity.RoleContributor},
		{ID: "subscriber", Name: entity.RoleSubscriber},
	}, nil
}

func (r *RoleRepo) GetByName(ctx context.Context, name entity.RoleName) (entity.Role, error) {
	roles, _ := r.List(ctx)
	for _, role := range roles {
		if role.Name == name {
			return role, nil
		}
	}
	return entity.Role{}, entity.ErrNotFound
}
