package role

import (
	"context"
)

type Repository interface {
	FindPermissions(ctx context.Context, r Role) ([]Permission, error)
}
