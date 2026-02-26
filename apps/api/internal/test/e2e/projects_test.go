//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"context"
	"testing"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	v1 "github.com/frostyeti/hyprship/apps/api/internal/routes/v1"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Re-using MockProjectStore from v1 tests but making a fresh one for the E2E simulation.
type MockProjectStore struct {
	ListFunc      func(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Project], error)
	GetFunc       func(ctx context.Context, id uuid.UUID) (*models.Project, error)
	GetBySlugFunc func(ctx context.Context, slug string) (*models.Project, error)
	CreateFunc    func(ctx context.Context, project *models.Project) error
	UpdateFunc    func(ctx context.Context, project *models.Project) error
	DeleteFunc    func(ctx context.Context, id uuid.UUID) error

	AddGroupFunc               func(ctx context.Context, projectID, groupID uuid.UUID, permissions int64) error
	RemoveGroupFunc            func(ctx context.Context, projectID, groupID uuid.UUID) error
	UpdateGroupPermissionsFunc func(ctx context.Context, projectID, groupID uuid.UUID, permissions int64) error
	ListGroupsFunc             func(ctx context.Context, projectID uuid.UUID) ([]models.ProjectGroup, error)
}

func (m *MockProjectStore) List(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Project], error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, opts)
	}
	return core.ListResult[models.Project]{}, nil
}
func (m *MockProjectStore) Get(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}
func (m *MockProjectStore) GetBySlug(ctx context.Context, slug string) (*models.Project, error) {
	if m.GetBySlugFunc != nil {
		return m.GetBySlugFunc(ctx, slug)
	}
	return nil, nil
}
func (m *MockProjectStore) Create(ctx context.Context, project *models.Project) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, project)
	}
	return nil
}
func (m *MockProjectStore) Update(ctx context.Context, project *models.Project) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, project)
	}
	return nil
}
func (m *MockProjectStore) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
func (m *MockProjectStore) AddGroup(ctx context.Context, projectID, groupID uuid.UUID, permissions int64) error {
	if m.AddGroupFunc != nil {
		return m.AddGroupFunc(ctx, projectID, groupID, permissions)
	}
	return nil
}
func (m *MockProjectStore) RemoveGroup(ctx context.Context, projectID, groupID uuid.UUID) error {
	if m.RemoveGroupFunc != nil {
		return m.RemoveGroupFunc(ctx, projectID, groupID)
	}
	return nil
}
func (m *MockProjectStore) UpdateGroupPermissions(ctx context.Context, projectID, groupID uuid.UUID, permissions int64) error {
	if m.UpdateGroupPermissionsFunc != nil {
		return m.UpdateGroupPermissionsFunc(ctx, projectID, groupID, permissions)
	}
	return nil
}
func (m *MockProjectStore) ListGroups(ctx context.Context, projectID uuid.UUID) ([]models.ProjectGroup, error) {
	if m.ListGroupsFunc != nil {
		return m.ListGroupsFunc(ctx, projectID)
	}
	return nil, nil
}

func TestE2E_Projects_CRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	apiV1 := r.Group("/api/v1")

	pID := uuid.New()

	mockStore := &MockProjectStore{
		GetFunc: func(ctx context.Context, id uuid.UUID) (*models.Project, error) {
			return &models.Project{
				ID:   pID,
				Name: "E2E Project",
				Slug: "e2e-project",
			}, nil
		},
		CreateFunc: func(ctx context.Context, proj *models.Project) error {
			proj.ID = pID
			return nil
		},
	}

	// We pass only the mock project store
	v1.RegisterRoutes(apiV1, nil, nil, nil, nil, mockStore, nil, nil, nil, nil, nil, nil)

	// 1. Create a project
	reqBody := v1.CreateProjectRequest{
		Name:     "E2E Project",
		Slug:     "e2e-project",
		IsActive: true,
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp routes.Response[models.Project]
	err := json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.NoError(t, err)
	assert.True(t, createResp.Ok)
	assert.Equal(t, "E2E Project", createResp.Value.Name)
	assert.NotEqual(t, uuid.Nil, createResp.Value.ID)

	// 2. Get the project we just created
	wGet := httptest.NewRecorder()
	reqGet, _ := http.NewRequest("GET", "/api/v1/projects/"+pID.String(), nil)
	r.ServeHTTP(wGet, reqGet)

	assert.Equal(t, http.StatusOK, wGet.Code)

	var getResp routes.Response[models.Project]
	err = json.Unmarshal(wGet.Body.Bytes(), &getResp)
	assert.NoError(t, err)
	assert.True(t, getResp.Ok)
	assert.Equal(t, "E2E Project", getResp.Value.Name)
	assert.Equal(t, pID, getResp.Value.ID)
}
