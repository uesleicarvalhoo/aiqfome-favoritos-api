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

type listClientsUseCase struct {
	repo user.Repository
}

func NewListClientsUseCase(repo user.Repository) usecase.ListClientsUseCase {
	return &listClientsUseCase{
		repo: repo,
	}
}

func (u *listClientsUseCase) Execute(ctx context.Context, p dto.ListClientsParams) (dto.PaginatedClients, error) {
	ctx, span := trace.NewSpan(ctx, "client.ListClients")
	defer span.End()

	if p.PageSize == 0 {
		p.PageSize = 10
	}

	if p.PageSize == 0 {
		p.PageSize = 10
	}

	if err := p.Validate(); err != nil {
		logger.ErrorF(ctx, "invalid params", logger.Fields{
			"params": p,
			"error":  err.Error(),
		})

		return dto.PaginatedClients{}, err
	}

	uu, total, err := u.repo.Paginate(ctx, p.Page, p.PageSize)
	if err != nil {
		logger.ErrorF(ctx, "error while trying to paginate clients", logger.Fields{
			"error": err.Error(),
		})
		return dto.PaginatedClients{}, domainerror.Wrap(err, domainerror.DependecyError, "erro ao paginar clientes", map[string]any{
			"error":  err.Error(),
			"params": p,
		})
	}

	cc := make([]dto.Client, 0, len(uu))
	for _, u := range uu {
		cc = append(cc, dto.FromDomain(u))
	}

	pages := (total + p.PageSize - 1) / p.PageSize

	return dto.PaginatedClients{
		Clients: cc,
		Total:   total,
		Pages:   pages,
	}, nil
}
