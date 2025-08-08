package usecase

import (
	"context"
	"time"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/jwt"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
)

type RefreshTokenOptions struct {
	RefreshTokenDuration time.Duration
	AccessTokenDuration  time.Duration
}
type refreshTokenUseCase struct {
	access  jwt.Provider
	refresh jwt.Provider
	opts    RefreshTokenOptions
}

func NewRefreshTokenUseCase(opts RefreshTokenOptions, accessProvider, refreshProvider jwt.Provider) auth.RefreshTokenUseCase {
	return &refreshTokenUseCase{
		access:  accessProvider,
		refresh: refreshProvider,
		opts:    opts,
	}
}

func (u *refreshTokenUseCase) Execute(ctx context.Context, params dto.RefreshTokenParams) (dto.AuthTokens, error) {
	ctx, span := trace.NewSpan(ctx, "auth.refreshToken")
	defer span.End()

	c, err := u.refresh.Validate(ctx, params.RefreshToken)
	if err != nil {
		logger.InfoF(ctx, "invalid refresh token", logger.Fields{
			"error": err.Error(),
		})
		return dto.AuthTokens{}, err
	}

	return generateAuthTokens(ctx, c.UserID.String(), u.access, u.refresh, u.opts.AccessTokenDuration, u.opts.RefreshTokenDuration)
}
