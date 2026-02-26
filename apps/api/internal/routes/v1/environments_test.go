package v1

import (
	"bytes"
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

type MockEnvironmentStore struct {
	ListFunc      func(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.Environment], error)
	GetFunc       func(ctx context.Context, id uuid.UUID) (*models.Environment, error)
	GetByNameFunc func(ctx context.Context, projectID uuid.UUID, name string) (*models.Environment, error)
	CreateFunc    func(ctx context.Context, env *models.Environment) error
	UpdateFunc    func(ctx context.Context, env *models.Environment) error
	DeleteFunc    func(ctx context.Context, id uuid.UUID) error
}

func (m *MockEnvironmentStore) List(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.Environment], error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, projectID, opts)
	}
	return core.ListResult[models.Environment]{}, nil
}
func (m *MockEnvironmentStore) Get(ctx context.Context, id uuid.UUID) (*models.Environment, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}
func (m *MockEnvironmentStore) GetByName(ctx context.Context, projectID uuid.UUID, name string) (*models.Environment, error) {
	if m.GetByNameFunc != nil {
		return m.GetByNameFunc(ctx, projectID, name)
	}
	return nil, nil
}
func (m *MockEnvironmentStore) Create(ctx context.Context, env *models.Environment) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, env)
	}
	return nil
}
func (m *MockEnvironmentStore) Update(ctx context.Context, env *models.Environment) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, env)
	}
	return nil
}
func (m *MockEnvironmentStore) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestListEnvironments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStore := &MockEnvironmentStore{
		ListFunc: func(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.Environment], error) {
			return core.ListResult[models.Environment]{
				Items: []models.Environment{
					{ID: uuid.New(), Name: "Production"},
				},
				TotalCount: 1,
			}, nil
		},
	}

	r := gin.Default()
	RegisterEnvironmentRoutes(r.Group("/projects/:id/environments"), mockStore)

	w := httptest.NewRecorder()
	projectID := uuid.New().String()
	req, _ := http.NewRequest("GET", "/projects/"+projectID+"/environments", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp routes.Response[core.ListResult[models.Environment]]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
	assert.Equal(t, 1, resp.Value.TotalCount)
}

func TestCreateEnvironment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStore := &MockEnvironmentStore{
		CreateFunc: func(ctx context.Context, env *models.Environment) error {
			return nil
		},
	}

	r := gin.Default()
	RegisterEnvironmentRoutes(r.Group("/projects/:id/environments"), mockStore)

	reqBody := CreateEnvironmentRequest{
		Name: "Staging",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	projectID := uuid.New().String()
	req, _ := http.NewRequest("POST", "/projects/"+projectID+"/environments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp routes.Response[models.Environment]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
	assert.Equal(t, "Staging", resp.Value.Name)
}
