package fixture

import "github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"

type SignUpParamsBuilder struct {
	email    string
	name     string
	password string
}

func AnySignUpParams() SignUpParamsBuilder {
	return SignUpParamsBuilder{
		email:    "user@example.com",
		name:     "Ueslei Carvalho",
		password: "secret",
	}
}

func (b SignUpParamsBuilder) WithEmail(email string) SignUpParamsBuilder {
	b.email = email
	return b
}

func (b SignUpParamsBuilder) WithName(name string) SignUpParamsBuilder {
	b.name = name
	return b
}

func (b SignUpParamsBuilder) WithPassword(pw string) SignUpParamsBuilder {
	b.password = pw
	return b
}

func (b SignUpParamsBuilder) Build() dto.SignUpParams {
	params := dto.SignUpParams{
		Email:    b.email,
		Name:     b.name,
		Password: b.password,
	}
	_ = params.Validate()
	return params
}
