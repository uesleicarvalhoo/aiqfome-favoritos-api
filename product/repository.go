package product

import "context"

type Reader interface {
	Find(ctx context.Context, id int) (Product, error)
	FindMultiple(ctx context.Context, ids []int) ([]Product, error)
}

type Repository interface {
	Reader
}
