package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/otel"
)

type Options struct {
	Host              string
	Port              string
	User              string
	Password          string
	Name              string
	PoolSize          int
	ConnMaxTTL        time.Duration
	TimeoutSeconds    int
	LockTimeoutMillis int
}

func (o Options) dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d statement_timeout=%ds lock_timeout=%d", //nolint: lll
		o.Host, o.Port, o.User, o.Password, o.Name, o.TimeoutSeconds, o.TimeoutSeconds, o.LockTimeoutMillis)
}

func NewPostgres(opts Options) (*sql.DB, error) {
	db, err := otelsql.Open("pgx", opts.dsn(),
		otelsql.WithDBName(opts.Name),
		otelsql.WithTracerProvider(otel.GetTracerProvider()),
	)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(opts.PoolSize)
	db.SetMaxOpenConns(opts.PoolSize)
	db.SetConnMaxLifetime(opts.ConnMaxTTL)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPostgresWithMigration(cfg Options) (*sql.DB, error) {
	db, err := NewPostgres(cfg)
	if err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
