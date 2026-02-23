//go:build integration
// +build integration

package stores_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/db"
	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/frostyeti/hyprship/apps/api/internal/stores/mssql"
	"github.com/frostyeti/hyprship/apps/api/internal/stores/mysql"
	"github.com/frostyeti/hyprship/apps/api/internal/stores/pg"
	"github.com/frostyeti/hyprship/apps/api/internal/stores/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tc_mssql "github.com/testcontainers/testcontainers-go/modules/mssql"
	tc_mysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	tc_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tc_wait "github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
	_ "modernc.org/sqlite"
)

func runStoreSuite(t *testing.T, userStore stores.UserStore, roleStore stores.RoleStore) {
	ctx := context.Background()

	// Test Users
	email := "test@example.com"
	name := "Test User"
	user := &models.User{
		ID:           uuid.New(),
		PrimaryEmail: &email,
		Name:         &name,
		IsBanned:     false,
		CreatedAt:    time.Now().Truncate(time.Second), // truncate to match SQL precision
	}

	err := userStore.Create(ctx, user)
	require.NoError(t, err)

	u, err := userStore.Get(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, u)
	assert.Equal(t, *user.Name, *u.Name)

	u2, err := userStore.GetByEmail(ctx, "TEST@example.com")
	require.NoError(t, err)
	require.NotNil(t, u2)
	assert.Equal(t, user.ID, u2.ID)

	// List users
	res, err := userStore.List(ctx, core.ListOptions{Filter: "test"})
	require.NoError(t, err)
	assert.Equal(t, int64(1), res.Total)
	assert.Len(t, res.Items, 1)

	// User Claims
	claim := &models.UserClaim{
		UserID: user.ID,
		Type:   "scope",
		Value:  "read:users",
	}
	err = userStore.AddClaim(ctx, claim)
	require.NoError(t, err)
	assert.NotZero(t, claim.ID)

	claims, err := userStore.ListClaims(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, claims, 1)

	err = userStore.RemoveClaim(ctx, claim.ID)
	require.NoError(t, err)

	// Test Roles
	role := &models.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator",
	}
	err = roleStore.Create(ctx, role)
	require.NoError(t, err)

	r, err := roleStore.Get(ctx, role.ID)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "admin", r.Name)

	role.Description = "Super Admin"
	err = roleStore.Update(ctx, role)
	require.NoError(t, err)

	roleRes, err := roleStore.List(ctx, core.ListOptions{})
	require.NoError(t, err)
	assert.Equal(t, int64(1), roleRes.Total)

	err = userStore.Delete(ctx, user.ID)
	require.NoError(t, err)
}

func TestStore_SQLite(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer conn.Close()

	err = db.Migrate(conn, "sqlite", "")
	require.NoError(t, err)

	runStoreSuite(t, sqlite.NewUserStore(conn), sqlite.NewRoleStore(conn))
}

func TestStore_Postgres(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := tc_postgres.Run(ctx,
		"postgres:15-alpine",
		tc_postgres.WithDatabase("testdb"),
		tc_postgres.WithUsername("testuser"),
		tc_postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			tc_wait.ForLog("database system is ready to accept connections").
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
	require.NoError(t, err)

	runStoreSuite(t, pg.NewUserStore(conn), pg.NewRoleStore(conn))
}

func TestStore_MySQL(t *testing.T) {
	ctx := context.Background()

	mysqlContainer, err := tc_mysql.Run(ctx,
		"mysql:8",
		tc_mysql.WithDatabase("testdb"),
		tc_mysql.WithUsername("testuser"),
		tc_mysql.WithPassword("testpass"),
	)
	require.NoError(t, err)
	defer func() { _ = mysqlContainer.Terminate(ctx) }()

	connStr, err := mysqlContainer.ConnectionString(ctx)
	require.NoError(t, err)
	connStr += "?multiStatements=true&parseTime=true"

	conn, err := sql.Open("mysql", connStr)
	require.NoError(t, err)
	defer conn.Close()

	err = db.Migrate(conn, "mysql", "testdb")
	require.NoError(t, err)

	runStoreSuite(t, mysql.NewUserStore(conn), mysql.NewRoleStore(conn))
}

func TestStore_MSSQL(t *testing.T) {
	ctx := context.Background()

	mssqlContainer, err := tc_mssql.Run(ctx,
		"mcr.microsoft.com/mssql/server:2022-latest",
		tc_mssql.WithAcceptEULA(),
		tc_mssql.WithPassword("StrongP@ssw0rd!"),
	)
	require.NoError(t, err)
	defer func() { _ = mssqlContainer.Terminate(ctx) }()

	connStr, err := mssqlContainer.ConnectionString(ctx)
	require.NoError(t, err)

	conn, err := sql.Open("sqlserver", connStr)
	require.NoError(t, err)
	defer conn.Close()

	err = db.Migrate(conn, "mssql", "")
	require.NoError(t, err)

	runStoreSuite(t, mssql.NewUserStore(conn), mssql.NewRoleStore(conn))
}
