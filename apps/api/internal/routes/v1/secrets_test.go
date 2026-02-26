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

type MockSecretStore struct {
	ListSecretNamesFunc     func(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.SecretName], error)
	GetSecretNameFunc       func(ctx context.Context, id uuid.UUID) (*models.SecretName, error)
	GetSecretNameByNameFunc func(ctx context.Context, projectID uuid.UUID, name string) (*models.SecretName, error)
	CreateSecretNameFunc    func(ctx context.Context, name *models.SecretName) error
	DeleteSecretNameFunc    func(ctx context.Context, id uuid.UUID) error

	GetSecretFunc    func(ctx context.Context, nameID, envID uuid.UUID) (*models.Secret, error)
	UpsertSecretFunc func(ctx context.Context, secret *models.Secret) error

	ListCertificateNamesFunc     func(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.CertificateName], error)
	GetCertificateNameFunc       func(ctx context.Context, id uuid.UUID) (*models.CertificateName, error)
	GetCertificateNameByNameFunc func(ctx context.Context, projectID uuid.UUID, name string) (*models.CertificateName, error)
	CreateCertificateNameFunc    func(ctx context.Context, name *models.CertificateName) error
	DeleteCertificateNameFunc    func(ctx context.Context, id uuid.UUID) error

	GetCertificateFunc    func(ctx context.Context, nameID, envID uuid.UUID) (*models.Certificate, error)
	UpsertCertificateFunc func(ctx context.Context, cert *models.Certificate) error

	ListSSHKeyNamesFunc     func(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.SSHKeyName], error)
	GetSSHKeyNameFunc       func(ctx context.Context, id uuid.UUID) (*models.SSHKeyName, error)
	GetSSHKeyNameByNameFunc func(ctx context.Context, projectID uuid.UUID, name string) (*models.SSHKeyName, error)
	CreateSSHKeyNameFunc    func(ctx context.Context, name *models.SSHKeyName) error
	DeleteSSHKeyNameFunc    func(ctx context.Context, id uuid.UUID) error

	GetSSHKeyFunc    func(ctx context.Context, nameID, envID uuid.UUID) (*models.SSHKey, error)
	UpsertSSHKeyFunc func(ctx context.Context, key *models.SSHKey) error
}

func (m *MockSecretStore) ListSecretNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.SecretName], error) {
	if m.ListSecretNamesFunc != nil {
		return m.ListSecretNamesFunc(ctx, projectID, opts)
	}
	return core.ListResult[models.SecretName]{}, nil
}
func (m *MockSecretStore) GetSecretName(ctx context.Context, id uuid.UUID) (*models.SecretName, error) {
	return nil, nil
}
func (m *MockSecretStore) GetSecretNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.SecretName, error) {
	return nil, nil
}
func (m *MockSecretStore) CreateSecretName(ctx context.Context, name *models.SecretName) error {
	if m.CreateSecretNameFunc != nil {
		return m.CreateSecretNameFunc(ctx, name)
	}
	return nil
}
func (m *MockSecretStore) DeleteSecretName(ctx context.Context, id uuid.UUID) error { return nil }

func (m *MockSecretStore) GetSecret(ctx context.Context, nameID, envID uuid.UUID) (*models.Secret, error) {
	return nil, nil
}
func (m *MockSecretStore) UpsertSecret(ctx context.Context, secret *models.Secret) error { return nil }

func (m *MockSecretStore) ListCertificateNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.CertificateName], error) {
	return core.ListResult[models.CertificateName]{}, nil
}
func (m *MockSecretStore) GetCertificateName(ctx context.Context, id uuid.UUID) (*models.CertificateName, error) {
	return nil, nil
}
func (m *MockSecretStore) GetCertificateNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.CertificateName, error) {
	return nil, nil
}
func (m *MockSecretStore) CreateCertificateName(ctx context.Context, name *models.CertificateName) error {
	return nil
}
func (m *MockSecretStore) DeleteCertificateName(ctx context.Context, id uuid.UUID) error { return nil }

func (m *MockSecretStore) GetCertificate(ctx context.Context, nameID, envID uuid.UUID) (*models.Certificate, error) {
	return nil, nil
}
func (m *MockSecretStore) UpsertCertificate(ctx context.Context, cert *models.Certificate) error {
	return nil
}

func (m *MockSecretStore) ListSSHKeyNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.SSHKeyName], error) {
	return core.ListResult[models.SSHKeyName]{}, nil
}
func (m *MockSecretStore) GetSSHKeyName(ctx context.Context, id uuid.UUID) (*models.SSHKeyName, error) {
	return nil, nil
}
func (m *MockSecretStore) GetSSHKeyNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.SSHKeyName, error) {
	return nil, nil
}
func (m *MockSecretStore) CreateSSHKeyName(ctx context.Context, name *models.SSHKeyName) error {
	return nil
}
func (m *MockSecretStore) DeleteSSHKeyName(ctx context.Context, id uuid.UUID) error { return nil }

func (m *MockSecretStore) GetSSHKey(ctx context.Context, nameID, envID uuid.UUID) (*models.SSHKey, error) {
	return nil, nil
}
func (m *MockSecretStore) UpsertSSHKey(ctx context.Context, key *models.SSHKey) error { return nil }

func TestListSecretNames(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStore := &MockSecretStore{
		ListSecretNamesFunc: func(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.SecretName], error) {
			return core.ListResult[models.SecretName]{
				Items: []models.SecretName{
					{ID: uuid.New(), Name: "db_password", CreatedAt: time.Now()},
				},
				TotalCount: 1,
			}, nil
		},
	}

	r := gin.Default()
	projGroup := r.Group("/projects/:id")
	envGroup := projGroup.Group("/environments/:envId")
	RegisterSecretRoutes(projGroup, envGroup, mockStore)

	w := httptest.NewRecorder()
	projectID := uuid.New().String()
	req, _ := http.NewRequest("GET", "/projects/"+projectID+"/secrets", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp routes.Response[core.ListResult[models.SecretName]]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
	assert.Equal(t, 1, resp.Value.TotalCount)
}

func TestCreateSecretName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStore := &MockSecretStore{
		CreateSecretNameFunc: func(ctx context.Context, name *models.SecretName) error {
			return nil
		},
	}

	r := gin.Default()
	projGroup := r.Group("/projects/:id")
	envGroup := projGroup.Group("/environments/:envId")
	RegisterSecretRoutes(projGroup, envGroup, mockStore)

	reqBody := CreateNameRequest{
		Name: "api_key",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	projectID := uuid.New().String()
	req, _ := http.NewRequest("POST", "/projects/"+projectID+"/secrets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp routes.Response[models.SecretName]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
	assert.Equal(t, "api_key", resp.Value.Name)
}
