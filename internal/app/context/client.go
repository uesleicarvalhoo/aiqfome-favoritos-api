package context

import (
	"context"

	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

var UserContext contextKey = "client_context"

func ContextWithUser(ctx context.Context, u user.User) context.Context {
	return context.WithValue(ctx, UserContext, u)
}

func GetClient(ctx context.Context) (user.User, error) {
	v := ctx.Value(UserContext)
	if c, ok := v.(user.User); ok {
		return c, nil
	}

	return user.User{}, domainerror.New(domainerror.AutenticationNotFound, "client não localizado no contexto", nil)
}
