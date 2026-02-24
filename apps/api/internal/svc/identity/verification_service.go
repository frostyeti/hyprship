package identity

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"
	"strconv"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/google/uuid"
)

var (
	ErrEmailNotFound      = errors.New("email not found")
	ErrPhoneNotFound      = errors.New("phone not found")
	ErrInvalidVerifyToken = errors.New("invalid verification token")
)

type VerificationService interface {
	RequestEmailVerification(ctx context.Context, userID uuid.UUID, email string) error
	ConfirmEmail(ctx context.Context, userID uuid.UUID, email string, token string) error

	RequestPhoneVerification(ctx context.Context, userID uuid.UUID, phone string) error
	ConfirmPhone(ctx context.Context, userID uuid.UUID, phone string, token string) error
}

type verificationService struct {
	emailStore stores.UserEmailStore
	phoneStore stores.UserPhoneStore
	authStore  stores.UserPasswordAuthStore
}

func NewVerificationService(emailStore stores.UserEmailStore, phoneStore stores.UserPhoneStore, authStore stores.UserPasswordAuthStore) VerificationService {
	return &verificationService{
		emailStore: emailStore,
		phoneStore: phoneStore,
		authStore:  authStore,
	}
}

func (s *verificationService) RequestEmailVerification(ctx context.Context, userID uuid.UUID, email string) error {
	userEmail, err := s.emailStore.GetByEmail(ctx, email)
	if err != nil || userEmail == nil {
		return ErrEmailNotFound
	}
	if userEmail.UserID != userID {
		return ErrEmailNotFound
	}
	if userEmail.IsVerified {
		return nil // Already verified
	}

	auth, err := s.authStore.Get(ctx, userID)
	if err != nil || auth == nil {
		return err
	}

	// Generate a secure token (e.g. 32 bytes base64)
	b := make([]byte, 32)
	rand.Read(b)
	token := base64.URLEncoding.EncodeToString(b)

	expires := time.Now().Add(24 * time.Hour)
	auth.OtpLinkToken = &token
	auth.OtpExpiresAt = &expires

	if err := s.authStore.Update(ctx, auth); err != nil {
		return err
	}

	// TODO: Dispatch email via messaging service
	// For now, it silently updates the DB token.
	return nil
}

func (s *verificationService) ConfirmEmail(ctx context.Context, userID uuid.UUID, email string, token string) error {
	userEmail, err := s.emailStore.GetByEmail(ctx, email)
	if err != nil || userEmail == nil || userEmail.UserID != userID {
		return ErrEmailNotFound
	}

	auth, err := s.authStore.Get(ctx, userID)
	if err != nil || auth == nil {
		return err
	}

	if auth.OtpLinkToken == nil || *auth.OtpLinkToken != token {
		return ErrInvalidVerifyToken
	}
	if auth.OtpExpiresAt == nil || time.Now().After(*auth.OtpExpiresAt) {
		return ErrInvalidVerifyToken
	}

	// Token matches, clear it
	auth.OtpLinkToken = nil
	auth.OtpExpiresAt = nil
	if err := s.authStore.Update(ctx, auth); err != nil {
		return err
	}

	// Mark email as verified
	userEmail.IsVerified = true
	return s.emailStore.Update(ctx, userEmail)
}

func (s *verificationService) RequestPhoneVerification(ctx context.Context, userID uuid.UUID, phone string) error {
	userPhone, err := s.phoneStore.GetByPhone(ctx, phone)
	if err != nil || userPhone == nil {
		return ErrPhoneNotFound
	}
	if userPhone.UserID != userID {
		return ErrPhoneNotFound
	}
	if userPhone.IsVerified {
		return nil
	}

	auth, err := s.authStore.Get(ctx, userID)
	if err != nil || auth == nil {
		return err
	}

	// Generate a 6 digit code
	max := big.NewInt(999999)
	n, _ := rand.Int(rand.Reader, max)
	code := strconv.FormatInt(n.Int64(), 10)
	for len(code) < 6 {
		code = "0" + code
	}

	expires := time.Now().Add(10 * time.Minute)
	auth.OtpLinkToken = &code // we reuse OtpLinkToken for the SMS cleartext code
	auth.OtpExpiresAt = &expires

	if err := s.authStore.Update(ctx, auth); err != nil {
		return err
	}

	// TODO: Dispatch SMS via messaging service
	return nil
}

func (s *verificationService) ConfirmPhone(ctx context.Context, userID uuid.UUID, phone string, token string) error {
	userPhone, err := s.phoneStore.GetByPhone(ctx, phone)
	if err != nil || userPhone == nil || userPhone.UserID != userID {
		return ErrPhoneNotFound
	}

	auth, err := s.authStore.Get(ctx, userID)
	if err != nil || auth == nil {
		return err
	}

	if auth.OtpLinkToken == nil || *auth.OtpLinkToken != token {
		return ErrInvalidVerifyToken
	}
	if auth.OtpExpiresAt == nil || time.Now().After(*auth.OtpExpiresAt) {
		return ErrInvalidVerifyToken
	}

	// Code matches, clear it
	auth.OtpLinkToken = nil
	auth.OtpExpiresAt = nil
	if err := s.authStore.Update(ctx, auth); err != nil {
		return err
	}

	// Mark phone as verified
	userPhone.IsVerified = true
	return s.phoneStore.Update(ctx, userPhone)
}
