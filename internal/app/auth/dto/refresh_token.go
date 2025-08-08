package dto

import "github.com/uesleicarvalhoo/aiqfome/pkg/validator"

type RefreshTokenParams struct {
	RefreshToken string `json:"refreshToken"`
}

func (p RefreshTokenParams) Validate() error {
	v := validator.New()

	if p.RefreshToken == "" {
		v.AddError("refreshToken", "campo obrigatório")
	}

	return v.Validate()
}
