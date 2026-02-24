package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockUserStore struct {
	GetFunc    func(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateFunc func(ctx context.Context, user *models.User) error
}

func (m *MockUserStore) List(ctx context.Context, opts core.ListOptions) (core.ListResult[models.User], error) {
	return core.ListResult[models.User]{}, nil
}
func (m *MockUserStore) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}
func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *MockUserStore) Create(ctx context.Context, user *models.User) error { return nil }
func (m *MockUserStore) Update(ctx context.Context, user *models.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}
func (m *MockUserStore) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockUserStore) ListClaims(ctx context.Context, userID uuid.UUID) ([]models.UserClaim, error) {
	return nil, nil
}
func (m *MockUserStore) AddClaim(ctx context.Context, claim *models.UserClaim) error { return nil }
func (m *MockUserStore) RemoveClaim(ctx context.Context, claimID int32) error        { return nil }
func (m *MockUserStore) Export(opts core.ListOptions) ([]models.User, error)         { return nil, nil }
func (m *MockUserStore) Import(items []models.User) error                            { return nil }

func TestGetMe_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockID := uuid.New()
	email := "test@example.com"
	mockUser := &models.User{
		ID:           mockID,
		PrimaryEmail: &email,
	}

	mockStore := &MockUserStore{
		GetFunc: func(ctx context.Context, id uuid.UUID) (*models.User, error) {
			assert.Equal(t, mockID, id)
			return mockUser, nil
		},
	}

	RegisterUserRoutes(r.Group("/api/v1/users"), nil, mockStore, nil, nil, nil, nil, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
	// mock auth middleware
	req.Header.Set("Authorization", "Bearer "+mockID.String())

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response routes.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)

	valMap := response.Value.(map[string]interface{})
	assert.Equal(t, mockID.String(), valMap["id"])
	assert.Equal(t, email, valMap["primaryEmail"])
}

func TestGetMe_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	RegisterUserRoutes(r.Group("/api/v1/users"), nil, &MockUserStore{}, nil, nil, nil, nil, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
	// no auth header

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
