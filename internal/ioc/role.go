package ioc

import (
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/role"
	"github.com/uesleicarvalhoo/aiqfome/role/fixed"
)

var (
	roleRepo     role.Repository
	roleRepoOnce sync.Once
)

func RoleRepository() role.Repository {
	roleRepoOnce.Do(func() {
		roleRepo = fixed.NewRepository()
	})

	return roleRepo
}
