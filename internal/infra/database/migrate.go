package database

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	_ "github.com/golang-migrate/migrate/v4/source/file" // Load migration files
	"github.com/pressly/goose"
)

func Migrate(db *sql.DB) error {
	if err := goose.Up(db, getMigrationsPath()); err != nil {
		return fmt.Errorf("failed to apply migrations up: %w", err)
	}

	return nil
}

func getMigrationsPath() string {
	_, filename, _, _ := runtime.Caller(1)

	basePath := filepath.Dir(filename)
	basePath = filepath.Dir(filepath.Dir(filepath.Dir(basePath)))

	return filepath.Join(basePath, ".migrations")
}
