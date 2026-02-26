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

func TestListProjects(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStore := &MockProjectStore{
		ListFunc: func(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Project], error) {
			return core.ListResult[models.Project]{
				Items: []models.Project{
					{ID: uuid.New(), Name: "Project 1", Slug: "project-1"},
				},
				TotalCount: 1,
			}, nil
		},
	}

	r := gin.Default()
	RegisterProjectRoutes(r.Group("/projects"), mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp routes.Response[core.ListResult[models.Project]]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
	assert.Equal(t, 1, resp.Value.TotalCount)
}

func TestGetProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	projID := uuid.New()

	mockStore := &MockProjectStore{
		GetFunc: func(ctx context.Context, id uuid.UUID) (*models.Project, error) {
			return &models.Project{ID: projID, Name: "Project 1"}, nil
		},
	}

	r := gin.Default()
	RegisterProjectRoutes(r.Group("/projects"), mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects/"+projID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp routes.Response[models.Project]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
	assert.Equal(t, projID, resp.Value.ID)
	assert.Equal(t, "Project 1", resp.Value.Name)
}

func TestCreateProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStore := &MockProjectStore{
		CreateFunc: func(ctx context.Context, proj *models.Project) error {
			return nil
		},
	}

	r := gin.Default()
	RegisterProjectRoutes(r.Group("/projects"), mockStore)

	reqBody := CreateProjectRequest{
		Name:     "New Project",
		Slug:     "new-project",
		IsActive: true,
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp routes.Response[models.Project]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
	assert.Equal(t, "New Project", resp.Value.Name)
}
