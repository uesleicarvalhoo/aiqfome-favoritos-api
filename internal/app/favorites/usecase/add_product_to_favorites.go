package usecase

import (
	"context"

	"github.com/uesleicarvalhoo/aiqfome/favorite"
	usecase "github.com/uesleicarvalhoo/aiqfome/internal/app/favorites"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
	"github.com/uesleicarvalhoo/aiqfome/product"
)

type addProductToFavoritesUseCase struct {
	products  product.Reader
	favorites favorite.Repository
}

func NewAddProductToFavoritesUseCase(productReader product.Reader, favoriteRepo favorite.Repository) usecase.AddProductToFavoritesUseCase {
	return &addProductToFavoritesUseCase{
		products:  productReader,
		favorites: favoriteRepo,
	}
}

func (u *addProductToFavoritesUseCase) Execute(ctx context.Context, p dto.AddProductToFavoritesParams) (dto.ProductFavorite, error) {
	ctx, span := trace.NewSpan(ctx, "favorites.addProductToFavorite")
	defer span.End()

	if err := p.Validate(); err != nil {
		logger.ErrorF(ctx, "invalid params", logger.Fields{
			"error":  err.Error(),
			"params": p,
		})

		return dto.ProductFavorite{}, err
	}

	pd, err := u.products.Find(ctx, p.ProductID)
	if err != nil {
		logger.ErrorF(ctx, "error while trying to find product", logger.Fields{
			"product_id": p.ProductID,
			"error":      err.Error(),
		})

		if _, ok := err.(*product.ErrNotFound); ok {
			return dto.ProductFavorite{}, domainerror.Wrap(err, domainerror.ResourceNotFound, "produto não encontrado", map[string]any{
				"product_id": p.ProductID,
			})
		}

		return dto.ProductFavorite{}, domainerror.Wrap(err, domainerror.DependecyError, "erro ao obter dados do produto", map[string]any{
			"product_id": p.ProductID,
			"error":      err.Error(),
		})
	}

	if _, err := u.favorites.Find(ctx, p.ClientID, p.ProductID); err == nil {
		return dto.ProductFavorite{}, domainerror.New(domainerror.ProductAlreadyIsFavorite, "o produto já está nos favoritos", map[string]any{
			"client_id":  p.ClientID,
			"product_id": p.ProductID,
		})
	}

	f, err := favorite.New(p.ClientID, p.ProductID)
	if err != nil {
		logger.ErrorF(ctx, "invalid favorite params", logger.Fields{
			"client_id":  p.ClientID,
			"product_id": p.ProductID,
			"error":      err.Error(),
		})

		return dto.ProductFavorite{}, err
	}

	if err := u.favorites.Create(ctx, f); err != nil {
		return dto.ProductFavorite{}, domainerror.Wrap(err, domainerror.DependecyError, "erro ao adicionar o produto aos favoritos", map[string]any{
			"client_id":  p.ClientID,
			"product_id": p.ProductID,
			"error":      err.Error(),
		})
	}

	return dto.ProductFavorite{
		ClientID: f.ClientID,
		Product:  pd,
	}, nil
}
