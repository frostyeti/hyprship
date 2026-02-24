package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockPasskeyService struct {
	BeginRegistrationFunc  func(ctx context.Context, user *models.User) (*protocol.CredentialCreation, *webauthn.SessionData, error)
	FinishRegistrationFunc func(ctx context.Context, user *models.User, sessionData webauthn.SessionData, req *http.Request) (*webauthn.Credential, error)
	BeginLoginFunc         func(ctx context.Context, user *models.User) (*protocol.CredentialAssertion, *webauthn.SessionData, error)
	FinishLoginFunc        func(ctx context.Context, user *models.User, sessionData webauthn.SessionData, req *http.Request) (*webauthn.Credential, error)
	ListPasskeysFunc       func(ctx context.Context, userID uuid.UUID) ([]models.UserPasskey, error)
	DeletePasskeyFunc      func(ctx context.Context, userID uuid.UUID, passkeyID uuid.UUID) error
}

func (m *MockPasskeyService) BeginRegistration(ctx context.Context, user *models.User) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	if m.BeginRegistrationFunc != nil {
		return m.BeginRegistrationFunc(ctx, user)
	}
	return nil, nil, nil
}

func (m *MockPasskeyService) FinishRegistration(ctx context.Context, user *models.User, sessionData webauthn.SessionData, req *http.Request) (*webauthn.Credential, error) {
	if m.FinishRegistrationFunc != nil {
		return m.FinishRegistrationFunc(ctx, user, sessionData, req)
	}
	return nil, nil
}

func (m *MockPasskeyService) BeginLogin(ctx context.Context, user *models.User) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	if m.BeginLoginFunc != nil {
		return m.BeginLoginFunc(ctx, user)
	}
	return nil, nil, nil
}

func (m *MockPasskeyService) FinishLogin(ctx context.Context, user *models.User, sessionData webauthn.SessionData, req *http.Request) (*webauthn.Credential, error) {
	if m.FinishLoginFunc != nil {
		return m.FinishLoginFunc(ctx, user, sessionData, req)
	}
	return nil, nil
}

func (m *MockPasskeyService) ListPasskeys(ctx context.Context, userID uuid.UUID) ([]models.UserPasskey, error) {
	if m.ListPasskeysFunc != nil {
		return m.ListPasskeysFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockPasskeyService) DeletePasskey(ctx context.Context, userID uuid.UUID, passkeyID uuid.UUID) error {
	if m.DeletePasskeyFunc != nil {
		return m.DeletePasskeyFunc(ctx, userID, passkeyID)
	}
	return nil
}

type MockUserStoreForPasskey struct {
	stores.UserStore
	GetFunc        func(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmailFunc func(ctx context.Context, email string) (*models.User, error)
}

func (m *MockUserStoreForPasskey) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserStoreForPasskey) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func setupPasskeyRouter(mockSvc *MockPasskeyService, mockStore *MockUserStoreForPasskey) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := r.Group("/api/v1")
	usersGroup := api.Group("/users")
	RegisterUserRoutes(usersGroup, nil, mockStore, nil, nil, mockSvc, nil)

	authGroup := api.Group("/auth")
	RegisterAuthRoutes(authGroup, nil, mockStore, mockSvc, nil)

	return r
}

func TestListMyPasskeys(t *testing.T) {
	mockSvc := &MockPasskeyService{}
	mockStore := &MockUserStoreForPasskey{}
	r := setupPasskeyRouter(mockSvc, mockStore)

	userID := uuid.New()
	passkeyID := uuid.New()

	mockSvc.ListPasskeysFunc = func(ctx context.Context, u uuid.UUID) ([]models.UserPasskey, error) {
		assert.Equal(t, userID, u)
		return []models.UserPasskey{
			{
				ID:        passkeyID,
				UserID:    userID,
				CreatedAt: time.Now(),
			},
		}, nil
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/me/passkeys", nil)
	req.Header.Set("Authorization", "Bearer "+userID.String())
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteMyPasskey(t *testing.T) {
	mockSvc := &MockPasskeyService{}
	mockStore := &MockUserStoreForPasskey{}
	r := setupPasskeyRouter(mockSvc, mockStore)

	userID := uuid.New()
	passkeyID := uuid.New()

	mockSvc.DeletePasskeyFunc = func(ctx context.Context, u uuid.UUID, pid uuid.UUID) error {
		assert.Equal(t, userID, u)
		assert.Equal(t, passkeyID, pid)
		return nil
	}

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/me/passkeys/"+passkeyID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+userID.String())
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
