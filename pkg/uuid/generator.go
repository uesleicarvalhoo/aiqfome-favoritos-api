package uuid

import "github.com/google/uuid"

type Generator interface {
	NextID() ID
}

type generator struct{}

func NewGenerator() Generator {
	return &generator{}
}

func (g *generator) NextID() ID {
	return ID(uuid.New())
}

var idGenerator Generator = NewGenerator()
