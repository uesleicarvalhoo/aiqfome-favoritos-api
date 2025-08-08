package ioc

import (
	"strings"
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/config"
	"github.com/uesleicarvalhoo/aiqfome/internal/infra/requester"
	"github.com/uesleicarvalhoo/aiqfome/product"
	"github.com/uesleicarvalhoo/aiqfome/product/fakestoreapi"
)

var (
	productRepo     product.Repository
	productRepoOnce sync.Once
)

func ProductRepository() product.Repository {
	productRepoOnce.Do(func() {
		productRepo = fakestoreapi.NewRepository(
			requester.New(HttpClient()),
			fakestoreapi.Options{
				BaseUrl:         strings.TrimSuffix(config.GetString("FAKE_STORE_API_URL"), "/"),
				GetByIdEndpoint: config.GetString("FAKE_STORE_API_GET_BY_ID_ENDPOINT"),
				GetAllEndpoint:  config.GetString("FAKE_STORE_API_GET_ALL"),
			},
		)
	})

	return productRepo
}
