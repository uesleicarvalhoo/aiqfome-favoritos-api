package fixture

import "github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"

type SignInParamsBuilder struct {
	email    string
	password string
}

func AnySignInParams() SignInParamsBuilder {
	return SignInParamsBuilder{
		email:    "user@example.com",
		password: "secret",
	}
}

func (b SignInParamsBuilder) WithEmail(email string) SignInParamsBuilder {
	b.email = email
	return b
}

func (b SignInParamsBuilder) WithPassword(pw string) SignInParamsBuilder {
	b.password = pw
	return b
}

func (b SignInParamsBuilder) Build() dto.SignInParams {
	return dto.SignInParams{
		Email:    b.email,
		Password: b.password,
	}
}
