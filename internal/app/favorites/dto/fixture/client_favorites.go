package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/product"
	fixtureProduct "github.com/uesleicarvalhoo/aiqfome/product/fixture"
)

type ClientFavoritesBuilder struct {
	clientID uuid.ID
	products []product.Product
	total    int
	pages    int
}

func AnyClientFavorites() ClientFavoritesBuilder {
	return ClientFavoritesBuilder{
		clientID: uuid.NextID(),
		products: []product.Product{fixtureProduct.AnyProduct().Build()},
		total:    1,
		pages:    1,
	}
}

func (b ClientFavoritesBuilder) WithClientID(id uuid.ID) ClientFavoritesBuilder {
	b.clientID = id
	return b
}

func (b ClientFavoritesBuilder) WithProducts(prods []product.Product) ClientFavoritesBuilder {
	b.products = prods
	return b
}

func (b ClientFavoritesBuilder) WithTotal(total int) ClientFavoritesBuilder {
	b.total = total
	return b
}

func (b ClientFavoritesBuilder) WithPages(pages int) ClientFavoritesBuilder {
	b.pages = pages
	return b
}

func (b ClientFavoritesBuilder) Build() dto.ClientFavorites {
	return dto.ClientFavorites{
		ClientID: b.clientID,
		Products: b.products,
		Total:    b.total,
		Pages:    b.pages,
	}
}
