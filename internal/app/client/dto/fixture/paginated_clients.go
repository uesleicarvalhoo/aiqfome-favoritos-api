package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
)

type PaginatedClientsBuilder struct {
	clients []dto.Client
	total   int
	pages   int
}

func AnyPaginatedClients() PaginatedClientsBuilder {
	return PaginatedClientsBuilder{
		clients: []dto.Client{AnyClient().Build()},
		total:   1,
		pages:   1,
	}
}

func (b PaginatedClientsBuilder) WithProducts(clients []dto.Client) PaginatedClientsBuilder {
	b.clients = clients
	return b
}

func (b PaginatedClientsBuilder) WithTotal(total int) PaginatedClientsBuilder {
	b.total = total
	return b
}

func (b PaginatedClientsBuilder) WithPages(pages int) PaginatedClientsBuilder {
	b.pages = pages
	return b
}

func (b PaginatedClientsBuilder) Build() dto.PaginatedClients {
	return dto.PaginatedClients{
		Clients: b.clients,
		Total:   b.total,
		Pages:   b.pages,
	}
}
