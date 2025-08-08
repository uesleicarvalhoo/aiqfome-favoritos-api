package http

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites"
	"github.com/uesleicarvalhoo/aiqfome/internal/http/middleware"
	"github.com/uesleicarvalhoo/aiqfome/internal/http/routes"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// @title Aiqfome api challenge
// @version 1.0.0
// @description This is api for magalu challenge
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Enter your Bearer token in the format: `Bearer {token}`"
func StartHttpServer(opts Options,
	authenticateUc auth.AuthenticateUseCase,
	authorizeUc auth.AuthorizeUseCase,
	signInUc auth.SignInUseCase,
	signUpUc auth.SignUpUseCase,
	refreshTokenUc auth.RefreshTokenUseCase,
	getClientFavoritesUc favorites.GetClientFavoritesUseCase,
	addProductToFavoritesUc favorites.AddProductToFavoritesUseCase,
	removeProductFromFavoritesUc favorites.RemoveProductFromFavoritesUseCase,
	listClientsUc client.ListClientsUseCase,
	updateClientUc client.UpdateClientUseCase,
	deleteClientUc client.DeleteClientUseCase,
) error {
	app := fiber.New(fiber.Config{
		AppName:               opts.ServiceName,
		DisableStartupMessage: true,
	})

	app.Use(
		recover.New(),
		cors.New(),
		requestid.New(),
		middleware.Context(),
		middleware.Logger(),
		middleware.Otel(opts.ServiceName),
	)

	routes.Swagger(app)
	routes.Auth(app.Group("/auth"), signInUc, signUpUc, refreshTokenUc)

	protected := app.Group("/", middleware.Authentication(authenticateUc))

	routes.Me(
		protected.Group("/me"),
		getClientFavoritesUc, addProductToFavoritesUc, removeProductFromFavoritesUc,
	)

	routes.Clients(
		protected.Group("/clients"),
		authorizeUc,
		listClientsUc,
		updateClientUc,
		deleteClientUc,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		defer cancel()
		logger.Debug(context.Background(), "http server running on port: %d", opts.Port)

		err := app.Listen(fmt.Sprintf(":%d", opts.Port))
		if err != nil {
			logger.ErrorF(context.Background(), "error while starting http server", logger.Fields{
				"error": err.Error(),
			})
		}
	}()

	<-ctx.Done()
	logger.Debug(context.Background(), "shutting down http server")
	if err := app.Shutdown(); err != nil {
		logger.ErrorF(context.Background(), "failed to shutting down http server", logger.Fields{
			"error": err.Error(),
		})

		return err
	}

	return nil
}
