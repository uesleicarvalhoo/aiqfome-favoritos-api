package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/product"
)

type ProductFavorite struct {
	ClientID uuid.ID         `json:"clientId"`
	Product  product.Product `json:"product"`
}
