package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type ClientBuilder struct {
	id     uuid.ID
	name   string
	email  string
	active bool
}

func AnyClient() ClientBuilder {
	return ClientBuilder{
		id:     uuid.NextID(),
		name:   "client name",
		email:  "client@email.com",
		active: true,
	}
}

func (b ClientBuilder) WithID(id uuid.ID) ClientBuilder {
	b.id = id
	return b
}

func (b ClientBuilder) WithName(n string) ClientBuilder {
	b.name = n
	return b
}

func (b ClientBuilder) WithEmail(e string) ClientBuilder {
	b.email = e
	return b
}

func (b ClientBuilder) WithActive(a bool) ClientBuilder {
	b.active = a
	return b
}

func (b ClientBuilder) Build() dto.Client {
	return dto.Client{
		ID:     b.id,
		Name:   b.name,
		Email:  b.email,
		Active: b.active,
	}
}
