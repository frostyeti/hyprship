package db

import (
	"database/sql"
	"fmt"

	"github.com/frostyeti/hyprship/apps/api/config"

	// Drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
)

func Connect(cfg *config.DBConfig) (*sql.DB, error) {
	driver := cfg.Driver
	if driver == "sqlite" {
		driver = "sqlite3"
	} else if driver == "pg" || driver == "postgres" {
		driver = "postgres"
	} else if driver == "sqlserver" || driver == "mssql" {
		driver = "sqlserver"
	}

	db, err := sql.Open(driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
