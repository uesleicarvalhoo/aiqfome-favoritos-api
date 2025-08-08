package ioc

import (
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/client"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/usecase"
)

var (
	listClientsUcOnce sync.Once
	listClientsUc     client.ListClientsUseCase
)

func ListClientsUseCase() client.ListClientsUseCase {
	listClientsUcOnce.Do(func() {
		listClientsUc = usecase.NewListClientsUseCase(UserRepository())
	})

	return listClientsUc
}

var (
	updateClientUcOnce sync.Once
	updateClientsUc    client.UpdateClientUseCase
)

func UpdateClientsUseCase() client.UpdateClientUseCase {
	updateClientUcOnce.Do(func() {
		updateClientsUc = usecase.NewUpdateClientUseCase(UserRepository())
	})

	return updateClientsUc
}

var (
	deleteClientUcOnce sync.Once
	deleteClientUc     client.DeleteClientUseCase
)

func DeleteClientUseCase() client.DeleteClientUseCase {
	deleteClientUcOnce.Do(func() {
		deleteClientUc = usecase.NewDeleteClientUseCase(UserRepository())
	})

	return deleteClientUc
}
