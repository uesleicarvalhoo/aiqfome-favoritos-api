package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
)

type DeleteClientParams struct {
	ClientID uuid.ID `json:"clientId"`
}

func (p DeleteClientParams) Validate() error {
	v := validator.New()
	if p.ClientID.IsZero() {
		v.AddError("clientId", "campo obrigatório")
	}

	return v.Validate()
}
