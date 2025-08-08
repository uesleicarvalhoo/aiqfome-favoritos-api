package fixed

import (
	"context"

	"github.com/uesleicarvalhoo/aiqfome/role"
)

type repository struct {
	roles map[role.Role][]role.Permission
}

func NewRepository() role.Repository {
	return &repository{
		roles: map[role.Role][]role.Permission{
			role.RoleAdmin: {
				{
					Resource: role.ResourceClient,
					Action:   role.ActionManage,
				},
				{
					Resource: role.ResourceFavorites,
					Action:   role.ActionManage,
				},
			},
			role.RoleClient: {
				{
					Resource: role.ResourceMe,
					Action:   role.ActionManage,
				},
			},
		},
	}
}

func (r *repository) FindPermissions(_ context.Context, rl role.Role) ([]role.Permission, error) {
	if p, ok := r.roles[rl]; ok {
		return p, nil
	}

	return []role.Permission{}, &role.ErrNotFound{Name: string(rl)}
}
