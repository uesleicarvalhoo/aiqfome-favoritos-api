package ioc

import (
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/user"
	"github.com/uesleicarvalhoo/aiqfome/user/postgres"
)

var (
	userRepo     user.Repository
	userRepoOnce sync.Once
)

func UserRepository() user.Repository {
	userRepoOnce.Do(func() {
		userRepo = postgres.NewRepository(Database())
	})

	return userRepo
}
