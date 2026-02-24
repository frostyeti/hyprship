package stores

import (
	"context"

	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserPasskeyStore interface {
	Get(ctx context.Context, id uuid.UUID) (*models.UserPasskey, error)
	GetByCredentialID(ctx context.Context, credentialID []byte) (*models.UserPasskey, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserPasskey, error)
	Create(ctx context.Context, passkey *models.UserPasskey) error
	Delete(ctx context.Context, id uuid.UUID) error
}
