package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
)

func Logger(pathsToIgnore ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		for _, path := range pathsToIgnore {
			if strings.Contains(c.Path(), path) {
				return c.Next()
			}
		}

		if err := c.Next(); err != nil {
			return err
		}

		statusCode := c.Response().StatusCode()

		switch {
		case statusCode >= http.StatusInternalServerError:
			logger.ErrorF(c.UserContext(), fmt.Sprintf("%s - %s", c.Method(), c.Path()), logger.Fields{"status_code": statusCode})
		case statusCode >= http.StatusBadRequest:
			logger.WarnF(c.UserContext(), fmt.Sprintf("%s - %s", c.Method(), c.Path()), logger.Fields{"status_code": statusCode})
		default:
			logger.InfoF(c.UserContext(), fmt.Sprintf("%s - %s", c.Method(), c.Path()), logger.Fields{"status_code": statusCode})
		}

		return nil
	}
}
