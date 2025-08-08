package main

import (
	"context"

	"github.com/uesleicarvalhoo/aiqfome/config"
	"github.com/uesleicarvalhoo/aiqfome/internal/http"
	"github.com/uesleicarvalhoo/aiqfome/internal/ioc"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
)

func main() {
	provider, err := trace.NewProvider(
		trace.ProviderConfig{
			Endpoint:       config.GetString("TRACER_ENDPOINT"),
			ServiceName:    config.GetString("SERVICE_NAME"),
			ServiceVersion: config.GetString("SERVICE_VERSION"),
			Environment:    config.GetString("ENVIRONMENT"),
			Disabled:       !config.GetBool("TRACE_ENABLED"),
		})
	if err != nil {
		logger.Fatal(context.Background(), "couldn't connect to provider: %s", err)
	}
	defer provider.Close(context.Background())

	authenticateUc := ioc.AuthenticateUseCase()
	authorizeUc := ioc.AuthorizeUseCase()
	signInUc := ioc.SignInUseCase()
	signUpUc := ioc.SignUpUseCase()
	refreshTokenUc := ioc.RefreshTokenUseCase()
	getClientFavoritesUc := ioc.GetClientFavoritesUseCase()
	addProductToFavoritesUc := ioc.AddProductToFavoritesUseCase()
	removeProductFromFavoritesUc := ioc.RemoveProductFromFavoritesUseCase()
	listClientsUc := ioc.ListClientsUseCase()
	updateClientUc := ioc.UpdateClientsUseCase()
	deleteClientUc := ioc.DeleteClientUseCase()

	err = http.StartHttpServer(http.Options{
		ServiceName: config.GetString("SERVICE_NAME"),
		Port:        config.GetInt("HTTP_SERVER_PORT"),
	},
		authenticateUc,
		authorizeUc,
		signInUc,
		signUpUc,
		refreshTokenUc,
		getClientFavoritesUc,
		addProductToFavoritesUc,
		removeProductFromFavoritesUc,
		listClientsUc,
		updateClientUc,
		deleteClientUc,
	)
	if err != nil {
		panic(err)
	}
}
