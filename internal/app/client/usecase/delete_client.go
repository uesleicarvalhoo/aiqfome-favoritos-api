package usecase

import (
	"context"
	"errors"

	usecase "github.com/uesleicarvalhoo/aiqfome/internal/app/client"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type deleteClientUseCase struct {
	repo user.Repository
}

func NewDeleteClientUseCase(repo user.Repository) usecase.DeleteClientUseCase {
	return &deleteClientUseCase{
		repo: repo,
	}
}

func (u *deleteClientUseCase) Execute(ctx context.Context, p dto.DeleteClientParams) error {
	ctx, span := trace.NewSpan(ctx, "client.deleteClient")
	defer span.End()

	if err := p.Validate(); err != nil {
		logger.ErrorF(ctx, "invalid params", logger.Fields{
			"params": p,
			"error":  err.Error(),
		})

		return err
	}

	usr, err := u.repo.Find(ctx, p.ClientID)
	if err != nil {
		logger.ErrorF(ctx, "error while trying to find client", logger.Fields{
			"error":     err.Error(),
			"client_id": p.ClientID,
		})
		if errors.Is(err, user.ErrNotFound) {
			return domainerror.New(domainerror.ResourceNotFound, "cliente não encontrado", map[string]any{
				"client_id": p.ClientID,
			})
		}

		return domainerror.Wrap(err, domainerror.DependecyError, "erro ao buscar cliente", map[string]any{
			"client_id": p.ClientID,
		})
	}

	if err := u.repo.Delete(ctx, usr); err != nil {
		logger.ErrorF(ctx, "error while trying to paginate clients", logger.Fields{
			"error": err.Error(),
		})
		return domainerror.Wrap(err, domainerror.DependecyError, "erro ao deletar cliente", map[string]any{
			"error":     err.Error(),
			"client_id": usr.ID,
		})
	}

	return nil
}
