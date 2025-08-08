package client

import (
	"context"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
)

type ListClientsUseCase interface {
	Execute(ctx context.Context, p dto.ListClientsParams) (dto.PaginatedClients, error)
}

type UpdateClientUseCase interface {
	Execute(ctx context.Context, p dto.UpdateClientParams) (dto.Client, error)
}

type DeleteClientUseCase interface {
	Execute(ctx context.Context, p dto.DeleteClientParams) error
}
