package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
	"github.com/uesleicarvalhoo/aiqfome/product"
)

type GetClientFavoritesParams struct {
	ClientID uuid.ID `json:"-"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}

func (p GetClientFavoritesParams) Validate() error {
	v := validator.New()

	if p.ClientID.IsZero() {
		v.AddError("clientId", "campo obrigatório")
	}

	if p.PageSize < 1 {
		v.AddError("pageSize", "deve ser maior do que 1")
	}

	if p.Page < 0 {
		v.AddError("page", "não pode ser negativo")
	}

	return v.Validate()
}

type ClientFavorites struct {
	ClientID uuid.ID           `json:"clientId"`
	Products []product.Product `json:"products"`
	Total    int               `json:"total"`
	Pages    int               `json:"pages"`
}
