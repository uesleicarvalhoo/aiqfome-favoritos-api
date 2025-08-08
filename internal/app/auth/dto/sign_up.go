package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
)

type SignUpParams struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (p SignUpParams) Validate() error {
	v := validator.New()

	if !validator.IsEmailValid(p.Email) {
		v.AddError("email", "email inválido")
	}

	if p.Name == "" {
		v.AddError("nome", "campo obrigatório")
	}

	if p.Password == "" {
		v.AddError("password", "campo obrigatório")
	}

	return v.Validate()
}
