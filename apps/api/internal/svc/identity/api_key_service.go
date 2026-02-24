package identity

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/crypto"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/google/uuid"
)

type APIKeyService interface {
	CreateAPIKey(ctx context.Context, userID uuid.UUID, name string, comment *string, expiresAt *time.Time) (*models.UserAPIKey, string, error)
	RevokeAPIKey(ctx context.Context, userID uuid.UUID, keyID int32) error
	ListAPIKeys(ctx context.Context, userID uuid.UUID) ([]models.UserAPIKey, error)
}

type apiKeyService struct {
	store  stores.UserAPIKeyStore
	hasher crypto.PasswordHasher
}

func NewAPIKeyService(store stores.UserAPIKeyStore, hasher crypto.PasswordHasher) APIKeyService {
	if hasher == nil {
		hasher = crypto.GetPasswordHasher("pbkdf2")
	}
	return &apiKeyService{
		store:  store,
		hasher: hasher,
	}
}

func (s *apiKeyService) CreateAPIKey(ctx context.Context, userID uuid.UUID, name string, comment *string, expiresAt *time.Time) (*models.UserAPIKey, string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, "", err
	}

	rawKey := fmt.Sprintf("hypr_%s", base64.RawURLEncoding.EncodeToString(b))

	digest, err := s.hasher.Hash(rawKey)
	if err != nil {
		return nil, "", err
	}

	hint := fmt.Sprintf("%s...%s", rawKey[:8], rawKey[len(rawKey)-4:])

	apiKey := &models.UserAPIKey{
		UserID:     userID,
		Name:       name,
		NameUpcase: strings.ToUpper(name),
		KeyHint:    hint,
		KeyDigest:  digest,
		Comment:    comment,
		ExpiresAt:  expiresAt,
		CreatedAt:  time.Now(),
		IsLocked:   false,
		IsRevoked:  false,
	}

	err = s.store.Create(ctx, apiKey)
	if err != nil {
		return nil, "", err
	}

	return apiKey, rawKey, nil
}

func (s *apiKeyService) RevokeAPIKey(ctx context.Context, userID uuid.UUID, keyID int32) error {
	key, err := s.store.Get(ctx, keyID)
	if err != nil {
		return err
	}

	if key.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.store.Delete(ctx, keyID)
}

func (s *apiKeyService) ListAPIKeys(ctx context.Context, userID uuid.UUID) ([]models.UserAPIKey, error) {
	return s.store.ListByUser(ctx, userID)
}
