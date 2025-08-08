package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/product"
	fixtureProduct "github.com/uesleicarvalhoo/aiqfome/product/fixture"
)

type ProductFavoriteBuilder struct {
	clientID uuid.ID
	product  product.Product
}

func AnyProductFavorite() ProductFavoriteBuilder {
	return ProductFavoriteBuilder{
		clientID: uuid.NextID(),
		product:  fixtureProduct.AnyProduct().Build(),
	}
}

func (b ProductFavoriteBuilder) WithClientID(id uuid.ID) ProductFavoriteBuilder {
	b.clientID = id
	return b
}

func (b ProductFavoriteBuilder) WithProduct(p product.Product) ProductFavoriteBuilder {
	b.product = p
	return b
}

func (b ProductFavoriteBuilder) Build() dto.ProductFavorite {
	return dto.ProductFavorite{
		ClientID: b.clientID,
		Product:  b.product,
	}
}
