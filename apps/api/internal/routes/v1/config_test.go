package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockConfigStore struct {
	ListConfigFileNamesFunc     func(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.ConfigFileName], error)
	GetConfigFileNameFunc       func(ctx context.Context, id uuid.UUID) (*models.ConfigFileName, error)
	GetConfigFileNameByNameFunc func(ctx context.Context, projectID uuid.UUID, name string) (*models.ConfigFileName, error)
	CreateConfigFileNameFunc    func(ctx context.Context, name *models.ConfigFileName) error
	DeleteConfigFileNameFunc    func(ctx context.Context, id uuid.UUID) error

	GetConfigFileFunc    func(ctx context.Context, nameID, envID uuid.UUID) (*models.ConfigFile, error)
	UpsertConfigFileFunc func(ctx context.Context, file *models.ConfigFile) error

	ListEnvVariableNamesFunc     func(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.EnvVariableName], error)
	GetEnvVariableNameFunc       func(ctx context.Context, id uuid.UUID) (*models.EnvVariableName, error)
	GetEnvVariableNameByNameFunc func(ctx context.Context, projectID uuid.UUID, name string) (*models.EnvVariableName, error)
	CreateEnvVariableNameFunc    func(ctx context.Context, name *models.EnvVariableName) error
	DeleteEnvVariableNameFunc    func(ctx context.Context, id uuid.UUID) error

	GetEnvVariableFunc    func(ctx context.Context, nameID, envID uuid.UUID) (*models.EnvVariable, error)
	UpsertEnvVariableFunc func(ctx context.Context, v *models.EnvVariable) error
}

func (m *MockConfigStore) ListConfigFileNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.ConfigFileName], error) {
	if m.ListConfigFileNamesFunc != nil {
		return m.ListConfigFileNamesFunc(ctx, projectID, opts)
	}
	return core.ListResult[models.ConfigFileName]{}, nil
}
func (m *MockConfigStore) GetConfigFileName(ctx context.Context, id uuid.UUID) (*models.ConfigFileName, error) {
	return nil, nil
}
func (m *MockConfigStore) GetConfigFileNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.ConfigFileName, error) {
	return nil, nil
}
func (m *MockConfigStore) CreateConfigFileName(ctx context.Context, name *models.ConfigFileName) error {
	if m.CreateConfigFileNameFunc != nil {
		return m.CreateConfigFileNameFunc(ctx, name)
	}
	return nil
}
func (m *MockConfigStore) DeleteConfigFileName(ctx context.Context, id uuid.UUID) error { return nil }

func (m *MockConfigStore) GetConfigFile(ctx context.Context, nameID, envID uuid.UUID) (*models.ConfigFile, error) {
	return nil, nil
}
func (m *MockConfigStore) UpsertConfigFile(ctx context.Context, file *models.ConfigFile) error {
	return nil
}

func (m *MockConfigStore) ListEnvVariableNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.EnvVariableName], error) {
	return core.ListResult[models.EnvVariableName]{}, nil
}
func (m *MockConfigStore) GetEnvVariableName(ctx context.Context, id uuid.UUID) (*models.EnvVariableName, error) {
	return nil, nil
}
func (m *MockConfigStore) GetEnvVariableNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.EnvVariableName, error) {
	return nil, nil
}
func (m *MockConfigStore) CreateEnvVariableName(ctx context.Context, name *models.EnvVariableName) error {
	return nil
}
func (m *MockConfigStore) DeleteEnvVariableName(ctx context.Context, id uuid.UUID) error { return nil }

func (m *MockConfigStore) GetEnvVariable(ctx context.Context, nameID, envID uuid.UUID) (*models.EnvVariable, error) {
	return nil, nil
}
func (m *MockConfigStore) UpsertEnvVariable(ctx context.Context, v *models.EnvVariable) error {
	return nil
}

func TestListConfigFileNames(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStore := &MockConfigStore{
		ListConfigFileNamesFunc: func(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.ConfigFileName], error) {
			return core.ListResult[models.ConfigFileName]{
				Items: []models.ConfigFileName{
					{ID: uuid.New(), Name: "app.json", CreatedAt: time.Now()},
				},
				TotalCount: 1,
			}, nil
		},
	}

	r := gin.Default()
	projGroup := r.Group("/projects/:id")
	envGroup := projGroup.Group("/environments/:envId")
	RegisterConfigRoutes(projGroup, envGroup, mockStore)

	w := httptest.NewRecorder()
	projectID := uuid.New().String()
	req, _ := http.NewRequest("GET", "/projects/"+projectID+"/config-files", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp routes.Response[core.ListResult[models.ConfigFileName]]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
	assert.Equal(t, 1, resp.Value.TotalCount)
}

func TestCreateConfigFileName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStore := &MockConfigStore{
		CreateConfigFileNameFunc: func(ctx context.Context, name *models.ConfigFileName) error {
			return nil
		},
	}

	r := gin.Default()
	projGroup := r.Group("/projects/:id")
	envGroup := projGroup.Group("/environments/:envId")
	RegisterConfigRoutes(projGroup, envGroup, mockStore)

	reqBody := CreateConfigNameRequest{
		Name: "settings.yaml",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	projectID := uuid.New().String()
	req, _ := http.NewRequest("POST", "/projects/"+projectID+"/config-files", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp routes.Response[models.ConfigFileName]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
	assert.Equal(t, "settings.yaml", resp.Value.Name)
}
