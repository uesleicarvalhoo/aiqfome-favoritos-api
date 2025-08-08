package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type GetClientFavoritesParamsBuilder struct {
	clientID uuid.ID
	page     int
	pageSize int
}

func AnyGetClientFavoritesParams() GetClientFavoritesParamsBuilder {
	return GetClientFavoritesParamsBuilder{
		clientID: uuid.NextID(),
		page:     1,
		pageSize: 20,
	}
}

func (b GetClientFavoritesParamsBuilder) WithClientID(id uuid.ID) GetClientFavoritesParamsBuilder {
	b.clientID = id
	return b
}

func (b GetClientFavoritesParamsBuilder) WithPage(p int) GetClientFavoritesParamsBuilder {
	b.page = p
	return b
}

func (b GetClientFavoritesParamsBuilder) WithPageSize(size int) GetClientFavoritesParamsBuilder {
	b.pageSize = size
	return b
}

func (b GetClientFavoritesParamsBuilder) Build() dto.GetClientFavoritesParams {
	return dto.GetClientFavoritesParams{
		ClientID: b.clientID,
		Page:     b.page,
		PageSize: b.pageSize,
	}
}
