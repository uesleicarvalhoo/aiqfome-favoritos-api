package user

import (
	"context"

	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type Reader interface {
	Find(ctx context.Context, id uuid.ID) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	Paginate(ctx context.Context, page, pageSize int) ([]User, int, error)
}

type Writer interface {
	Create(ctx context.Context, c User) error
	Delete(ctx context.Context, c User) error
	Update(ctx context.Context, c User) error
}

type Repository interface {
	Reader
	Writer
}
