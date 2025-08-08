package usecase

import (
	"context"

	usecase "github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
	"github.com/uesleicarvalhoo/aiqfome/role"
)

type authorizeUseCase struct {
	repo role.Repository
}

func NewAuthorizeUseCase(repo role.Repository) usecase.AuthorizeUseCase {
	return &authorizeUseCase{
		repo: repo,
	}
}

func (u *authorizeUseCase) Execute(ctx context.Context, params dto.AuthorizeParams) error {
	ctx, span := trace.NewSpan(ctx, "auth.authorize")
	defer span.End()

	// TODO: Implementar cache
	pp, err := u.repo.FindPermissions(ctx, params.User.Role)
	if err != nil {
		if _, ok := err.(*role.ErrNotFound); ok {
			logger.ErrorF(ctx, "role not found", logger.Fields{
				"role": params.User.Role,
			})

			return domainerror.New(domainerror.ResourceNotFound, "role não encontrada", map[string]any{"role": params.User.Role})
		}

		return domainerror.Wrap(err, domainerror.DependecyError, "ocorreu um erro ao obter as permissões", map[string]any{"role": params.User.Role, "error": err.Error()})
	}

	if params.User.Role == role.RoleAdmin {
		logger.InfoF(ctx, "access grant by role admin", logger.Fields{
			"client_id":   params.User.ID,
			"client_role": params.User.Role,
		})

		return nil
	}

	for _, p := range pp {
		if p.Resource != params.Resource {
			continue
		}

		switch {
		case p.Action == role.ActionManage:
			logger.InfoF(ctx, "authorization granted by manage", logger.Fields{
				"client_id": params.User.ID,
				"role":      params.User.Role,
				"resource":  params.Resource,
				"action":    params.Action,
			})

			return nil

		case p.Action == params.Action:
			logger.InfoF(ctx, "authorization granded", logger.Fields{
				"client_id": params.User.ID,
				"role":      params.User.Role,
				"resource":  params.Resource,
				"action":    params.Action,
			})

			return nil

		case p.Action == role.ActionWrite && params.Action == role.ActionRead:
			logger.InfoF(ctx, "authorization granded by write > read", logger.Fields{
				"client_id": params.User.ID,
				"role":      params.User.Role,
				"resource":  params.Resource,
				"action":    params.Action,
			})
			return nil
		}
	}

	logger.ErrorF(ctx, "permission denied", logger.Fields{
		"client_id":   params.User.ID,
		"role":        params.User.Role,
		"resource":    params.Resource,
		"action":      params.Action,
		"permissions": pp,
	})

	return domainerror.New(
		domainerror.OperationNotAllowed,
		"permissão negada",
		map[string]any{
			"client_id":   params.User.ID,
			"role":        params.User.Role,
			"resource":    params.Resource,
			"action":      params.Action,
			"permissions": pp,
		},
	)
}
