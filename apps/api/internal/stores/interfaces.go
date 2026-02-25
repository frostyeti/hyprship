package stores

import (
	"context"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserStore interface {
	// Basic CRUD
	List(ctx context.Context, opts core.ListOptions) (core.ListResult[models.User], error)
	Get(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Claims Management
	ListClaims(ctx context.Context, userID uuid.UUID) ([]models.UserClaim, error)
	AddClaim(ctx context.Context, claim *models.UserClaim) error
	RemoveClaim(ctx context.Context, claimID int32) error
}

type UserPasswordAuthStore interface {
	Get(ctx context.Context, userID uuid.UUID) (*models.UserPasswordAuth, error)
	Create(ctx context.Context, auth *models.UserPasswordAuth) error
	Update(ctx context.Context, auth *models.UserPasswordAuth) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type UserPasswordHistoryStore interface {
	ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]string, error)
	Create(ctx context.Context, userID uuid.UUID, passwordDigest string) error
	DeleteOld(ctx context.Context, userID uuid.UUID, keepCount int) error
}

type UserSessionStore interface {
	Get(ctx context.Context, sessionID uuid.UUID) (*models.UserSession, error)
	GetByToken(ctx context.Context, token string) (*models.UserSession, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserSession, error)
	Create(ctx context.Context, session *models.UserSession) error
	Update(ctx context.Context, session *models.UserSession) error
	Delete(ctx context.Context, sessionID uuid.UUID) error
	DeleteByUser(ctx context.Context, userID uuid.UUID) error
}

type RoleStore interface {
	// Basic CRUD
	List(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Role], error)
	Get(ctx context.Context, id uuid.UUID) (*models.Role, error)
	Create(ctx context.Context, role *models.Role) error
	Update(ctx context.Context, role *models.Role) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Claims Management
	ListClaims(ctx context.Context, roleID uuid.UUID) ([]models.RoleClaim, error)
	AddClaim(ctx context.Context, claim *models.RoleClaim) error
	RemoveClaim(ctx context.Context, claimID int32) error
}

type GroupStore interface {
	List(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Group], error)
	Get(ctx context.Context, id uuid.UUID) (*models.Group, error)
	GetByName(ctx context.Context, name string) (*models.Group, error)
	Create(ctx context.Context, group *models.Group) error
	Update(ctx context.Context, group *models.Group) error
	Delete(ctx context.Context, id uuid.UUID) error

	AddUser(ctx context.Context, groupID, userID uuid.UUID) error
	RemoveUser(ctx context.Context, groupID, userID uuid.UUID) error
	ListUsers(ctx context.Context, groupID uuid.UUID) ([]models.GroupUser, error)
	ListGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]models.GroupUser, error)

	AddAdmin(ctx context.Context, groupID, userID uuid.UUID) error
	RemoveAdmin(ctx context.Context, groupID, userID uuid.UUID) error
	ListAdmins(ctx context.Context, groupID uuid.UUID) ([]models.GroupAdmin, error)

	AddRole(ctx context.Context, groupID, roleID uuid.UUID) error
	RemoveRole(ctx context.Context, groupID, roleID uuid.UUID) error
	ListRoles(ctx context.Context, groupID uuid.UUID) ([]models.GroupRole, error)
}
