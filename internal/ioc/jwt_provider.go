package ioc

import (
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/config"
	"github.com/uesleicarvalhoo/aiqfome/pkg/jwt"
)

var (
	accessTokenProvider     jwt.Provider
	accessTokenProviderOnce sync.Once

	refreshTokenProvider     jwt.Provider
	refreshTokenProviderOnce sync.Once
)

func AccessTokenProvider() jwt.Provider {
	accessTokenProviderOnce.Do(func() {
		key := config.GetString("ACESS_TOKEN_SECRET_KEY")
		if key == "" {
			panic("failed to setup access token provider, env `ACCESS_TOKEN_SECRET_KEY` is missing")
		}

		accessTokenProvider = jwt.NewProvider(jwt.Options{
			Issuer:    config.GetString("SERVICE_NAME"),
			Audiencer: config.GetString("SERVICE_NAME"),
			Secret:    key,
		})
	})

	return accessTokenProvider
}

func RefreshTokenProvider() jwt.Provider {
	refreshTokenProviderOnce.Do(func() {
		key := config.GetString("REFRESH_TOKEN_SECRET_KEY")
		if key == "" {
			panic("failed to setup refresh token provider, env `REFRESH_TOKEN_SECRET_KEY` is missing")
		}

		refreshTokenProvider = jwt.NewProvider(jwt.Options{
			Issuer:    config.GetString("SERVICE_NAME"),
			Audiencer: config.GetString("SERVICE_NAME"),
			Secret:    key,
		})
	})

	return refreshTokenProvider
}
