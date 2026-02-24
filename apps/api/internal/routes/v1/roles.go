package v1

import (
	"net/http"
	"strconv"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/frostyeti/hyprship/apps/api/internal/svc/identity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleHandler struct {
	svc   identity.IdentityService
	store stores.RoleStore
}

func RegisterRoleRoutes(r *gin.RouterGroup, svc identity.IdentityService, store stores.RoleStore) {
	h := &RoleHandler{svc: svc, store: store}

	r.GET("", h.ListRoles)
	r.GET("/:id", h.GetRole)
	r.POST("", h.CreateRole)
	r.PUT("/:id", h.UpdateRole)
	r.DELETE("/:id", h.DeleteRole)

	r.GET("/:id/claims", h.ListRoleClaims)
	r.POST("/:id/claims", h.AddRoleClaim)
	r.DELETE("/:id/claims/:claimId", h.RemoveRoleClaim)
}

func (h *RoleHandler) ListRoles(c *gin.Context) {
	opts := core.ListOptions{
		Page:     1,
		PageSize: 50,
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

func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_id",
			Message: "Invalid UUID format",
		}))
		return
	}

	role, err := h.store.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "Role not found",
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(role))
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	role := &models.Role{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.store.Create(c.Request.Context(), role); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "create_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(role))
}

type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_id",
			Message: "Invalid UUID format",
		}))
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	role, err := h.store.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "Role not found",
		}))
		return
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := h.store.Update(c.Request.Context(), role); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "update_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(role))
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
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

	c.JSON(http.StatusOK, routes.SuccessResponse("ok"))
}

func (h *RoleHandler) ListRoleClaims(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid UUID format"}))
		return
	}

	claims, err := h.store.ListClaims(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(claims))
}

type AddRoleClaimRequest struct {
	Type  string `json:"type" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func (h *RoleHandler) AddRoleClaim(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid UUID format"}))
		return
	}

	var req AddRoleClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	claim := &models.RoleClaim{
		RoleID: id,
		Type:   req.Type,
		Value:  req.Value,
	}

	if err := h.store.AddClaim(c.Request.Context(), claim); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "add_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(claim))
}

func (h *RoleHandler) RemoveRoleClaim(c *gin.Context) {
	claimIdStr := c.Param("claimId")
	claimID, err := strconv.ParseInt(claimIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_parameter", Message: "Invalid claim ID"}))
		return
	}

	if err := h.store.RemoveClaim(c.Request.Context(), int32(claimID)); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "remove_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse("ok"))
}
