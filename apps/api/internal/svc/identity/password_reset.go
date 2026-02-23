package identity

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"
)

func generateSecureToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (s *identityService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userStore.GetByEmail(ctx, email)
	if err != nil || user == nil {
		// Do not leak existence
		return nil
	}

	auth, err := s.authStore.Get(ctx, user.ID)
	if err != nil {
		// Unlikely to hit unless user is fully broken, but treat it similarly.
		return nil
	}

	token := generateSecureToken()
	digest, err := s.hasher.Hash(token)
	if err != nil {
		return err
	}

	auth.OtpDigest = &digest
	expires := time.Now().Add(time.Duration(s.config.Identity.Reset.TokenTTLMinutes) * time.Minute)
	auth.OtpExpiresAt = &expires
	// In reality, auth.OtpLinkToken could just be used. But we typically hash it.

	if err := s.authStore.Update(ctx, auth); err != nil {
		return err
	}

	// TODO: Send email containing the cleartext token
	return nil
}

func (s *identityService) ResetPassword(ctx context.Context, email string, token string, newPassword string) error {
	user, err := s.userStore.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return ErrInvalidCredentials
	}

	auth, err := s.authStore.Get(ctx, user.ID)
	if err != nil {
		return ErrInvalidCredentials
	}

	if auth.OtpDigest == nil || auth.OtpExpiresAt == nil || time.Now().After(*auth.OtpExpiresAt) {
		return ErrInvalidCredentials
	}

	match, err := s.hasher.Verify(*auth.OtpDigest, token)
	if err != nil || !match {
		return ErrInvalidCredentials
	}

	// Consume OTP immediately so it can't be used twice
	auth.OtpDigest = nil
	auth.OtpExpiresAt = nil
	_ = s.authStore.Update(ctx, auth)

	// Now set the new password
	return s.SetPassword(ctx, user.ID, newPassword)
}
