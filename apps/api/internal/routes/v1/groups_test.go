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

type MockGroupStore struct {
	ListFunc      func(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Group], error)
	GetFunc       func(ctx context.Context, id uuid.UUID) (*models.Group, error)
	GetByNameFunc func(ctx context.Context, name string) (*models.Group, error)
	CreateFunc    func(ctx context.Context, group *models.Group) error
	UpdateFunc    func(ctx context.Context, group *models.Group) error
	DeleteFunc    func(ctx context.Context, id uuid.UUID) error

	AddUserFunc     func(ctx context.Context, groupID, userID uuid.UUID) error
	RemoveUserFunc  func(ctx context.Context, groupID, userID uuid.UUID) error
	ListUsersFunc   func(ctx context.Context, groupID uuid.UUID) ([]models.GroupUser, error)
	AddAdminFunc    func(ctx context.Context, groupID, userID uuid.UUID) error
	RemoveAdminFunc func(ctx context.Context, groupID, userID uuid.UUID) error
	ListAdminsFunc  func(ctx context.Context, groupID uuid.UUID) ([]models.GroupAdmin, error)
	AddRoleFunc     func(ctx context.Context, groupID, roleID uuid.UUID) error
	RemoveRoleFunc  func(ctx context.Context, groupID, roleID uuid.UUID) error
	ListRolesFunc   func(ctx context.Context, groupID uuid.UUID) ([]models.GroupRole, error)
}

func (m *MockGroupStore) List(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Group], error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, opts)
	}
	return core.ListResult[models.Group]{}, nil
}
func (m *MockGroupStore) Get(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}
func (m *MockGroupStore) GetByName(ctx context.Context, name string) (*models.Group, error) {
	if m.GetByNameFunc != nil {
		return m.GetByNameFunc(ctx, name)
	}
	return nil, nil
}
func (m *MockGroupStore) Create(ctx context.Context, group *models.Group) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, group)
	}
	return nil
}
func (m *MockGroupStore) Update(ctx context.Context, group *models.Group) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, group)
	}
	return nil
}
func (m *MockGroupStore) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockGroupStore) AddUser(ctx context.Context, groupID, userID uuid.UUID) error {
	if m.AddUserFunc != nil {
		return m.AddUserFunc(ctx, groupID, userID)
	}
	return nil
}

func (m *MockGroupStore) RemoveUser(ctx context.Context, groupID, userID uuid.UUID) error {
	if m.RemoveUserFunc != nil {
		return m.RemoveUserFunc(ctx, groupID, userID)
	}
	return nil
}

func (m *MockGroupStore) ListUsers(ctx context.Context, groupID uuid.UUID) ([]models.GroupUser, error) {
	if m.ListUsersFunc != nil {
		return m.ListUsersFunc(ctx, groupID)
	}
	return nil, nil
}

func (m *MockGroupStore) AddAdmin(ctx context.Context, groupID, userID uuid.UUID) error {
	if m.AddAdminFunc != nil {
		return m.AddAdminFunc(ctx, groupID, userID)
	}
	return nil
}

func (m *MockGroupStore) RemoveAdmin(ctx context.Context, groupID, userID uuid.UUID) error {
	if m.RemoveAdminFunc != nil {
		return m.RemoveAdminFunc(ctx, groupID, userID)
	}
	return nil
}

func (m *MockGroupStore) ListAdmins(ctx context.Context, groupID uuid.UUID) ([]models.GroupAdmin, error) {
	if m.ListAdminsFunc != nil {
		return m.ListAdminsFunc(ctx, groupID)
	}
	return nil, nil
}

func (m *MockGroupStore) AddRole(ctx context.Context, groupID, roleID uuid.UUID) error {
	if m.AddRoleFunc != nil {
		return m.AddRoleFunc(ctx, groupID, roleID)
	}
	return nil
}

func (m *MockGroupStore) RemoveRole(ctx context.Context, groupID, roleID uuid.UUID) error {
	if m.RemoveRoleFunc != nil {
		return m.RemoveRoleFunc(ctx, groupID, roleID)
	}
	return nil
}

func (m *MockGroupStore) ListRoles(ctx context.Context, groupID uuid.UUID) ([]models.GroupRole, error) {
	if m.ListRolesFunc != nil {
		return m.ListRolesFunc(ctx, groupID)
	}
	return nil, nil
}

func setupGroupTestRouter(store *MockGroupStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	RegisterGroupRoutes(r.Group("/groups"), store)
	return r
}

func TestGroupHandler_ListGroups(t *testing.T) {
	mockStore := &MockGroupStore{
		ListFunc: func(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Group], error) {
			return core.ListResult[models.Group]{
				TotalCount: 1,
				Items: []models.Group{
					{ID: uuid.New(), Name: "Test Group"},
				},
			}, nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/groups", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response routes.Response[core.ListResult[models.Group]]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)
	assert.Equal(t, 1, response.Value.TotalCount)
	assert.Len(t, response.Value.Items, 1)
}

func TestGroupHandler_CreateGroup(t *testing.T) {
	mockStore := &MockGroupStore{
		CreateFunc: func(ctx context.Context, group *models.Group) error {
			assert.Equal(t, "New Group", group.Name)
			return nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	reqBody := CreateGroupRequest{
		Name:     "New Group",
		IsActive: true,
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/groups", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response routes.Response[models.Group]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)
	assert.Equal(t, "New Group", response.Value.Name)
}

func TestGroupHandler_GetGroup(t *testing.T) {
	groupID := uuid.New()
	mockStore := &MockGroupStore{
		GetFunc: func(ctx context.Context, id uuid.UUID) (*models.Group, error) {
			assert.Equal(t, groupID, id)
			return &models.Group{ID: id, Name: "Found Group"}, nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/groups/"+groupID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response routes.Response[models.Group]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)
	assert.Equal(t, "Found Group", response.Value.Name)
}

func TestGroupHandler_UpdateGroup(t *testing.T) {
	groupID := uuid.New()
	mockStore := &MockGroupStore{
		GetFunc: func(ctx context.Context, id uuid.UUID) (*models.Group, error) {
			return &models.Group{ID: id, Name: "Old Name"}, nil
		},
		UpdateFunc: func(ctx context.Context, group *models.Group) error {
			assert.Equal(t, "Updated Name", group.Name)
			return nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	reqBody := UpdateGroupRequest{
		Name:     "Updated Name",
		IsActive: true,
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/groups/"+groupID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response routes.Response[models.Group]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)
	assert.Equal(t, "Updated Name", response.Value.Name)
}

func TestGroupHandler_DeleteGroup(t *testing.T) {
	groupID := uuid.New()
	mockStore := &MockGroupStore{
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			assert.Equal(t, groupID, id)
			return nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/groups/"+groupID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response routes.Response[any]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)
}

func TestGroupHandler_AddUser(t *testing.T) {
	groupID := uuid.New()
	userID := uuid.New()
	mockStore := &MockGroupStore{
		AddUserFunc: func(ctx context.Context, gID, uID uuid.UUID) error {
			assert.Equal(t, groupID, gID)
			assert.Equal(t, userID, uID)
			return nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/groups/"+groupID.String()+"/users/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGroupHandler_RemoveUser(t *testing.T) {
	groupID := uuid.New()
	userID := uuid.New()
	mockStore := &MockGroupStore{
		RemoveUserFunc: func(ctx context.Context, gID, uID uuid.UUID) error {
			assert.Equal(t, groupID, gID)
			assert.Equal(t, userID, uID)
			return nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/groups/"+groupID.String()+"/users/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGroupHandler_ListUsers(t *testing.T) {
	groupID := uuid.New()
	mockStore := &MockGroupStore{
		ListUsersFunc: func(ctx context.Context, gID uuid.UUID) ([]models.GroupUser, error) {
			assert.Equal(t, groupID, gID)
			return []models.GroupUser{{GroupID: gID, UserID: uuid.New()}}, nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/groups/"+groupID.String()+"/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGroupHandler_AddAdmin(t *testing.T) {
	groupID := uuid.New()
	userID := uuid.New()
	mockStore := &MockGroupStore{
		AddAdminFunc: func(ctx context.Context, gID, uID uuid.UUID) error {
			assert.Equal(t, groupID, gID)
			assert.Equal(t, userID, uID)
			return nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/groups/"+groupID.String()+"/admins/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGroupHandler_RemoveAdmin(t *testing.T) {
	groupID := uuid.New()
	userID := uuid.New()
	mockStore := &MockGroupStore{
		RemoveAdminFunc: func(ctx context.Context, gID, uID uuid.UUID) error {
			assert.Equal(t, groupID, gID)
			assert.Equal(t, userID, uID)
			return nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/groups/"+groupID.String()+"/admins/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGroupHandler_ListAdmins(t *testing.T) {
	groupID := uuid.New()
	mockStore := &MockGroupStore{
		ListAdminsFunc: func(ctx context.Context, gID uuid.UUID) ([]models.GroupAdmin, error) {
			assert.Equal(t, groupID, gID)
			return []models.GroupAdmin{{GroupID: gID, UserID: uuid.New()}}, nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/groups/"+groupID.String()+"/admins", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGroupHandler_AddRole(t *testing.T) {
	groupID := uuid.New()
	roleID := uuid.New()
	mockStore := &MockGroupStore{
		AddRoleFunc: func(ctx context.Context, gID, rID uuid.UUID) error {
			assert.Equal(t, groupID, gID)
			assert.Equal(t, roleID, rID)
			return nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/groups/"+groupID.String()+"/roles/"+roleID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGroupHandler_RemoveRole(t *testing.T) {
	groupID := uuid.New()
	roleID := uuid.New()
	mockStore := &MockGroupStore{
		RemoveRoleFunc: func(ctx context.Context, gID, rID uuid.UUID) error {
			assert.Equal(t, groupID, gID)
			assert.Equal(t, roleID, rID)
			return nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/groups/"+groupID.String()+"/roles/"+roleID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGroupHandler_ListRoles(t *testing.T) {
	groupID := uuid.New()
	mockStore := &MockGroupStore{
		ListRolesFunc: func(ctx context.Context, gID uuid.UUID) ([]models.GroupRole, error) {
			assert.Equal(t, groupID, gID)
			return []models.GroupRole{{GroupID: gID, RoleID: uuid.New()}}, nil
		},
	}
	router := setupGroupTestRouter(mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/groups/"+groupID.String()+"/roles", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
