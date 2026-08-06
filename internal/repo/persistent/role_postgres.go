package persistent

import (
	"context"

	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type RoleRepo struct {
	*postgres.Postgres
}

func NewRoleRepo(pg *postgres.Postgres) *RoleRepo {
	return &RoleRepo{Postgres: pg}
}

func (r *RoleRepo) List(ctx context.Context) ([]usermodule.Role, error) {
	_ = ctx
	return []usermodule.Role{
		{ID: "administrator", Name: usermodule.RoleAdministrator},
		{ID: "editor", Name: usermodule.RoleEditor},
		{ID: "author", Name: usermodule.RoleAuthor},
		{ID: "contributor", Name: usermodule.RoleContributor},
		{ID: "subscriber", Name: usermodule.RoleSubscriber},
	}, nil
}

func (r *RoleRepo) GetByName(ctx context.Context, name usermodule.RoleName) (usermodule.Role, error) {
	roles, _ := r.List(ctx)
	for _, role := range roles {
		if role.Name == name {
			return role, nil
		}
	}
	return usermodule.Role{}, usermodule.ErrNotFound
}
