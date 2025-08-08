package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/role"
)

type UpdateClientParamsBuilder struct {
	clientID uuid.ID
	name     *string
	active   *bool
	role     *role.Role
}

func AnyUpdateClientParams() UpdateClientParamsBuilder {
	return UpdateClientParamsBuilder{
		clientID: uuid.NextID(),
		name:     nil,
		active:   nil,
		role:     nil,
	}
}

func (b UpdateClientParamsBuilder) WithClientID(id uuid.ID) UpdateClientParamsBuilder {
	b.clientID = id
	return b
}

func (b UpdateClientParamsBuilder) WithName(name string) UpdateClientParamsBuilder {
	b.name = &name
	return b
}

func (b UpdateClientParamsBuilder) WithActive(active bool) UpdateClientParamsBuilder {
	b.active = &active
	return b
}

func (b UpdateClientParamsBuilder) WithRole(role role.Role) UpdateClientParamsBuilder {
	b.role = &role
	return b
}

func (b UpdateClientParamsBuilder) Build() dto.UpdateClientParams {
	return dto.UpdateClientParams{
		ClientID: b.clientID,
		Name:     b.name,
		Active:   b.active,
		Role:     b.role,
	}
}
