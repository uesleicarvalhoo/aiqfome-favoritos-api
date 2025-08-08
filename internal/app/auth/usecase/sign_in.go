package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/jwt"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/password"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type SignInOptions struct {
	RefreshTokenDuration time.Duration
	AccessTokenDuration  time.Duration
}

type signInUseCase struct {
	repo    user.Repository
	hasher  password.Hasher
	opts    SignInOptions
	access  jwt.Provider
	refresh jwt.Provider
}

func NewSignUseCase(repo user.Repository, hasher password.Hasher, opts SignInOptions, accessProvider, refreshProvider jwt.Provider) auth.SignInUseCase {
	return &signInUseCase{
		repo:    repo,
		hasher:  hasher,
		opts:    opts,
		access:  accessProvider,
		refresh: refreshProvider,
	}
}

func (u *signInUseCase) Execute(ctx context.Context, p dto.SignInParams) (dto.AuthTokens, error) {
	ctx, span := trace.NewSpan(ctx, "auth.signIn")
	defer span.End()

	if err := p.Validate(); err != nil {
		logger.InfoF(ctx, "invalid sign in params", logger.Fields{
			"error": err.Error(),
		})

		return dto.AuthTokens{}, err
	}

	usr, err := u.repo.FindByEmail(ctx, p.Email)
	if err != nil {
		logger.InfoF(ctx, "error while trying to find user", logger.Fields{
			"user_email": p.Email,
			"error":      err.Error(),
		})

		if errors.Is(err, user.ErrNotFound) {
			return dto.AuthTokens{}, domainerror.New(domainerror.ResourceNotFound, "user not found", logger.Fields{
				"user_email": p.Email,
			})
		}

		return dto.AuthTokens{}, domainerror.Wrap(err, domainerror.DependecyError, "error while trying to find user", logger.Fields{
			"user_email": p.Email,
			"error":      err.Error(),
		})
	}

	pwd := fmt.Sprintf("%s:%s", usr.ID.String(), p.Password)
	if err := u.hasher.Compare(usr.PasswordHash, pwd); err != nil {
		return dto.AuthTokens{}, domainerror.New(domainerror.InvalidPassword, "senha invalida", map[string]any{
			"error": err.Error(),
		})
	}

	return generateAuthTokens(ctx, usr.ID.String(), u.access, u.refresh, u.opts.AccessTokenDuration, u.opts.RefreshTokenDuration)
}
