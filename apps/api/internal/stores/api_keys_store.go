package stores

import (
	"context"

	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserAPIKeyStore interface {
	Get(ctx context.Context, id int32) (*models.UserAPIKey, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserAPIKey, error)
	Create(ctx context.Context, apiKey *models.UserAPIKey) error
	Delete(ctx context.Context, id int32) error
}
