package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/context"
	"github.com/uesleicarvalhoo/aiqfome/internal/http/utils"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/role"
)

func Authentication(uc auth.AuthenticateUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.WriteError(c, domainerror.New(
				domainerror.AutenticationNotFound,
				"token não informado ou formato inválido",
				nil,
			))
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		cl, err := uc.Execute(c.UserContext(), token)
		if err != nil {
			return utils.WriteError(c, err)
		}

		ctx := context.ContextWithUser(c.UserContext(), cl)
		ctx = logger.ContextWithFields(ctx, logger.Fields{
			"client_id": cl.ID,
		})

		c.SetUserContext(ctx)

		return c.Next()
	}
}

func Authorize(uc auth.AuthorizeUseCase, resource role.Resource, action role.Action) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl, err := context.GetClient(c.UserContext())
		if err != nil {
			return utils.WriteError(c, err)
		}

		if err := uc.Execute(c.UserContext(), dto.AuthorizeParams{
			User:     cl,
			Resource: resource,
			Action:   action,
		}); err != nil {
			return utils.WriteError(c, err)
		}

		return c.Next()
	}
}
