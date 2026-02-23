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

	// Importer/Exporter implementations
	core.Exporter[models.User]
	core.Importer[models.User]
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

	// Importer/Exporter implementations
	core.Exporter[models.Role]
	core.Importer[models.Role]
}
