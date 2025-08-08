package ioc

import (
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/config"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/usecase"
)

var (
	signInUcOnce sync.Once
	signUc       auth.SignInUseCase
)

func SignInUseCase() auth.SignInUseCase {
	signInUcOnce.Do(func() {
		signUc = usecase.NewSignUseCase(
			UserRepository(),
			PasswordHasher(),
			usecase.SignInOptions{
				AccessTokenDuration:  config.GetDuration("ACCESS_TOKEN_DURATION"),
				RefreshTokenDuration: config.GetDuration("REFRESH_TOKEN_DURATION"),
			},
			AccessTokenProvider(),
			RefreshTokenProvider(),
		)
	})

	return signUc
}

var (
	signUpUcOnce sync.Once
	signUpUc     auth.SignUpUseCase
)

func SignUpUseCase() auth.SignUpUseCase {
	signUpUcOnce.Do(func() {
		signUpUc = usecase.NewSignUpUseCase(
			IDGenerator(),
			PasswordHasher(),
			UserRepository(),
			usecase.SignUpOptions{
				MinPasswordLength: config.GetInt("MIN_PASSWORD_LENGTH"),
			})
	})

	return signUpUc
}

var (
	refreshTokenUcOnce sync.Once
	refreshTokenUc     auth.RefreshTokenUseCase
)

func RefreshTokenUseCase() auth.RefreshTokenUseCase {
	refreshTokenUcOnce.Do(func() {
		refreshTokenUc = usecase.NewRefreshTokenUseCase(usecase.RefreshTokenOptions{
			AccessTokenDuration:  config.GetDuration("ACCESS_TOKEN_DURATION"),
			RefreshTokenDuration: config.GetDuration("REFRESH_TOKEN_DURATION"),
		},
			AccessTokenProvider(),
			RefreshTokenProvider(),
		)
	})

	return refreshTokenUc
}

var (
	authenticateUc     auth.AuthenticateUseCase
	authenticateUcOnce sync.Once
)

func AuthenticateUseCase() auth.AuthenticateUseCase {
	authenticateUcOnce.Do(func() {
		authenticateUc = usecase.NewAuthenticateUseCase(UserRepository(), AccessTokenProvider())
	})

	return authenticateUc
}

var (
	authorizeUc     auth.AuthorizeUseCase
	authorizeUcOnce sync.Once
)

func AuthorizeUseCase() auth.AuthorizeUseCase {
	authorizeUcOnce.Do(func() {
		authorizeUc = usecase.NewAuthorizeUseCase(RoleRepository())
	})

	return authorizeUc
}
