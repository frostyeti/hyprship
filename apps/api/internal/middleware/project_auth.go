package middleware

import (
	"net/http"

	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ProjectPermissionsKey = "projectPermissions"
)

// ProjectPermissionsMiddleware calculates the highest bitwise permission level
// for the authenticated user based on their group associations to the specified project.
// It assumes ClaimsMiddleware has already run and populated the user's group IDs.
func ProjectPermissionsMiddleware(projectStore stores.ProjectStore, projectIdParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user's groups from context (populated by ClaimsMiddleware)
		groupsVal, exists := c.Get(GroupsKey)
		if !exists {
			// Not authenticated or no groups, permissions = 0
			c.Set(ProjectPermissionsKey, int64(0))
			c.Next()
			return
		}

		userGroups, ok := groupsVal.([]string)
		if !ok || len(userGroups) == 0 {
			c.Set(ProjectPermissionsKey, int64(0))
			c.Next()
			return
		}

		// Parse Project ID from path parameter
		projectIDStr := c.Param(projectIdParam)
		if projectIDStr == "" {
			c.Set(ProjectPermissionsKey, int64(0))
			c.Next()
			return
		}

		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			// Try looking up project by slug if it fails UUID parse
			project, slugErr := projectStore.GetBySlug(c.Request.Context(), projectIDStr)
			if slugErr != nil || project == nil {
				c.Set(ProjectPermissionsKey, int64(0))
				c.Next()
				return
			}
			projectID = project.ID
		}

		// Fetch the groups mapped to this project
		projectGroups, err := projectStore.ListGroups(c.Request.Context(), projectID)
		if err != nil {
			c.Set(ProjectPermissionsKey, int64(0))
			c.Next()
			return
		}

		// Calculate bitwise permissions
		var highestPerms int64 = 0
		userGroupMap := make(map[string]bool)
		for _, gID := range userGroups {
			userGroupMap[gID] = true
		}

		for _, pg := range projectGroups {
			if userGroupMap[pg.GroupID.String()] {
				highestPerms |= pg.Permissions
			}
		}

		c.Set(ProjectPermissionsKey, highestPerms)
		c.Next()
	}
}

// RequireProjectPermission ensures the user has specific permission bits for the project.
func RequireProjectPermission(requiredPerms int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		permsVal, exists := c.Get(ProjectPermissionsKey)
		if !exists {
			c.JSON(http.StatusForbidden, routes.ErrorResponse(&routes.ApiError{
				Code:    "forbidden",
				Message: "Missing project permissions",
			}))
			c.Abort()
			return
		}

		perms, ok := permsVal.(int64)
		if !ok {
			c.JSON(http.StatusForbidden, routes.ErrorResponse(&routes.ApiError{
				Code:    "forbidden",
				Message: "Missing project permissions",
			}))
			c.Abort()
			return
		}

		if (perms & requiredPerms) != requiredPerms {
			c.JSON(http.StatusForbidden, routes.ErrorResponse(&routes.ApiError{
				Code:    "forbidden",
				Message: "Insufficient project permissions",
			}))
			c.Abort()
			return
		}

		c.Next()
	}
}
