package auth

import (
	"context"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type SignInUseCase interface {
	Execute(ctx context.Context, params dto.SignInParams) (dto.AuthTokens, error)
}

type SignUpUseCase interface {
	Execute(ctx context.Context, params dto.SignUpParams) (user.User, error)
}

type RefreshTokenUseCase interface {
	Execute(ctx context.Context, params dto.RefreshTokenParams) (dto.AuthTokens, error)
}

type AuthenticateUseCase interface {
	Execute(ctx context.Context, token string) (user.User, error)
}

type AuthorizeUseCase interface {
	Execute(ctx context.Context, params dto.AuthorizeParams) error
}
