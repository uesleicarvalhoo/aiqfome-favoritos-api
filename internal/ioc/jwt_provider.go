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
		accessTokenProvider = jwt.NewProvider(jwt.Options{
			Issuer:    config.GetString("SERVICE_NAME"),
			Audiencer: config.GetString("SERVICE_NAME"),
			Secret:    config.GetString("ACESS_TOKEN_SECRET_KEY"),
		})
	})

	return accessTokenProvider
}

func RefreshTokenProvider() jwt.Provider {
	refreshTokenProviderOnce.Do(func() {
		refreshTokenProvider = jwt.NewProvider(jwt.Options{
			Issuer:    config.GetString("SERVICE_NAME"),
			Audiencer: config.GetString("SERVICE_NAME"),
			Secret:    config.GetString("REFRESH_TOKEN_SECRET_KEY"),
		})
	})

	return refreshTokenProvider
}
