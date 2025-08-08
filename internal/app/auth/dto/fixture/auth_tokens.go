package fixture

import "github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"

type AuthTokensBuilder struct {
	accessToken  string
	refreshToken string
}

func AnyAuthTokens() AuthTokensBuilder {
	return AuthTokensBuilder{
		accessToken:  "access.token",
		refreshToken: "refresh.token",
	}
}

func (b AuthTokensBuilder) WithAccessToken(token string) AuthTokensBuilder {
	b.accessToken = token
	return b
}

func (b AuthTokensBuilder) WithRefreshToken(token string) AuthTokensBuilder {
	b.refreshToken = token
	return b
}

func (b AuthTokensBuilder) Build() dto.AuthTokens {
	return dto.AuthTokens{
		AccessToken:  b.accessToken,
		RefreshToken: b.refreshToken,
	}
}
