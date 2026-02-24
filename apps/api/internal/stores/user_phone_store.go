package stores

import (
	"context"

	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserPhoneStore interface {
	Get(ctx context.Context, id uuid.UUID) (*models.UserPhone, error)
	GetByPhone(ctx context.Context, phone string) (*models.UserPhone, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserPhone, error)
	Create(ctx context.Context, phone *models.UserPhone) error
	Update(ctx context.Context, phone *models.UserPhone) error
	Delete(ctx context.Context, id uuid.UUID) error
}
