package routes

import (
	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
	_ "github.com/uesleicarvalhoo/aiqfome/docs"
)

func Swagger(r fiber.Router) {
	r.Get("/docs/*", swagger.WrapHandler)
}
