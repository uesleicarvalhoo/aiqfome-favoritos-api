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

type getClientFavoritesUseCase struct {
	favorites favorite.Repository
	products  product.Repository
}

func NewGetClientFavoritesUseCase(favoritesRepo favorite.Repository, productsRepo product.Repository) usecase.GetClientFavoritesUseCase {
	return &getClientFavoritesUseCase{
		favorites: favoritesRepo,
		products:  productsRepo,
	}
}

func (u *getClientFavoritesUseCase) Execute(ctx context.Context, p dto.GetClientFavoritesParams) (dto.ClientFavorites, error) {
	ctx, span := trace.NewSpan(ctx, "favorites.getClientFavorites")
	defer span.End()

	if p.PageSize == 0 {
		p.PageSize = 10
	}

	if err := p.Validate(); err != nil {
		logger.ErrorF(ctx, "invalid params", logger.Fields{
			"params": p,
			"error":  err.Error(),
		})

		return dto.ClientFavorites{}, err
	}

	fvs, total, err := u.favorites.PaginateByClientID(ctx, p.ClientID, p.Page, p.PageSize)
	if err != nil {
		logger.ErrorF(ctx, "error while trying to paginate favorites", logger.Fields{
			"error": err.Error(),
		})
		return dto.ClientFavorites{}, domainerror.Wrap(err, domainerror.DependecyError, "erro ao paginar favoritos", map[string]any{
			"error":  err.Error(),
			"params": p,
		})
	}

	pIds := make([]int, 0, len(fvs))

	for _, f := range fvs {
		pIds = append(pIds, f.ProductID)
	}

	pds, err := u.getProducts(ctx, pIds)
	if err != nil {
		return dto.ClientFavorites{}, err
	}

	pages := (total + p.PageSize - 1) / p.PageSize

	return dto.ClientFavorites{
		ClientID: p.ClientID,
		Products: pds,
		Total:    total,
		Pages:    pages,
	}, nil
}

func (u *getClientFavoritesUseCase) getProducts(ctx context.Context, ids []int) ([]product.Product, error) {
	pp, err := u.products.FindMultiple(ctx, ids)
	if err != nil {
		if nfErr, ok := err.(*product.ErrProductsNotFound); ok {
			logger.ErrorF(ctx, "products not found", logger.Fields{
				"products_not_found": nfErr.IDs,
			})

			return []product.Product{}, domainerror.Wrap(err, domainerror.ResourceNotFound, "produtos não encontrados", map[string]any{
				"products_not_found": nfErr.IDs,
			})
		}

		logger.ErrorF(ctx, "error while trying to get products", logger.Fields{
			"error":       err.Error(),
			"product_ids": ids,
		})

		return []product.Product{}, domainerror.Wrap(err, domainerror.DependecyError, "erro ao buscar produtos", map[string]any{
			"error":       err.Error(),
			"product_ids": ids,
		})
	}

	return pp, nil
}
