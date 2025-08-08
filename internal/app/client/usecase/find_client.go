package usecase

import (
	"context"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/client"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type findClientUseCase struct {
	repo user.Repository
}

func NewFindClientUseCase(repo user.Repository) client.FindClientUseCase {
	return &findClientUseCase{
		repo: repo,
	}
}

func (u *findClientUseCase) Execute(ctx context.Context, id uuid.ID) (dto.Client, error) {
	ctx, span := trace.NewSpan(ctx, "client.FindClient")
	defer span.End()

	usr, err := u.repo.Find(ctx, id)
	if err != nil {
		logger.ErrorF(ctx, "error while trying to find client", logger.Fields{
			"error":   err.Error(),
			"user_id": id,
		})
	}

	return dto.FromDomain(usr), nil
}
