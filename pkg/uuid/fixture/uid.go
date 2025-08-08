package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type UUIDBuilder struct {
	uid uuid.ID
}

func AnyUUID() UUIDBuilder {
	return UUIDBuilder{
		uid: uuid.MustParse("c1c43927-8c55-4276-a9c9-f27cae2ee332"),
	}
}

func (b UUIDBuilder) Build() uuid.ID {
	return b.uid
}
