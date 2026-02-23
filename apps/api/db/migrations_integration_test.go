//go:build integration
// +build integration

package db_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/frostyeti/hyprship/apps/api/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mssql"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
	_ "modernc.org/sqlite"
)

func TestMigrations_SQLite(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer conn.Close()

	err = db.Migrate(conn, "sqlite", "")
	assert.NoError(t, err)
}

func TestMigrations_Postgres(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second)),
	)
	require.NoError(t, err)
	defer func() { _ = pgContainer.Terminate(ctx) }()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	conn, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer conn.Close()

	err = db.Migrate(conn, "postgres", "testdb")
	assert.NoError(t, err)
}

func TestMigrations_MySQL(t *testing.T) {
	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8",
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
	)
	require.NoError(t, err)
	defer func() { _ = mysqlContainer.Terminate(ctx) }()

	connStr, err := mysqlContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Enable multiStatements for migrations
	connStr += "?multiStatements=true"

	conn, err := sql.Open("mysql", connStr)
	require.NoError(t, err)
	defer conn.Close()

	err = db.Migrate(conn, "mysql", "testdb")
	assert.NoError(t, err)
}

func TestMigrations_MSSQL(t *testing.T) {
	ctx := context.Background()

	mssqlContainer, err := mssql.Run(ctx,
		"mcr.microsoft.com/mssql/server:2022-latest",
		mssql.WithAcceptEULA(),
		mssql.WithPassword("StrongP@ssw0rd!"),
	)
	require.NoError(t, err)
	defer func() { _ = mssqlContainer.Terminate(ctx) }()

	connStr, err := mssqlContainer.ConnectionString(ctx)
	require.NoError(t, err)

	conn, err := sql.Open("sqlserver", connStr)
	require.NoError(t, err)
	defer conn.Close()

	err = db.Migrate(conn, "mssql", "")
	assert.NoError(t, err)
}
