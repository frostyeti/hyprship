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

type ConfigHandler struct {
	store stores.ConfigStore
}

func RegisterConfigRoutes(projectGroup *gin.RouterGroup, envGroup *gin.RouterGroup, store stores.ConfigStore) {
	h := &ConfigHandler{store: store}

	// Project level names
	projectGroup.GET("/config-files", h.ListConfigFileNames)
	projectGroup.POST("/config-files", h.CreateConfigFileName)
	projectGroup.DELETE("/config-files/:nameId", h.DeleteConfigFileName)

	projectGroup.GET("/env-variables", h.ListEnvVariableNames)
	projectGroup.POST("/env-variables", h.CreateEnvVariableName)
	projectGroup.DELETE("/env-variables/:nameId", h.DeleteEnvVariableName)

	// Environment level values
	envGroup.GET("/config-files/:nameId", h.GetConfigFile)
	envGroup.PUT("/config-files/:nameId", h.UpsertConfigFile)

	envGroup.GET("/env-variables/:nameId", h.GetEnvVariable)
	envGroup.PUT("/env-variables/:nameId", h.UpsertEnvVariable)
}

func (h *ConfigHandler) ListConfigFileNames(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	opts := core.ListOptions{Offset: 0, Limit: 50}
	result, err := h.store.ListConfigFileNames(c.Request.Context(), projectID, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(result))
}

type CreateConfigNameRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *ConfigHandler) CreateConfigFileName(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	var req CreateConfigNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	cfgName := &models.ConfigFileName{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateConfigFileName(c.Request.Context(), cfgName); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "create_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(cfgName))
}

func (h *ConfigHandler) DeleteConfigFileName(c *gin.Context) {
	nameID, err := uuid.Parse(c.Param("nameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid name UUID format"}))
		return
	}

	if err := h.store.DeleteConfigFileName(c.Request.Context(), nameID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "delete_failed", Message: err.Error()}))
		return
	}
	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}

func (h *ConfigHandler) GetConfigFile(c *gin.Context) {
	nameID, err := uuid.Parse(c.Param("nameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid name UUID"}))
		return
	}
	envID, err := uuid.Parse(c.Param("envId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid env UUID"}))
		return
	}

	file, err := h.store.GetConfigFile(c.Request.Context(), nameID, envID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "get_failed", Message: err.Error()}))
		return
	}
	if file == nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{Code: "not_found", Message: "Config file not found"}))
		return
	}
	c.JSON(http.StatusOK, routes.SuccessResponse(file))
}

type UpsertConfigFileRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *ConfigHandler) UpsertConfigFile(c *gin.Context) {
	nameID, err := uuid.Parse(c.Param("nameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid name UUID"}))
		return
	}
	envID, err := uuid.Parse(c.Param("envId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid env UUID"}))
		return
	}

	var req UpsertConfigFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	now := time.Now()
	file := &models.ConfigFile{
		ID:               uuid.New(),
		ConfigFileNameID: nameID,
		EnvironmentID:    envID,
		Content:          req.Content,
		CreatedAt:        now,
		UpdatedAt:        &now,
	}

	if err := h.store.UpsertConfigFile(c.Request.Context(), file); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "upsert_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(file))
}

// Env Variables
func (h *ConfigHandler) ListEnvVariableNames(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	opts := core.ListOptions{Offset: 0, Limit: 50}
	result, err := h.store.ListEnvVariableNames(c.Request.Context(), projectID, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(result))
}

func (h *ConfigHandler) CreateEnvVariableName(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	var req CreateConfigNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	cfgName := &models.EnvVariableName{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateEnvVariableName(c.Request.Context(), cfgName); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "create_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(cfgName))
}

func (h *ConfigHandler) DeleteEnvVariableName(c *gin.Context) {
	nameID, err := uuid.Parse(c.Param("nameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid name UUID format"}))
		return
	}

	if err := h.store.DeleteEnvVariableName(c.Request.Context(), nameID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "delete_failed", Message: err.Error()}))
		return
	}
	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}

func (h *ConfigHandler) GetEnvVariable(c *gin.Context) {
	nameID, err := uuid.Parse(c.Param("nameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid name UUID"}))
		return
	}
	envID, err := uuid.Parse(c.Param("envId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid env UUID"}))
		return
	}

	v, err := h.store.GetEnvVariable(c.Request.Context(), nameID, envID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "get_failed", Message: err.Error()}))
		return
	}
	if v == nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{Code: "not_found", Message: "Env variable not found"}))
		return
	}
	c.JSON(http.StatusOK, routes.SuccessResponse(v))
}

type UpsertEnvVariableRequest struct {
	Value string `json:"value" binding:"required"`
}

func (h *ConfigHandler) UpsertEnvVariable(c *gin.Context) {
	nameID, err := uuid.Parse(c.Param("nameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid name UUID"}))
		return
	}
	envID, err := uuid.Parse(c.Param("envId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid env UUID"}))
		return
	}

	var req UpsertEnvVariableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	now := time.Now()
	v := &models.EnvVariable{
		ID:                uuid.New(),
		EnvVariableNameID: nameID,
		EnvironmentID:     envID,
		Value:             req.Value,
		CreatedAt:         now,
		UpdatedAt:         &now,
	}

	if err := h.store.UpsertEnvVariable(c.Request.Context(), v); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "upsert_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(v))
}
