package stores

import (
	"context"

	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserTotpStore interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserTotp, error)
	Upsert(ctx context.Context, totp *models.UserTotp) error
	Delete(ctx context.Context, userID uuid.UUID) error
}
