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
	UserStore                UserStore
	RoleStore                RoleStore
	GroupStore               GroupStore
	ProjectStore             ProjectStore
	EnvironmentStore         EnvironmentStore
	ConfigStore              ConfigStore
	SecretStore              SecretStore
	UserPasswordAuthStore    UserPasswordAuthStore
	UserSessionStore         UserSessionStore
	UserPasswordHistoryStore UserPasswordHistoryStore
	UserAPIKeyStore          UserAPIKeyStore
	UserPasskeyStore         UserPasskeyStore
	UserTotpStore            UserTotpStore
	UserEmailStore           UserEmailStore
	UserPhoneStore           UserPhoneStore
}

func NewStoreFactory(driver string, db *sql.DB) (*StoreFactory, error) {
	switch driver {
	case "sqlite", "sqlite3":
		return &StoreFactory{
			UserStore:                sqlite.NewUserStore(db),
			RoleStore:                sqlite.NewRoleStore(db),
			GroupStore:               sqlite.NewGroupStore(db),
			ProjectStore:             sqlite.NewProjectStore(db),
			EnvironmentStore:         sqlite.NewEnvironmentStore(db),
			ConfigStore:              sqlite.NewConfigStore(db),
			SecretStore:              sqlite.NewSecretStore(db),
			UserPasswordAuthStore:    sqlite.NewUserPasswordAuthStore(db),
			UserSessionStore:         sqlite.NewUserSessionStore(db),
			UserPasswordHistoryStore: sqlite.NewUserPasswordHistoryStore(db),
			UserAPIKeyStore:          sqlite.NewUserAPIKeyStore(db),
			UserPasskeyStore:         sqlite.NewUserPasskeyStore(db),
			UserTotpStore:            sqlite.NewUserTotpStore(db),
			UserEmailStore:           sqlite.NewUserEmailStore(db),
			UserPhoneStore:           sqlite.NewUserPhoneStore(db),
		}, nil
	case "pg", "postgres":
		return &StoreFactory{
			UserStore:                pg.NewUserStore(db),
			RoleStore:                pg.NewRoleStore(db),
			GroupStore:               pg.NewGroupStore(db),
			ProjectStore:             pg.NewProjectStore(db),
			EnvironmentStore:         pg.NewEnvironmentStore(db),
			ConfigStore:              pg.NewConfigStore(db),
			SecretStore:              pg.NewSecretStore(db),
			UserPasswordAuthStore:    pg.NewUserPasswordAuthStore(db),
			UserSessionStore:         pg.NewUserSessionStore(db),
			UserPasswordHistoryStore: pg.NewUserPasswordHistoryStore(db),
			UserAPIKeyStore:          pg.NewUserAPIKeyStore(db),
			UserPasskeyStore:         pg.NewUserPasskeyStore(db),
			UserTotpStore:            pg.NewUserTotpStore(db),
			UserEmailStore:           pg.NewUserEmailStore(db),
			UserPhoneStore:           pg.NewUserPhoneStore(db),
		}, nil
	case "mysql", "mariadb":
		return &StoreFactory{
			UserStore:                mysql.NewUserStore(db),
			RoleStore:                mysql.NewRoleStore(db),
			GroupStore:               mysql.NewGroupStore(db),
			ProjectStore:             mysql.NewProjectStore(db),
			EnvironmentStore:         mysql.NewEnvironmentStore(db),
			ConfigStore:              mysql.NewConfigStore(db),
			SecretStore:              mysql.NewSecretStore(db),
			UserPasswordAuthStore:    mysql.NewUserPasswordAuthStore(db),
			UserSessionStore:         mysql.NewUserSessionStore(db),
			UserPasswordHistoryStore: mysql.NewUserPasswordHistoryStore(db),
			UserAPIKeyStore:          mysql.NewUserAPIKeyStore(db),
			UserPasskeyStore:         mysql.NewUserPasskeyStore(db),
			UserTotpStore:            mysql.NewUserTotpStore(db),
			UserEmailStore:           mysql.NewUserEmailStore(db),
			UserPhoneStore:           mysql.NewUserPhoneStore(db),
		}, nil
	case "mssql", "sqlserver":
		return &StoreFactory{
			UserStore:                mssql.NewUserStore(db),
			RoleStore:                mssql.NewRoleStore(db),
			GroupStore:               mssql.NewGroupStore(db),
			ProjectStore:             mssql.NewProjectStore(db),
			EnvironmentStore:         mssql.NewEnvironmentStore(db),
			ConfigStore:              mssql.NewConfigStore(db),
			SecretStore:              mssql.NewSecretStore(db),
			UserPasswordAuthStore:    mssql.NewUserPasswordAuthStore(db),
			UserSessionStore:         mssql.NewUserSessionStore(db),
			UserPasswordHistoryStore: mssql.NewUserPasswordHistoryStore(db),
			UserAPIKeyStore:          mssql.NewUserAPIKeyStore(db),
			UserPasskeyStore:         mssql.NewUserPasskeyStore(db),
			UserTotpStore:            mssql.NewUserTotpStore(db),
			UserEmailStore:           mssql.NewUserEmailStore(db),
			UserPhoneStore:           mssql.NewUserPhoneStore(db),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}
