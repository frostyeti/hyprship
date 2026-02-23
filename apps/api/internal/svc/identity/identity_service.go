package identity

import (
	"context"
	"errors"
	"time"

	"github.com/frostyeti/hyprship/apps/api/config"
	"github.com/frostyeti/hyprship/apps/api/internal/crypto"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrAccountLocked         = errors.New("account is locked")
	ErrPasswordExpired       = errors.New("password has expired")
	ErrPasswordHistoryReuse  = errors.New("password was used recently")
	ErrPasswordPolicyFailed  = errors.New("password does not meet policy requirements")
	ErrMaxConcurrentSessions = errors.New("maximum concurrent sessions reached")
)

type IdentityService interface {
	LoginWithPassword(ctx context.Context, email string, password string, ipAddress *string, userAgent *string) (*models.UserSession, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string) error
	// SetPassword bypasses the old password check, usually used for resets or admin actions
	SetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
}

type identityService struct {
	config       *config.Config
	userStore    stores.UserStore
	authStore    stores.UserPasswordAuthStore
	sessionStore stores.UserSessionStore
	historyStore stores.UserPasswordHistoryStore
	hasher       crypto.PasswordHasher
}

func NewIdentityService(
	cfg *config.Config,
	userStore stores.UserStore,
	authStore stores.UserPasswordAuthStore,
	sessionStore stores.UserSessionStore,
	historyStore stores.UserPasswordHistoryStore,
	hasher crypto.PasswordHasher,
) IdentityService {
	if hasher == nil {
		hasher = crypto.GetPasswordHasher("pbkdf2")
	}
	return &identityService{
		config:       cfg,
		userStore:    userStore,
		authStore:    authStore,
		sessionStore: sessionStore,
		historyStore: historyStore,
		hasher:       hasher,
	}
}

func (s *identityService) validatePasswordPolicy(password string) error {
	cfg := s.config.Identity.Password
	if len(password) < cfg.MinLength || len(password) > cfg.MaxLength {
		return ErrPasswordPolicyFailed
	}
	// TODO: Add upper, lower, number, symbol checks based on config
	return nil
}

func (s *identityService) LoginWithPassword(ctx context.Context, email string, password string, ipAddress *string, userAgent *string) (*models.UserSession, error) {
	user, err := s.userStore.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials // Generic error to prevent user enumeration
	}

	auth, err := s.authStore.Get(ctx, user.ID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check Lockout
	if auth.IsLocked {
		return nil, ErrAccountLocked
	}
	if auth.AttemptCount >= s.config.Identity.Lockout.MaxAttempts {
		if auth.LastAttemptedAt != nil {
			lockoutDuration := time.Duration(s.config.Identity.Lockout.DurationMinutes) * time.Minute
			if time.Since(*auth.LastAttemptedAt) < lockoutDuration {
				return nil, ErrAccountLocked
			}
			// Lockout expired, reset attempt count (we do this below if successful, or we just proceed)
		}
	}

	// Verify password
	match, err := s.hasher.Verify(auth.PasswordDigest, password)
	if err != nil || !match {
		// Increment attempts
		now := time.Now()
		auth.LastAttemptedAt = &now
		auth.AttemptCount++
		if auth.AttemptCount >= s.config.Identity.Lockout.MaxAttempts {
			auth.IsLocked = true
		}
		if err := s.authStore.Update(ctx, auth); err != nil {
			// Log error in a real app, returning invalid credentials anyway
			return nil, ErrInvalidCredentials
		}
		return nil, ErrInvalidCredentials
	}

	// Password valid, check expiration
	if auth.PasswordExpiresAt != nil && time.Now().After(*auth.PasswordExpiresAt) {
		return nil, ErrPasswordExpired
	}

	// Reset attempts on success
	auth.AttemptCount = 0
	auth.IsLocked = false
	if err := s.authStore.Update(ctx, auth); err != nil {
		return nil, err
	}

	// Session management
	if s.config.Identity.Sessions.MaxConcurrent > 0 {
		sessions, err := s.sessionStore.ListByUser(ctx, user.ID)
		if err == nil && len(sessions) >= s.config.Identity.Sessions.MaxConcurrent {
			return nil, ErrMaxConcurrentSessions
			// Alternatively, could delete oldest session
		}
	}

	// Create Session
	session := &models.UserSession{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     uuid.New().String(), // Placeholder for token generation, possibly JWT
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(time.Duration(s.config.Identity.Sessions.TTLMinutes) * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := s.sessionStore.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *identityService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionStore.Delete(ctx, sessionID)
}

func (s *identityService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string) error {
	auth, err := s.authStore.Get(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	match, err := s.hasher.Verify(auth.PasswordDigest, oldPassword)
	if err != nil || !match {
		return ErrInvalidCredentials
	}

	return s.SetPassword(ctx, userID, newPassword)
}

func (s *identityService) SetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	if err := s.validatePasswordPolicy(newPassword); err != nil {
		return err
	}

	historyCount := s.config.Identity.Password.HistoryCount
	if historyCount > 0 && s.historyStore != nil {
		history, err := s.historyStore.ListByUser(ctx, userID, historyCount)
		if err == nil {
			for _, hash := range history {
				match, err := s.hasher.Verify(hash, newPassword)
				if err == nil && match {
					return ErrPasswordHistoryReuse
				}
			}
		}
	}

	digest, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	auth, err := s.authStore.Get(ctx, userID)
	isNew := false
	if err != nil {
		// Doesn't exist, create new
		isNew = true
		auth = &models.UserPasswordAuth{
			UserID:    userID,
			CreatedAt: time.Now(),
		}
	}

	auth.PasswordDigest = digest
	if s.config.Identity.Password.ExpireDays > 0 {
		expires := time.Now().AddDate(0, 0, s.config.Identity.Password.ExpireDays)
		auth.PasswordExpiresAt = &expires
	} else {
		auth.PasswordExpiresAt = nil
	}

	auth.AttemptCount = 0
	auth.IsLocked = false
	auth.LastAttemptedAt = nil

	if isNew {
		err = s.authStore.Create(ctx, auth)
	} else {
		now := time.Now()
		auth.UpdatedAt = &now
		err = s.authStore.Update(ctx, auth)
	}

	if err == nil && s.historyStore != nil && historyCount > 0 {
		_ = s.historyStore.Create(ctx, userID, digest)
		_ = s.historyStore.DeleteOld(ctx, userID, historyCount)
	}

	return err
}
