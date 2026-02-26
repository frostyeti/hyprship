package v1

import (
	"net/http"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EnvironmentHandler struct {
	store stores.EnvironmentStore
}

func RegisterEnvironmentRoutes(r *gin.RouterGroup, store stores.EnvironmentStore) {
	h := &EnvironmentHandler{store: store}

	r.GET("", h.ListEnvironments)
	r.GET("/:envId", h.GetEnvironment)
	r.POST("", h.CreateEnvironment)
	r.PUT("/:envId", h.UpdateEnvironment)
	r.DELETE("/:envId", h.DeleteEnvironment)
}

func (h *EnvironmentHandler) ListEnvironments(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	opts := core.ListOptions{
		Offset: 0,
		Limit:  50,
	}

	result, err := h.store.List(c.Request.Context(), projectID, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(result))
}

func (h *EnvironmentHandler) GetEnvironment(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	idOrName := c.Param("envId")
	var env *models.Environment

	if id, parseErr := uuid.Parse(idOrName); parseErr == nil {
		env, err = h.store.Get(c.Request.Context(), id)
	} else {
		env, err = h.store.GetByName(c.Request.Context(), projectID, idOrName)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "get_failed", Message: err.Error()}))
		return
	}

	if env == nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{Code: "not_found", Message: "Environment not found"}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(env))
}

type CreateEnvironmentRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
}

func (h *EnvironmentHandler) CreateEnvironment(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	var req CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	env := &models.Environment{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := h.store.Create(c.Request.Context(), env); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "create_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(env))
}

type UpdateEnvironmentRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
}

func (h *EnvironmentHandler) UpdateEnvironment(c *gin.Context) {
	envID, err := uuid.Parse(c.Param("envId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid environment UUID format"}))
		return
	}

	var req UpdateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	env, err := h.store.Get(c.Request.Context(), envID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "update_failed", Message: err.Error()}))
		return
	}
	if env == nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{Code: "not_found", Message: "Environment not found"}))
		return
	}

	env.Name = req.Name
	env.Description = req.Description

	if err := h.store.Update(c.Request.Context(), env); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "update_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(env))
}

func (h *EnvironmentHandler) DeleteEnvironment(c *gin.Context) {
	envID, err := uuid.Parse(c.Param("envId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid environment UUID format"}))
		return
	}

	if err := h.store.Delete(c.Request.Context(), envID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "delete_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}
