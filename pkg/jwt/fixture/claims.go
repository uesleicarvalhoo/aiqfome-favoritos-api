package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/jwt"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type ClaimsBuilder struct {
	clientID uuid.ID
}

func AnyClaims() ClaimsBuilder {
	return ClaimsBuilder{
		clientID: uuid.NextID(),
	}
}

func (b ClaimsBuilder) WithClientID(id uuid.ID) ClaimsBuilder {
	b.clientID = id
	return b
}

func (b ClaimsBuilder) Build() jwt.Claims {
	return jwt.Claims{
		UserID: b.clientID,
	}
}
