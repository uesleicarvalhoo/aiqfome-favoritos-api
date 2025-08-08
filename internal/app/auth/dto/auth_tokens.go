package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
	"github.com/uesleicarvalhoo/aiqfome/role"
)

type AuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthClaims struct {
	UserID uuid.ID   `json:"userId"`
	Role   role.Role `json:"role"`
}

func (a AuthClaims) Validate() error {
	v := validator.New()

	if a.UserID.IsZero() {
		v.AddError("userId", "ID do usuário é obrigatório")
	}

	return v.Validate()
}
