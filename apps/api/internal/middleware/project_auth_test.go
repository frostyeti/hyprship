package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/middleware"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockProjectStore struct {
	projectGroups []models.ProjectGroup
}

func (m *mockProjectStore) ListGroups(ctx context.Context, projectID uuid.UUID) ([]models.ProjectGroup, error) {
	return m.projectGroups, nil
}
func (m *mockProjectStore) GetBySlug(ctx context.Context, slug string) (*models.Project, error) {
	return &models.Project{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")}, nil
}
func (m *mockProjectStore) List(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Project], error) {
	return core.ListResult[models.Project]{}, nil
}
func (m *mockProjectStore) Get(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return nil, nil
}
func (m *mockProjectStore) Create(ctx context.Context, project *models.Project) error { return nil }
func (m *mockProjectStore) Update(ctx context.Context, project *models.Project) error { return nil }
func (m *mockProjectStore) Delete(ctx context.Context, id uuid.UUID) error            { return nil }
func (m *mockProjectStore) AddGroup(ctx context.Context, projectID, groupID uuid.UUID, permissions int64) error {
	return nil
}
func (m *mockProjectStore) UpdateGroupPermissions(ctx context.Context, projectID, groupID uuid.UUID, permissions int64) error {
	return nil
}
func (m *mockProjectStore) RemoveGroup(ctx context.Context, projectID, groupID uuid.UUID) error {
	return nil
}

func TestProjectPermissionsMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	projectID := uuid.New()
	groupID1 := uuid.New()
	groupID2 := uuid.New()

	store := &mockProjectStore{
		projectGroups: []models.ProjectGroup{
			{ProjectID: projectID, GroupID: groupID1, Permissions: 1}, // Bit 0
			{ProjectID: projectID, GroupID: groupID2, Permissions: 4}, // Bit 2
		},
	}

	tests := []struct {
		name           string
		userGroups     []string
		expectedPerms  int64
		requirePerms   int64
		expectedStatus int
	}{
		{
			name:           "No groups",
			userGroups:     nil,
			expectedPerms:  0,
			requirePerms:   1,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "In Group 1 only",
			userGroups:     []string{groupID1.String()},
			expectedPerms:  1,
			requirePerms:   1,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "In Group 2 only",
			userGroups:     []string{groupID2.String()},
			expectedPerms:  4,
			requirePerms:   4,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "In both groups",
			userGroups:     []string{groupID1.String(), groupID2.String()},
			expectedPerms:  5, // 1 | 4 = 5
			requirePerms:   5,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Lacking required perms",
			userGroups:     []string{groupID1.String()},
			expectedPerms:  1,
			requirePerms:   4,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()

			// Mock ClaimsMiddleware setting GroupsKey
			r.Use(func(c *gin.Context) {
				if tt.userGroups != nil {
					c.Set(middleware.GroupsKey, tt.userGroups)
				}
				c.Next()
			})

			r.Use(middleware.ProjectPermissionsMiddleware(store, "id"))
			r.Use(middleware.RequireProjectPermission(tt.requirePerms))

			r.GET("/projects/:id", func(c *gin.Context) {
				perms, _ := c.Get(middleware.ProjectPermissionsKey)
				assert.Equal(t, tt.expectedPerms, perms)
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/projects/"+projectID.String(), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
