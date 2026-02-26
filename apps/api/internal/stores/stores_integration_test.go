//go:build integration
// +build integration

package stores_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/db"
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
	res, err := userStore.List(ctx, core.ListOptions{Filter: []core.FilterOption{{Field: "name", Operator: "like", Value: "%TEST%"}}})
	require.NoError(t, err)
	assert.Equal(t, 1, res.TotalCount)
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
	assert.Equal(t, 1, roleRes.TotalCount)

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
	runGroupStoreSuite(t, sqlite.NewGroupStore(conn), sqlite.NewUserStore(conn), sqlite.NewRoleStore(conn))
	runProjectStoreSuite(t, sqlite.NewProjectStore(conn), sqlite.NewGroupStore(conn))
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
	runGroupStoreSuite(t, pg.NewGroupStore(conn), pg.NewUserStore(conn), pg.NewRoleStore(conn))
	runProjectStoreSuite(t, pg.NewProjectStore(conn), pg.NewGroupStore(conn))
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
	runGroupStoreSuite(t, mysql.NewGroupStore(conn), mysql.NewUserStore(conn), mysql.NewRoleStore(conn))
	runProjectStoreSuite(t, mysql.NewProjectStore(conn), mysql.NewGroupStore(conn))
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
	runGroupStoreSuite(t, mssql.NewGroupStore(conn), mssql.NewUserStore(conn), mssql.NewRoleStore(conn))
	runProjectStoreSuite(t, mssql.NewProjectStore(conn), mssql.NewGroupStore(conn))
}

func runGroupStoreSuite(t *testing.T, groupStore stores.GroupStore, userStore stores.UserStore, roleStore stores.RoleStore) {
	ctx := context.Background()

	// Test Groups
	name := "Admin Group"
	desc := "Admins only"
	imageUri := "https://example.com/admin.png"
	group := &models.Group{
		ID:          uuid.New(),
		Name:        name,
		Description: &desc,
		ImageURI:    &imageUri,
		IsActive:    true,
		CreatedAt:   time.Now().Truncate(time.Second), // truncate to match SQL precision
	}

	err := groupStore.Create(ctx, group)
	require.NoError(t, err)

	g, err := groupStore.Get(ctx, group.ID)
	require.NoError(t, err)
	require.NotNil(t, g)
	assert.Equal(t, group.Name, g.Name)
	assert.Equal(t, group.Description, g.Description)
	assert.Equal(t, group.ImageURI, g.ImageURI)

	g2, err := groupStore.GetByName(ctx, "admin group")
	require.NoError(t, err)
	require.NotNil(t, g2)
	assert.Equal(t, group.ID, g2.ID)

	// List groups
	res, err := groupStore.List(ctx, core.ListOptions{Filter: []core.FilterOption{{Field: "name", Operator: "like", Value: "%ADMIN%"}}})
	require.NoError(t, err)
	assert.Equal(t, 1, res.TotalCount)
	assert.Len(t, res.Items, 1)

	// Update Group
	newName := "Super Admins"
	group.Name = newName
	err = groupStore.Update(ctx, group)
	require.NoError(t, err)

	g3, err := groupStore.Get(ctx, group.ID)
	require.NoError(t, err)
	require.NotNil(t, g3)
	assert.Equal(t, newName, g3.Name)

	// Associations Setup
	uEmail := "group_user@example.com"
	uName := "Group User"
	user := &models.User{
		ID:           uuid.New(),
		PrimaryEmail: &uEmail,
		Name:         &uName,
		IsBanned:     false,
		CreatedAt:    time.Now().Truncate(time.Second),
	}
	err = userStore.Create(ctx, user)
	require.NoError(t, err)

	role := &models.Role{
		ID:          uuid.New(),
		Name:        "group_role",
		Description: "Group Role",
	}
	err = roleStore.Create(ctx, role)
	require.NoError(t, err)

	// Users
	err = groupStore.AddUser(ctx, group.ID, user.ID)
	require.NoError(t, err)

	users, err := groupStore.ListUsers(ctx, group.ID)
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, user.ID, users[0].UserID)

	err = groupStore.RemoveUser(ctx, group.ID, user.ID)
	require.NoError(t, err)

	users, err = groupStore.ListUsers(ctx, group.ID)
	require.NoError(t, err)
	assert.Len(t, users, 0)

	// Admins
	err = groupStore.AddAdmin(ctx, group.ID, user.ID)
	require.NoError(t, err)

	admins, err := groupStore.ListAdmins(ctx, group.ID)
	require.NoError(t, err)
	assert.Len(t, admins, 1)
	assert.Equal(t, user.ID, admins[0].UserID)

	err = groupStore.RemoveAdmin(ctx, group.ID, user.ID)
	require.NoError(t, err)

	admins, err = groupStore.ListAdmins(ctx, group.ID)
	require.NoError(t, err)
	assert.Len(t, admins, 0)

	// Roles
	err = groupStore.AddRole(ctx, group.ID, role.ID)
	require.NoError(t, err)

	roles, err := groupStore.ListRoles(ctx, group.ID)
	require.NoError(t, err)
	assert.Len(t, roles, 1)
	assert.Equal(t, role.ID, roles[0].RoleID)

	err = groupStore.RemoveRole(ctx, group.ID, role.ID)
	require.NoError(t, err)

	roles, err = groupStore.ListRoles(ctx, group.ID)
	require.NoError(t, err)
	assert.Len(t, roles, 0)

	err = groupStore.Delete(ctx, group.ID)
	require.NoError(t, err)

	g4, err := groupStore.Get(ctx, group.ID)
	require.NoError(t, err)
	require.Nil(t, g4)
}

func runProjectStoreSuite(t *testing.T, projectStore stores.ProjectStore, groupStore stores.GroupStore) {
	ctx := context.Background()

	// Test Projects
	desc := "My new project"
	project := &models.Project{
		ID:          uuid.New(),
		Name:        "Acme Project",
		Slug:        "acme-project",
		Description: &desc,
		IsActive:    true,
		CreatedAt:   time.Now().Truncate(time.Second),
	}

	err := projectStore.Create(ctx, project)
	require.NoError(t, err)

	p, err := projectStore.Get(ctx, project.ID)
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, project.Name, p.Name)
	assert.Equal(t, project.Slug, p.Slug)
	assert.Equal(t, project.Description, p.Description)

	p2, err := projectStore.GetBySlug(ctx, "acme-project")
	require.NoError(t, err)
	require.NotNil(t, p2)
	assert.Equal(t, project.ID, p2.ID)

	// List projects
	res, err := projectStore.List(ctx, core.ListOptions{Filter: []core.FilterOption{{Field: "name", Operator: "like", Value: "%ACME%"}}})
	require.NoError(t, err)
	assert.Equal(t, 1, res.TotalCount)
	assert.Len(t, res.Items, 1)

	// Update Project
	newSlug := "super-acme"
	project.Slug = newSlug
	err = projectStore.Update(ctx, project)
	require.NoError(t, err)

	p3, err := projectStore.Get(ctx, project.ID)
	require.NoError(t, err)
	require.NotNil(t, p3)
	assert.Equal(t, newSlug, p3.Slug)

	// Associations Setup
	gDesc := "Project Admins"
	group := &models.Group{
		ID:          uuid.New(),
		Name:        "Project Admins",
		Description: &gDesc,
		IsActive:    true,
		CreatedAt:   time.Now().Truncate(time.Second),
	}
	err = groupStore.Create(ctx, group)
	require.NoError(t, err)

	// Project Groups
	err = projectStore.AddGroup(ctx, project.ID, group.ID, 7) // e.g. permission bitmask
	require.NoError(t, err)

	pgs, err := projectStore.ListGroups(ctx, project.ID)
	require.NoError(t, err)
	assert.Len(t, pgs, 1)
	assert.Equal(t, group.ID, pgs[0].GroupID)
	assert.Equal(t, int64(7), pgs[0].Permissions)

	err = projectStore.UpdateGroupPermissions(ctx, project.ID, group.ID, 15)
	require.NoError(t, err)

	pgs, err = projectStore.ListGroups(ctx, project.ID)
	require.NoError(t, err)
	assert.Len(t, pgs, 1)
	assert.Equal(t, int64(15), pgs[0].Permissions)

	err = projectStore.RemoveGroup(ctx, project.ID, group.ID)
	require.NoError(t, err)

	pgs, err = projectStore.ListGroups(ctx, project.ID)
	require.NoError(t, err)
	assert.Len(t, pgs, 0)

	err = projectStore.Delete(ctx, project.ID)
	require.NoError(t, err)

	p4, err := projectStore.Get(ctx, project.ID)
	require.NoError(t, err)
	require.Nil(t, p4)
}
