package ioc

import (
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

var (
	idGeneratorOnce sync.Once
	idGenerator     uuid.Generator
)

func IDGenerator() uuid.Generator {
	idGeneratorOnce.Do(func() {
		idGenerator = uuid.NewGenerator()
	})

	return idGenerator
}
