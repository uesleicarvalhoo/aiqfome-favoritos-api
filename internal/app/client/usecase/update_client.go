package usecase

import (
	"context"

	usecase "github.com/uesleicarvalhoo/aiqfome/internal/app/client"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type updateClientUseCase struct {
	repo user.Repository
}

func NewUpdateClientUseCase(repo user.Repository) usecase.UpdateClientUseCase {
	return &updateClientUseCase{
		repo: repo,
	}
}

func (u *updateClientUseCase) Execute(ctx context.Context, p dto.UpdateClientParams) (dto.Client, error) {
	ctx, span := trace.NewSpan(ctx, "client.updateClient")
	defer span.End()

	if err := p.Validate(); err != nil {
		logger.ErrorF(ctx, "invalid params", logger.Fields{
			"params": p,
			"error":  err.Error(),
		})

		return dto.Client{}, err
	}

	usr, err := u.repo.Find(ctx, p.ClientID)
	if err != nil {
		logger.ErrorF(ctx, "error while trying to find client", logger.Fields{
			"error":     err.Error(),
			"client_id": p.ClientID,
		})
	}

	if err := u.update(&usr, p); err != nil {
		logger.ErrorF(ctx, "validation failed after update client", logger.Fields{
			"params": p,
			"error":  err.Error(),
		})
		return dto.Client{}, err
	}

	if err := u.repo.Update(ctx, usr); err != nil {
		return dto.Client{}, domainerror.Wrap(err, domainerror.DependecyError, "failed to update client", map[string]any{
			"error": err.Error(),
		})
	}

	return dto.NewFromDomain(usr), nil
}

func (u *updateClientUseCase) update(usr *user.User, p dto.UpdateClientParams) error {
	if p.Name != nil {
		usr.Name = *p.Name
	}

	if p.Active != nil {
		usr.Active = *p.Active
	}

	if p.Role != nil {
		usr.Role = *p.Role
	}

	return usr.Validate()
}
