package usecase

import (
	"context"
	"errors"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/jwt"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type authenticateUseCase struct {
	repo   user.Repository
	access jwt.Provider
}

func NewAuthenticateUseCase(repo user.Repository, access jwt.Provider) auth.AuthenticateUseCase {
	return &authenticateUseCase{
		repo:   repo,
		access: access,
	}
}

func (u *authenticateUseCase) Execute(ctx context.Context, token string) (user.User, error) {
	ctx, span := trace.NewSpan(ctx, "auth.authenticate")
	defer span.End()

	cl, err := u.access.Validate(ctx, token)
	if err != nil {
		logger.InfoF(ctx, "invalid access token", logger.Fields{
			"error": err.Error(),
		})

		return user.User{}, err
	}

	// TODO: Implementar cache
	c, err := u.repo.Find(ctx, cl.UserID)
	if err != nil {
		logger.InfoF(ctx, "error while trying to find user", logger.Fields{
			"client_id": cl.UserID,
			"error":     err.Error(),
		})

		if errors.Is(err, user.ErrNotFound) {
			return user.User{}, domainerror.Wrap(err, domainerror.ResourceNotFound, "user not found", logger.Fields{
				"client_id": cl.UserID,
			})
		}

		return user.User{}, domainerror.Wrap(err, domainerror.DependecyError, "error while trying to find user", logger.Fields{
			"client_id": cl.UserID,
			"error":     err.Error(),
		})
	}

	if !c.Active {
		logger.WarnF(ctx, "usuário bloqueado", logger.Fields{
			"client_id": c.ID,
		})
		return user.User{}, domainerror.New(domainerror.UserNotActive, "usuário bloqueado", map[string]any{
			"client_id": c.ID,
		})
	}

	return c, nil
}
