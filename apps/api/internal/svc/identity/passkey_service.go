package identity

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

type PasskeyService interface {
	BeginRegistration(ctx context.Context, user *models.User) (*protocol.CredentialCreation, *webauthn.SessionData, error)
	FinishRegistration(ctx context.Context, user *models.User, sessionData webauthn.SessionData, req *http.Request) (*webauthn.Credential, error)

	BeginLogin(ctx context.Context, user *models.User) (*protocol.CredentialAssertion, *webauthn.SessionData, error)
	FinishLogin(ctx context.Context, user *models.User, sessionData webauthn.SessionData, req *http.Request) (*webauthn.Credential, error)

	ListPasskeys(ctx context.Context, userID uuid.UUID) ([]models.UserPasskey, error)
	DeletePasskey(ctx context.Context, userID uuid.UUID, passkeyID uuid.UUID) error
}

type passkeyService struct {
	store stores.UserPasskeyStore
	wa    *webauthn.WebAuthn
}

type webauthnUser struct {
	*models.User
	credentials []webauthn.Credential
}

func (u *webauthnUser) WebAuthnID() []byte {
	return []byte(u.ID.String())
}

func (u *webauthnUser) WebAuthnName() string {
	if u.PrimaryEmail != nil {
		return *u.PrimaryEmail
	}
	return u.ID.String()
}

func (u *webauthnUser) WebAuthnDisplayName() string {
	if u.Name != nil {
		return *u.Name
	}
	return u.WebAuthnName()
}

func (u *webauthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func NewPasskeyService(store stores.UserPasskeyStore, rpDisplayName, rpID, rpOrigin string) (PasskeyService, error) {
	wa, err := webauthn.New(&webauthn.Config{
		RPDisplayName: rpDisplayName,
		RPID:          rpID,
		RPOrigins:     []string{rpOrigin},
	})
	if err != nil {
		return nil, err
	}

	return &passkeyService{
		store: store,
		wa:    wa,
	}, nil
}

func (s *passkeyService) getUserWithCredentials(ctx context.Context, user *models.User) (*webauthnUser, error) {
	passkeys, err := s.store.ListByUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	var creds []webauthn.Credential
	for _, pk := range passkeys {
		var cred webauthn.Credential
		if err := json.Unmarshal(pk.Data, &cred); err == nil {
			creds = append(creds, cred)
		}
	}

	return &webauthnUser{
		User:        user,
		credentials: creds,
	}, nil
}

func (s *passkeyService) BeginRegistration(ctx context.Context, user *models.User) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	wu, err := s.getUserWithCredentials(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return s.wa.BeginRegistration(wu)
}

func (s *passkeyService) FinishRegistration(ctx context.Context, user *models.User, sessionData webauthn.SessionData, req *http.Request) (*webauthn.Credential, error) {
	wu, err := s.getUserWithCredentials(ctx, user)
	if err != nil {
		return nil, err
	}

	credential, err := s.wa.FinishRegistration(wu, sessionData, req)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(credential)
	if err != nil {
		return nil, err
	}

	passkey := &models.UserPasskey{
		ID:           uuid.New(),
		UserID:       user.ID,
		CredentialID: credential.ID,
		Data:         data,
		CreatedAt:    time.Now(),
	}

	if err := s.store.Create(ctx, passkey); err != nil {
		return nil, err
	}

	return credential, nil
}

func (s *passkeyService) BeginLogin(ctx context.Context, user *models.User) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	wu, err := s.getUserWithCredentials(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return s.wa.BeginLogin(wu)
}

func (s *passkeyService) FinishLogin(ctx context.Context, user *models.User, sessionData webauthn.SessionData, req *http.Request) (*webauthn.Credential, error) {
	wu, err := s.getUserWithCredentials(ctx, user)
	if err != nil {
		return nil, err
	}

	return s.wa.FinishLogin(wu, sessionData, req)
}

func (s *passkeyService) ListPasskeys(ctx context.Context, userID uuid.UUID) ([]models.UserPasskey, error) {
	return s.store.ListByUser(ctx, userID)
}

func (s *passkeyService) DeletePasskey(ctx context.Context, userID uuid.UUID, passkeyID uuid.UUID) error {
	pk, err := s.store.Get(ctx, passkeyID)
	if err != nil {
		return err
	}
	if pk.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.store.Delete(ctx, passkeyID)
}

// SessionDataStore is a simple memory store for webauthn sessions

var sessionStore sync.Map

func SaveSessionData(sessionID string, data webauthn.SessionData) {
	sessionStore.Store(sessionID, data)
	// In production, we should have a background goroutine to clean up expired sessions
	go func() {
		time.Sleep(5 * time.Minute)
		sessionStore.Delete(sessionID)
	}()
}

func GetSessionData(sessionID string) (webauthn.SessionData, bool) {
	if val, ok := sessionStore.Load(sessionID); ok {
		return val.(webauthn.SessionData), true
	}
	return webauthn.SessionData{}, false
}

func DeleteSessionData(sessionID string) {
	sessionStore.Delete(sessionID)
}
