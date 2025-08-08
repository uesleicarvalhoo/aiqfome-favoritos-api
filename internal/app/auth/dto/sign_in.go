package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
)

type SignInParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p SignInParams) Validate() error {
	v := validator.New()

	if !validator.IsEmailValid(p.Email) {
		v.AddError("email", "email inválido")
	}

	if p.Password == "" {
		v.AddError("password", "campo obrigatório")
	}

	return v.Validate()
}
