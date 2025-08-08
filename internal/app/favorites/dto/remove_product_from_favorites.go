package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
)

type RemoveProductFromFavoritesParams struct {
	ClientID  uuid.ID `json:"clientId"`
	ProductID int     `json:"productId"`
}

func (p RemoveProductFromFavoritesParams) Validate() error {
	v := validator.New()

	if p.ClientID.IsZero() {
		v.AddError("clientId", "campo obrigatório")
	}

	if p.ProductID == 0 {
		v.AddError("productId", "campo obrigatório")
	}

	return v.Validate()
}
