package stores

import (
	"context"

	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserEmailStore interface {
	Get(ctx context.Context, id uuid.UUID) (*models.UserEmail, error)
	GetByEmail(ctx context.Context, email string) (*models.UserEmail, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserEmail, error)
	Create(ctx context.Context, email *models.UserEmail) error
	Update(ctx context.Context, email *models.UserEmail) error
	Delete(ctx context.Context, id uuid.UUID) error
}
