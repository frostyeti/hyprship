package v1

import (
	"net/http"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/middleware"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	store    stores.ProjectStore
	envStore stores.EnvironmentStore
	cfgStore stores.ConfigStore
	secStore stores.SecretStore
}

func RegisterProjectRoutes(r *gin.RouterGroup, store stores.ProjectStore, envStore stores.EnvironmentStore, cfgStore stores.ConfigStore, secStore stores.SecretStore) {
	h := &ProjectHandler{store: store, envStore: envStore, cfgStore: cfgStore, secStore: secStore}

	r.GET("", h.ListProjects)
	r.GET("/:id", h.GetProject)
	r.POST("", h.CreateProject)
	r.PUT("/:id", h.UpdateProject)
	r.DELETE("/:id", h.DeleteProject)

	projectSubGroup := r.Group("/:id")
	projectSubGroup.Use(middleware.ProjectPermissionsMiddleware(store, "id"))

	projectSubGroup.GET("/groups", h.ListProjectGroups)
	projectSubGroup.POST("/groups/:groupId", h.AddProjectGroup)
	projectSubGroup.PUT("/groups/:groupId", h.UpdateProjectGroup)
	projectSubGroup.DELETE("/groups/:groupId", h.RemoveProjectGroup)

	// Environments
	envGroup := projectSubGroup.Group("/environments")
	RegisterEnvironmentRoutes(envGroup, envStore)

	// Configs
	envSubGroup := envGroup.Group("/:envId")
	RegisterConfigRoutes(projectSubGroup, envSubGroup, cfgStore)

	// Secrets
	RegisterSecretRoutes(projectSubGroup, envSubGroup, secStore)
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	opts := core.ListOptions{
		Offset: 0,
		Limit:  50,
	}

	result, err := h.store.List(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "list_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(result))
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	idOrSlug := c.Param("id")

	var project *models.Project
	var err error

	if id, parseErr := uuid.Parse(idOrSlug); parseErr == nil {
		project, err = h.store.Get(c.Request.Context(), id)
	} else {
		project, err = h.store.GetBySlug(c.Request.Context(), idOrSlug)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "get_failed",
			Message: err.Error(),
		}))
		return
	}

	if project == nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "Project not found",
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(project))
}

type CreateProjectRequest struct {
	Name        string  `json:"name" binding:"required"`
	Slug        string  `json:"slug" binding:"required"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"isActive"`
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	project := &models.Project{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if err := h.store.Create(c.Request.Context(), project); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "create_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(project))
}

type UpdateProjectRequest struct {
	Name        string  `json:"name" binding:"required"`
	Slug        string  `json:"slug" binding:"required"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"isActive"`
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_id",
			Message: "Invalid UUID format",
		}))
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	project, err := h.store.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "update_failed",
			Message: err.Error(),
		}))
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "Project not found",
		}))
		return
	}

	project.Name = req.Name
	project.Slug = req.Slug
	project.Description = req.Description
	project.IsActive = req.IsActive

	if err := h.store.Update(c.Request.Context(), project); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "update_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(project))
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_id",
			Message: "Invalid UUID format",
		}))
		return
	}

	if err := h.store.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "delete_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}

func (h *ProjectHandler) ListProjectGroups(c *gin.Context) {
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid UUID format"}))
		return
	}

	groups, err := h.store.ListGroups(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(groups))
}

type ProjectGroupRequest struct {
	Permissions int64 `json:"permissions"`
}

func (h *ProjectHandler) AddProjectGroup(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID format"}))
		return
	}

	groupID, err := uuid.Parse(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid group UUID format"}))
		return
	}

	var req ProjectGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	if err := h.store.AddGroup(c.Request.Context(), projectID, groupID, req.Permissions); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "add_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse[any](nil))
}

func (h *ProjectHandler) UpdateProjectGroup(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID format"}))
		return
	}

	groupID, err := uuid.Parse(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid group UUID format"}))
		return
	}

	var req ProjectGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	if err := h.store.UpdateGroupPermissions(c.Request.Context(), projectID, groupID, req.Permissions); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "update_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}

func (h *ProjectHandler) RemoveProjectGroup(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID format"}))
		return
	}

	groupID, err := uuid.Parse(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid group UUID format"}))
		return
	}

	if err := h.store.RemoveGroup(c.Request.Context(), projectID, groupID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "remove_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}
