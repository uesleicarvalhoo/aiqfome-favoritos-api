package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
)

type ListClientsParamsBuilder struct {
	page     int
	pageSize int
}

func AnyListClientsParams() ListClientsParamsBuilder {
	return ListClientsParamsBuilder{
		page:     1,
		pageSize: 20,
	}
}

func (b ListClientsParamsBuilder) WithPage(p int) ListClientsParamsBuilder {
	b.page = p
	return b
}

func (b ListClientsParamsBuilder) WithPageSize(size int) ListClientsParamsBuilder {
	b.pageSize = size
	return b
}

func (b ListClientsParamsBuilder) Build() dto.ListClientsParams {
	return dto.ListClientsParams{
		Page:     b.page,
		PageSize: b.pageSize,
	}
}
