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

type SecretHandler struct {
	store stores.SecretStore
}

func RegisterSecretRoutes(projectGroup *gin.RouterGroup, envGroup *gin.RouterGroup, store stores.SecretStore) {
	h := &SecretHandler{store: store}

	// Project level names
	projectGroup.GET("/secrets", h.ListSecretNames)
	projectGroup.POST("/secrets", h.CreateSecretName)
	projectGroup.DELETE("/secrets/:nameId", h.DeleteSecretName)

	projectGroup.GET("/certificates", h.ListCertificateNames)
	projectGroup.POST("/certificates", h.CreateCertificateName)
	projectGroup.DELETE("/certificates/:nameId", h.DeleteCertificateName)

	projectGroup.GET("/ssh-keys", h.ListSSHKeyNames)
	projectGroup.POST("/ssh-keys", h.CreateSSHKeyName)
	projectGroup.DELETE("/ssh-keys/:nameId", h.DeleteSSHKeyName)

	// Environment level values
	envGroup.GET("/secrets/:nameId", h.GetSecret)
	envGroup.PUT("/secrets/:nameId", h.UpsertSecret)

	envGroup.GET("/certificates/:nameId", h.GetCertificate)
	envGroup.PUT("/certificates/:nameId", h.UpsertCertificate)

	envGroup.GET("/ssh-keys/:nameId", h.GetSSHKey)
	envGroup.PUT("/ssh-keys/:nameId", h.UpsertSSHKey)
}

// Request models
type CreateNameRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpsertSecretRequest struct {
	Value string `json:"value" binding:"required"`
}

type UpsertCertificateRequest struct {
	CertificateData  string  `json:"certificateData" binding:"required"`
	PrivateKeySecret *string `json:"privateKeySecret,omitempty"`
}

type UpsertSSHKeyRequest struct {
	PublicKey        string `json:"publicKey" binding:"required"`
	PrivateKeySecret string `json:"privateKeySecret" binding:"required"`
}

// Secrets
func (h *SecretHandler) ListSecretNames(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	opts := core.ListOptions{Offset: 0, Limit: 50}
	result, err := h.store.ListSecretNames(c.Request.Context(), projectID, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(result))
}

func (h *SecretHandler) CreateSecretName(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	var req CreateNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	name := &models.SecretName{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateSecretName(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "create_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(name))
}

func (h *SecretHandler) DeleteSecretName(c *gin.Context) {
	nameID, err := uuid.Parse(c.Param("nameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid name UUID format"}))
		return
	}

	if err := h.store.DeleteSecretName(c.Request.Context(), nameID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "delete_failed", Message: err.Error()}))
		return
	}
	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}

func (h *SecretHandler) GetSecret(c *gin.Context) {
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

	v, err := h.store.GetSecret(c.Request.Context(), nameID, envID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "get_failed", Message: err.Error()}))
		return
	}
	if v == nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{Code: "not_found", Message: "Secret not found"}))
		return
	}
	c.JSON(http.StatusOK, routes.SuccessResponse(v))
}

func (h *SecretHandler) UpsertSecret(c *gin.Context) {
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

	var req UpsertSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	now := time.Now()
	v := &models.Secret{
		ID:            uuid.New(),
		SecretNameID:  nameID,
		EnvironmentID: envID,
		Value:         req.Value,
		CreatedAt:     now,
		UpdatedAt:     &now,
	}

	if err := h.store.UpsertSecret(c.Request.Context(), v); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "upsert_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(v))
}

// Certificates
func (h *SecretHandler) ListCertificateNames(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	opts := core.ListOptions{Offset: 0, Limit: 50}
	result, err := h.store.ListCertificateNames(c.Request.Context(), projectID, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(result))
}

func (h *SecretHandler) CreateCertificateName(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	var req CreateNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	name := &models.CertificateName{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateCertificateName(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "create_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(name))
}

func (h *SecretHandler) DeleteCertificateName(c *gin.Context) {
	nameID, err := uuid.Parse(c.Param("nameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid name UUID format"}))
		return
	}

	if err := h.store.DeleteCertificateName(c.Request.Context(), nameID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "delete_failed", Message: err.Error()}))
		return
	}
	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}

func (h *SecretHandler) GetCertificate(c *gin.Context) {
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

	v, err := h.store.GetCertificate(c.Request.Context(), nameID, envID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "get_failed", Message: err.Error()}))
		return
	}
	if v == nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{Code: "not_found", Message: "Certificate not found"}))
		return
	}
	c.JSON(http.StatusOK, routes.SuccessResponse(v))
}

func (h *SecretHandler) UpsertCertificate(c *gin.Context) {
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

	var req UpsertCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	now := time.Now()
	v := &models.Certificate{
		ID:                uuid.New(),
		CertificateNameID: nameID,
		EnvironmentID:     envID,
		CertificateData:   req.CertificateData,
		PrivateKeySecret:  req.PrivateKeySecret,
		CreatedAt:         now,
		UpdatedAt:         &now,
	}

	if err := h.store.UpsertCertificate(c.Request.Context(), v); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "upsert_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(v))
}

// SSH Keys
func (h *SecretHandler) ListSSHKeyNames(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	opts := core.ListOptions{Offset: 0, Limit: 50}
	result, err := h.store.ListSSHKeyNames(c.Request.Context(), projectID, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(result))
}

func (h *SecretHandler) CreateSSHKeyName(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid project UUID"}))
		return
	}

	var req CreateNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	name := &models.SSHKeyName{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateSSHKeyName(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "create_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(name))
}

func (h *SecretHandler) DeleteSSHKeyName(c *gin.Context) {
	nameID, err := uuid.Parse(c.Param("nameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid name UUID format"}))
		return
	}

	if err := h.store.DeleteSSHKeyName(c.Request.Context(), nameID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "delete_failed", Message: err.Error()}))
		return
	}
	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}

func (h *SecretHandler) GetSSHKey(c *gin.Context) {
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

	v, err := h.store.GetSSHKey(c.Request.Context(), nameID, envID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "get_failed", Message: err.Error()}))
		return
	}
	if v == nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{Code: "not_found", Message: "SSH key not found"}))
		return
	}
	c.JSON(http.StatusOK, routes.SuccessResponse(v))
}

func (h *SecretHandler) UpsertSSHKey(c *gin.Context) {
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

	var req UpsertSSHKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	now := time.Now()
	v := &models.SSHKey{
		ID:               uuid.New(),
		SSHKeyNameID:     nameID,
		EnvironmentID:    envID,
		PublicKey:        req.PublicKey,
		PrivateKeySecret: req.PrivateKeySecret,
		CreatedAt:        now,
		UpdatedAt:        &now,
	}

	if err := h.store.UpsertSSHKey(c.Request.Context(), v); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "upsert_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(v))
}
