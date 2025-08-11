package ioc

import (
	"fmt"
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/config"
	"github.com/uesleicarvalhoo/aiqfome/internal/infra/cache"
	cachePkg "github.com/uesleicarvalhoo/aiqfome/pkg/cache"
)

var (
	cacheOnce sync.Once
	cacheCli  cachePkg.Cache
)

func Cache() cachePkg.Cache {
	cacheOnce.Do(func() {
		cli, err := cache.NewRedis(cache.RedisOptions{
			Host:     config.GetString("REDIS_HOST"),
			Port:     config.GetString("REDIS_PORT"),
			User:     config.GetString("REDIS_USER"),
			Password: config.GetString("REDIS_PASSWORD"),
			UseSSL:   config.GetBool("REDIS_USE_SSL"),
		})
		if err != nil {
			panic(fmt.Sprintf("failed to setup redis: %s", err))
		}

		cacheCli = cli
	})

	return cacheCli
}
