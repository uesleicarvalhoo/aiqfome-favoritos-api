package ioc

import (
	"database/sql"
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/config"
	"github.com/uesleicarvalhoo/aiqfome/internal/infra/database"
)

var (
	dbOnce sync.Once
	db     *sql.DB
)

func Database() *sql.DB {
	dbOnce.Do(func() {
		con, err := database.NewPostgres(database.Options{
			Host:              config.GetString("DATABASE_HOST"),
			Port:              config.GetString("DATABASE_PORT"),
			User:              config.GetString("DATABASE_USER"),
			Password:          config.GetString("DATABASE_PASSWORD"),
			Name:              config.GetString("DATABASE_NAME"),
			PoolSize:          config.GetInt("DATABASE_POOL_SIZE"),
			ConnMaxTTL:        config.GetDuration("DATABASE_CONN_MAX_TTL"),
			TimeoutSeconds:    config.GetInt("DATABASE_TIMEOUT_SECONDS"),
			LockTimeoutMillis: config.GetInt("DATABASE_LOCK_TIMEOUT_MILLIS"),
		})
		if err != nil {
			panic(err)
		}

		db = con
	})

	return db
}
