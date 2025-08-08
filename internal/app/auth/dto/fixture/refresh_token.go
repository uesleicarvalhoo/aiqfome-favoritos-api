package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
)

type RefreshTokenParamsBuilder struct {
	refreshToken string
}

func AnyRefreshTokenParams() RefreshTokenParamsBuilder {
	return RefreshTokenParamsBuilder{
		refreshToken: "refresh.token",
	}
}

func (b RefreshTokenParamsBuilder) WithRefreshToken(token string) RefreshTokenParamsBuilder {
	b.refreshToken = token
	return b
}

func (b RefreshTokenParamsBuilder) Build() dto.RefreshTokenParams {
	return dto.RefreshTokenParams{
		RefreshToken: b.refreshToken,
	}
}
