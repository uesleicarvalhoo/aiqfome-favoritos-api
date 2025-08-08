package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type AddProductToFavoritesParamsBuilder struct {
	clientID  uuid.ID
	productID int
}

func AnyAddProductToFavoritesParams() AddProductToFavoritesParamsBuilder {
	return AddProductToFavoritesParamsBuilder{
		clientID:  uuid.NextID(),
		productID: 1,
	}
}

func (b AddProductToFavoritesParamsBuilder) WithClientID(id uuid.ID) AddProductToFavoritesParamsBuilder {
	b.clientID = id
	return b
}

func (b AddProductToFavoritesParamsBuilder) WithProductID(pid int) AddProductToFavoritesParamsBuilder {
	b.productID = pid
	return b
}

func (b AddProductToFavoritesParamsBuilder) Build() dto.AddProductToFavoritesParams {
	return dto.AddProductToFavoritesParams{
		ClientID:  b.clientID,
		ProductID: b.productID,
	}
}
