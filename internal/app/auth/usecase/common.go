package usecase

import (
	"context"
	"time"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/jwt"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
)

func generateAuthTokens(ctx context.Context, sub string, accessProvider, refreshProvider jwt.Provider, accessDuration, refreshDuration time.Duration) (dto.AuthTokens, error) {
	at, err := accessProvider.Generate(ctx, sub, accessDuration)
	if err != nil {
		logger.InfoF(ctx, "error while trying to generate access token", logger.Fields{
			"error": err.Error(),
		})

		return dto.AuthTokens{}, err
	}

	rt, err := refreshProvider.Generate(ctx, sub, refreshDuration)
	if err != nil {
		logger.InfoF(ctx, "error while trying to generate refresh token", logger.Fields{
			"error": err.Error(),
		})

		return dto.AuthTokens{}, err
	}

	return dto.AuthTokens{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}
