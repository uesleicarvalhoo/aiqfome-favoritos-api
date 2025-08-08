package favorite

import (
	"context"

	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type Reader interface {
	Find(ctx context.Context, clientID uuid.ID, productID int) (Favorite, error)
	PaginateByClientID(ctx context.Context, clientID uuid.ID, page, pageSize int) ([]Favorite, int, error)
}

type Writer interface {
	Create(ctx context.Context, f Favorite) error
	Remove(ctx context.Context, f Favorite) error
}

type Repository interface {
	Reader
	Writer
}
