package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
)

func Context() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fields := logger.Fields{
			"http.path":   c.Path(),
			"http.method": c.Method(),
			"request_id":  c.Locals(requestid.ConfigDefault.ContextKey),
		}

		ctx := logger.ContextWithFields(c.UserContext(), fields)

		c.SetUserContext(ctx)

		return c.Next()
	}
}
