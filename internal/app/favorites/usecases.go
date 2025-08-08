package favorites

import (
	"context"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
)

type GetClientFavoritesUseCase interface {
	Execute(ctx context.Context, p dto.GetClientFavoritesParams) (dto.ClientFavorites, error)
}

type AddProductToFavoritesUseCase interface {
	Execute(ctx context.Context, p dto.AddProductToFavoritesParams) (dto.ProductFavorite, error)
}

type RemoveProductFromFavoritesUseCase interface {
	Execute(ctx context.Context, p dto.RemoveProductFromFavoritesParams) error
}
