package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockIdentityService struct {
	LoginFunc          func(ctx context.Context, email string, password string, ipAddress *string, userAgent *string) (*models.UserSession, error)
	LogoutFunc         func(ctx context.Context, sessionID uuid.UUID) error
	ForgotPasswordFunc func(ctx context.Context, email string) error
	ResetPasswordFunc  func(ctx context.Context, email string, token string, newPassword string) error
}

func (m *MockIdentityService) LoginWithPassword(ctx context.Context, email string, password string, ipAddress *string, userAgent *string) (*models.UserSession, *string, error) {
	if m.LoginFunc != nil {
		sess, err := m.LoginFunc(ctx, email, password, ipAddress, userAgent)
		return sess, nil, err
	}
	return nil, nil, nil
}

func (m *MockIdentityService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	if m.LogoutFunc != nil {
		return m.LogoutFunc(ctx, sessionID)
	}
	return nil
}

func (m *MockIdentityService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string) error {
	return nil
}

func (m *MockIdentityService) ForgotPassword(ctx context.Context, email string) error {
	if m.ForgotPasswordFunc != nil {
		return m.ForgotPasswordFunc(ctx, email)
	}
	return nil
}

func (m *MockIdentityService) ResetPassword(ctx context.Context, email string, token string, newPassword string) error {
	if m.ResetPasswordFunc != nil {
		return m.ResetPasswordFunc(ctx, email, token, newPassword)
	}
	return nil
}

func (m *MockIdentityService) SetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	return nil
}

func TestAuthLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockSession := &models.UserSession{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "test-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	mockSvc := &MockIdentityService{
		LoginFunc: func(ctx context.Context, email, password string, ip, ua *string) (*models.UserSession, error) {
			assert.Equal(t, "test@test.com", email)
			assert.Equal(t, "Password123!", password)
			return mockSession, nil
		},
	}

	RegisterAuthRoutes(r.Group("/api/v1/auth"), mockSvc, nil, nil, nil)

	body, _ := json.Marshal(LoginRequest{
		Email:    "test@test.com",
		Password: "Password123!",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response routes.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)

	valMap := response.Value.(map[string]interface{})
	assert.Equal(t, mockSession.ID.String(), valMap["id"])
}

func TestAuthLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockSvc := &MockIdentityService{
		LogoutFunc: func(ctx context.Context, sessionID uuid.UUID) error {
			return nil
		},
	}

	RegisterAuthRoutes(r.Group("/api/v1/auth"), mockSvc, nil, nil, nil)

	body, _ := json.Marshal(LogoutRequest{
		SessionID: uuid.New().String(),
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response routes.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)
	assert.Equal(t, "ok", response.Value)
}

func (m *MockIdentityService) CreateSession(ctx context.Context, userID uuid.UUID, ipAddress *string, userAgent *string) (*models.UserSession, error) {
	return nil, nil
}
