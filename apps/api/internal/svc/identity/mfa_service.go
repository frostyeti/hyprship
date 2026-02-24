package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

type MfaService interface {
	GenerateTotpSecret(ctx context.Context, user *models.User) (string, string, error)
	VerifyAndEnableTotp(ctx context.Context, userID uuid.UUID, code string) error
	DisableTotp(ctx context.Context, userID uuid.UUID) error
	VerifyTotp(ctx context.Context, userID uuid.UUID, code string) (bool, error)
}

type mfaService struct {
	store     stores.UserTotpStore
	authStore stores.UserPasswordAuthStore
	issuer    string
}

func NewMfaService(store stores.UserTotpStore, authStore stores.UserPasswordAuthStore, issuer string) MfaService {
	return &mfaService{
		store:     store,
		authStore: authStore,
		issuer:    issuer,
	}
}

func (s *mfaService) GenerateTotpSecret(ctx context.Context, user *models.User) (string, string, error) {
	// Generate new key
	accountName := user.ID.String()
	if user.PrimaryEmail != nil {
		accountName = *user.PrimaryEmail
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: accountName,
	})
	if err != nil {
		return "", "", err
	}

	secret := key.Secret()
	url := key.URL()

	// Store it as unverified
	err = s.store.Upsert(ctx, &models.UserTotp{
		UserID:     user.ID,
		Secret:     secret,
		IsVerified: false,
	})
	if err != nil {
		return "", "", err
	}

	return secret, url, nil
}

func (s *mfaService) VerifyAndEnableTotp(ctx context.Context, userID uuid.UUID, code string) error {
	t, err := s.store.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.New("totp not configured")
	}

	valid := totp.Validate(code, t.Secret)
	if !valid {
		return errors.New("invalid code")
	}

	t.IsVerified = true
	if err := s.store.Upsert(ctx, t); err != nil {
		return err
	}

	auth, err := s.authStore.Get(ctx, userID)
	if err != nil {
		return err
	}
	auth.TotpEnabled = true
	return s.authStore.Update(ctx, auth)
}

func (s *mfaService) DisableTotp(ctx context.Context, userID uuid.UUID) error {
	if err := s.store.Delete(ctx, userID); err != nil {
		return err
	}

	auth, err := s.authStore.Get(ctx, userID)
	if err != nil {
		// Ignore if auth record doesn't exist for some reason
		return nil
	}
	auth.TotpEnabled = false
	return s.authStore.Update(ctx, auth)
}

func (s *mfaService) VerifyTotp(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	t, err := s.store.GetByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	if t == nil || !t.IsVerified {
		return false, fmt.Errorf("totp not verified")
	}

	return totp.Validate(code, t.Secret), nil
}
