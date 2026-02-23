package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*/*.sql
var migrationsFS embed.FS

// Migrate runs the database migrations using the provided DB connection and provider.
// Provider should be one of "sqlite", "postgres", "mysql", "mssql" (or "sqlserver").
func Migrate(db *sql.DB, provider string, dbName string) error {
	driverName := provider
	if provider == "mssql" {
		driverName = "sqlserver"
	}

	var d database.Driver
	var err error

	switch driverName {
	case "sqlite":
		d, err = sqlite.WithInstance(db, &sqlite.Config{})
	case "postgres":
		d, err = postgres.WithInstance(db, &postgres.Config{})
	case "mysql":
		d, err = mysql.WithInstance(db, &mysql.Config{})
	case "sqlserver":
		d, err = sqlserver.WithInstance(db, &sqlserver.Config{DatabaseName: dbName})
	default:
		return fmt.Errorf("unsupported driver: %s", driverName)
	}

	if err != nil {
		return fmt.Errorf("could not create database driver: %w", err)
	}

	sourceDriver, err := iofs.New(migrationsFS, "migrations/"+provider)
	if err != nil {
		return fmt.Errorf("could not create source driver: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		driverName,
		d,
	)
	if err != nil {
		return fmt.Errorf("could not create migrator: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
