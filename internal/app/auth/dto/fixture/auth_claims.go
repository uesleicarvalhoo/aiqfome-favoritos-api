package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type AuthClaimsBuilder struct {
	userID uuid.ID
}

func AnyAuthClaims() AuthClaimsBuilder {
	return AuthClaimsBuilder{
		userID: uuid.NextID(),
	}
}

func (b AuthClaimsBuilder) WithUserID(id uuid.ID) AuthClaimsBuilder {
	b.userID = id
	return b
}

func (b AuthClaimsBuilder) Build() dto.AuthClaims {
	return dto.AuthClaims{
		UserID: b.userID,
	}
}
