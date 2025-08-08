package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
)

type AddProductToFavoritesParams struct {
	ClientID  uuid.ID `json:"-"`
	ProductID int     `json:"productId"`
}

func (p AddProductToFavoritesParams) Validate() error {
	v := validator.New()

	if p.ClientID.IsZero() {
		v.AddError("clientId", "campo obrigatório")
	}

	if p.ProductID == 0 {
		v.AddError("productId", "campo obrigatório")
	}

	return v.Validate()
}
