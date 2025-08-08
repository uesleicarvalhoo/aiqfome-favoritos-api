package usecase

import (
	"context"

	"github.com/uesleicarvalhoo/aiqfome/favorite"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
)

type removeProductFromFavoritesUseCase struct {
	repo favorite.Repository
}

func NewRemoveProductFromFavoritesUseCase(repo favorite.Repository) favorites.RemoveProductFromFavoritesUseCase {
	return &removeProductFromFavoritesUseCase{
		repo: repo,
	}
}

func (u removeProductFromFavoritesUseCase) Execute(ctx context.Context, p dto.RemoveProductFromFavoritesParams) error {
	ctx, span := trace.NewSpan(ctx, "favorites.removeProductFromFavorites")
	defer span.End()

	if err := p.Validate(); err != nil {
		logger.ErrorF(ctx, "invalid params", logger.Fields{
			"params": p,
			"error":  err.Error(),
		})
		return err
	}

	f, err := u.repo.Find(ctx, p.ClientID, p.ProductID)
	if err != nil {
		if nfErr, ok := err.(*favorite.ErrFavoriteNotFound); ok {
			return domainerror.New(domainerror.ResourceNotFound, "favorito não encontrado", map[string]any{
				"client_id":  nfErr.ClientID,
				"product_id": nfErr.ProductID,
			})
		}

		return domainerror.Wrap(err, domainerror.DependecyError, "error while to trying find favorite", map[string]any{
			"client_id":  p.ClientID,
			"product_id": p.ProductID,
			"error":      err,
		})
	}

	if err := u.repo.Remove(ctx, f); err != nil {
		return domainerror.Wrap(err, domainerror.DependecyError, "error while trying to remove favorite", map[string]any{
			"client_id":  p.ClientID,
			"product_id": p.ProductID,
			"error":      err,
		})
	}

	return nil
}
