package middleware

import (
	"fmt"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
)

func Otel(serverName string) fiber.Handler {
	mid := otelfiber.Middleware(
		otelfiber.WithServerName(serverName),
		otelfiber.WithSpanNameFormatter(
			func(c *fiber.Ctx) string {
				return fmt.Sprintf("%s - %s", c.Route().Method, c.Route().Path)
			}))

	return func(c *fiber.Ctx) error {
		if c.Path() == "/health" {
			return c.Next()
		}

		ctx := c.UserContext()
		span := trace.SpanFromContext(ctx)

		logger.ContextWithFields(ctx, logger.Fields{
			"trace_id": span.SpanContext().TraceID().String(),
			"span_id":  span.SpanContext().SpanID().String(),
		})

		c.SetUserContext(ctx)

		return mid(c)
	}
}
