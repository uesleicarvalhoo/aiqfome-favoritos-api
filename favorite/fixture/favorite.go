package fixture

import (
	"time"

	"github.com/uesleicarvalhoo/aiqfome/favorite"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type FavoriteBuilder struct {
	clientID    uuid.ID
	productID   int
	registredAt time.Time
}

func AnyFavorite() FavoriteBuilder {
	return FavoriteBuilder{
		clientID:    uuid.NextID(),
		productID:   1,
		registredAt: time.Now(),
	}
}

func (b FavoriteBuilder) WithClientID(id uuid.ID) FavoriteBuilder {
	b.clientID = id
	return b
}

func (b FavoriteBuilder) WithProductID(pid int) FavoriteBuilder {
	b.productID = pid
	return b
}

func (b FavoriteBuilder) WithRegistredAt(t time.Time) FavoriteBuilder {
	b.registredAt = t
	return b
}

func (b FavoriteBuilder) Build() favorite.Favorite {
	return favorite.Favorite{
		ClientID:    b.clientID,
		ProductID:   b.productID,
		RegistredAt: b.registredAt,
	}
}
