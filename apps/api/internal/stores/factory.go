package stores

import (
	"database/sql"
	"fmt"

	"github.com/frostyeti/hyprship/apps/api/internal/stores/mssql"
	"github.com/frostyeti/hyprship/apps/api/internal/stores/mysql"
	"github.com/frostyeti/hyprship/apps/api/internal/stores/pg"
	"github.com/frostyeti/hyprship/apps/api/internal/stores/sqlite"
)

type StoreFactory struct {
	UserStore              UserStore
	RoleStore              RoleStore
	UserPasswordAuthStore  UserPasswordAuthStore
	UserSessionStore       UserSessionStore
	UserPasswordHistoryStore UserPasswordHistoryStore
}

func NewStoreFactory(driver string, db *sql.DB) (*StoreFactory, error) {
	switch driver {
	case "sqlite", "sqlite3":
		return &StoreFactory{
			UserStore:              sqlite.NewUserStore(db),
			RoleStore:              sqlite.NewRoleStore(db),
			UserPasswordAuthStore:  sqlite.NewUserPasswordAuthStore(db),
			UserSessionStore:       sqlite.NewUserSessionStore(db),
			UserPasswordHistoryStore: sqlite.NewUserPasswordHistoryStore(db),
		}, nil
	case "pg", "postgres":
		return &StoreFactory{
			UserStore:              pg.NewUserStore(db),
			RoleStore:              pg.NewRoleStore(db),
			UserPasswordAuthStore:  pg.NewUserPasswordAuthStore(db),
			UserSessionStore:       pg.NewUserSessionStore(db),
			UserPasswordHistoryStore: pg.NewUserPasswordHistoryStore(db),
		}, nil
	case "mysql", "mariadb":
		return &StoreFactory{
			UserStore:              mysql.NewUserStore(db),
			RoleStore:              mysql.NewRoleStore(db),
			UserPasswordAuthStore:  mysql.NewUserPasswordAuthStore(db),
			UserSessionStore:       mysql.NewUserSessionStore(db),
			UserPasswordHistoryStore: mysql.NewUserPasswordHistoryStore(db),
		}, nil
	case "mssql", "sqlserver":
		return &StoreFactory{
			UserStore:              mssql.NewUserStore(db),
			RoleStore:              mssql.NewRoleStore(db),
			UserPasswordAuthStore:  mssql.NewUserPasswordAuthStore(db),
			UserSessionStore:       mssql.NewUserSessionStore(db),
			UserPasswordHistoryStore: mssql.NewUserPasswordHistoryStore(db),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}
