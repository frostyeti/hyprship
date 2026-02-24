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

type MockAPIKeyService struct {
	CreateAPIKeyFunc func(ctx context.Context, userID uuid.UUID, name string, comment *string, expiresAt *time.Time) (*models.UserAPIKey, string, error)
	RevokeAPIKeyFunc func(ctx context.Context, userID uuid.UUID, keyID int32) error
	ListAPIKeysFunc  func(ctx context.Context, userID uuid.UUID) ([]models.UserAPIKey, error)
}

func (m *MockAPIKeyService) CreateAPIKey(ctx context.Context, userID uuid.UUID, name string, comment *string, expiresAt *time.Time) (*models.UserAPIKey, string, error) {
	if m.CreateAPIKeyFunc != nil {
		return m.CreateAPIKeyFunc(ctx, userID, name, comment, expiresAt)
	}
	return nil, "", nil
}

func (m *MockAPIKeyService) RevokeAPIKey(ctx context.Context, userID uuid.UUID, keyID int32) error {
	if m.RevokeAPIKeyFunc != nil {
		return m.RevokeAPIKeyFunc(ctx, userID, keyID)
	}
	return nil
}

func (m *MockAPIKeyService) ListAPIKeys(ctx context.Context, userID uuid.UUID) ([]models.UserAPIKey, error) {
	if m.ListAPIKeysFunc != nil {
		return m.ListAPIKeysFunc(ctx, userID)
	}
	return nil, nil
}

func setupAPIKeysRouter(mockSvc *MockAPIKeyService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	usersGroup := r.Group("/api/v1/users")
	RegisterUserRoutes(usersGroup, nil, nil, nil, mockSvc, nil, nil, nil, nil)
	return r
}

func TestListMyAPIKeys(t *testing.T) {
	userID := uuid.New()
	mockKey := models.UserAPIKey{
		ID:        1,
		UserID:    userID,
		Name:      "Test Key",
		KeyHint:   "hypr_123...456",
		CreatedAt: time.Now(),
	}

	mockSvc := &MockAPIKeyService{
		ListAPIKeysFunc: func(ctx context.Context, reqUserID uuid.UUID) ([]models.UserAPIKey, error) {
			assert.Equal(t, userID, reqUserID)
			return []models.UserAPIKey{mockKey}, nil
		},
	}

	r := setupAPIKeysRouter(mockSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/me/api-keys", nil)
	req.Header.Set("Authorization", "Bearer "+userID.String())
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response routes.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)

	keys := response.Value.([]interface{})
	assert.Len(t, keys, 1)
	keyMap := keys[0].(map[string]interface{})
	assert.Equal(t, "Test Key", keyMap["name"])
	assert.Equal(t, "hypr_123...456", keyMap["keyHint"])
}

func TestCreateMyAPIKey(t *testing.T) {
	userID := uuid.New()
	mockKey := &models.UserAPIKey{
		ID:        1,
		UserID:    userID,
		Name:      "New Key",
		KeyHint:   "hypr_abc...def",
		CreatedAt: time.Now(),
	}

	mockSvc := &MockAPIKeyService{
		CreateAPIKeyFunc: func(ctx context.Context, reqUserID uuid.UUID, name string, comment *string, expiresAt *time.Time) (*models.UserAPIKey, string, error) {
			assert.Equal(t, userID, reqUserID)
			assert.Equal(t, "New Key", name)
			return mockKey, "hypr_abcdef1234567890", nil
		},
	}

	r := setupAPIKeysRouter(mockSvc)

	reqBody := CreateAPIKeyRequest{
		Name: "New Key",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users/me/api-keys", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userID.String())
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response routes.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)

	valMap := response.Value.(map[string]interface{})
	assert.Equal(t, "hypr_abcdef1234567890", valMap["secret"])

	keyMap := valMap["apiKey"].(map[string]interface{})
	assert.Equal(t, "New Key", keyMap["name"])
}

func TestDeleteMyAPIKey(t *testing.T) {
	userID := uuid.New()

	mockSvc := &MockAPIKeyService{
		RevokeAPIKeyFunc: func(ctx context.Context, reqUserID uuid.UUID, keyID int32) error {
			assert.Equal(t, userID, reqUserID)
			assert.Equal(t, int32(1), keyID)
			return nil
		},
	}

	r := setupAPIKeysRouter(mockSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/me/api-keys/1", nil)
	req.Header.Set("Authorization", "Bearer "+userID.String())
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response routes.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)
	assert.Equal(t, "ok", response.Value)
}
