package ioc

import (
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/usecase"
)

var (
	addProductToFavoritesUc   favorites.AddProductToFavoritesUseCase
	addProductToFavoritesOnce sync.Once
)

func AddProductToFavoritesUseCase() favorites.AddProductToFavoritesUseCase {
	addProductToFavoritesOnce.Do(func() {
		addProductToFavoritesUc = usecase.NewAddProductToFavoritesUseCase(ProductRepository(), FavoriteRepository())
	})

	return addProductToFavoritesUc
}

var (
	getClientFavoritesUc   favorites.GetClientFavoritesUseCase
	getClientFavoritesOnce sync.Once
)

func GetClientFavoritesUseCase() favorites.GetClientFavoritesUseCase {
	getClientFavoritesOnce.Do(func() {
		getClientFavoritesUc = usecase.NewGetClientFavoritesUseCase(FavoriteRepository(), ProductRepository())
	})

	return getClientFavoritesUc
}

var (
	removeProductFromFavoritesUc   favorites.RemoveProductFromFavoritesUseCase
	removeProductFromFavoritesOnce sync.Once
)

func RemoveProductFromFavoritesUseCase() favorites.RemoveProductFromFavoritesUseCase {
	removeProductFromFavoritesOnce.Do(func() {
		removeProductFromFavoritesUc = usecase.NewRemoveProductFromFavoritesUseCase(FavoriteRepository())
	})

	return removeProductFromFavoritesUc
}
