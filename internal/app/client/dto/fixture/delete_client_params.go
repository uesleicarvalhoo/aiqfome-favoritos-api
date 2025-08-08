package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type DeleteClientParamsBuilder struct {
	clientID uuid.ID
}

func AnyDeleteClientParams() DeleteClientParamsBuilder {
	return DeleteClientParamsBuilder{
		clientID: uuid.NextID(),
	}
}

func (b DeleteClientParamsBuilder) WithClientID(id uuid.ID) DeleteClientParamsBuilder {
	b.clientID = id
	return b
}

func (b DeleteClientParamsBuilder) Build() dto.DeleteClientParams {
	return dto.DeleteClientParams{
		ClientID: b.clientID,
	}
}
