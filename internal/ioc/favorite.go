package ioc

import (
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/favorite"
	"github.com/uesleicarvalhoo/aiqfome/favorite/postgres"
)

var (
	favoriteRepo     favorite.Repository
	favoriteRepoOnce sync.Once
)

func FavoriteRepository() favorite.Repository {
	favoriteRepoOnce.Do(func() {
		favoriteRepo = postgres.NewRepository(Database())
	})

	return favoriteRepo
}
