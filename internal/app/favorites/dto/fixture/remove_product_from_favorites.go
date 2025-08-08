package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type RemoveProductFromFavoritesParamsBuilder struct {
	clientID  uuid.ID
	productID int
}

func AnyRemoveProductFromFavoritesParams() RemoveProductFromFavoritesParamsBuilder {
	return RemoveProductFromFavoritesParamsBuilder{
		clientID:  uuid.NextID(),
		productID: 1,
	}
}

func (b RemoveProductFromFavoritesParamsBuilder) WithClientID(id uuid.ID) RemoveProductFromFavoritesParamsBuilder {
	b.clientID = id
	return b
}

func (b RemoveProductFromFavoritesParamsBuilder) WithProductID(pid int) RemoveProductFromFavoritesParamsBuilder {
	b.productID = pid
	return b
}

func (b RemoveProductFromFavoritesParamsBuilder) Build() dto.RemoveProductFromFavoritesParams {
	return dto.RemoveProductFromFavoritesParams{
		ClientID:  b.clientID,
		ProductID: b.productID,
	}
}
