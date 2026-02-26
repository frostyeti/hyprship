package v1

import (
	"net/http"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupHandler struct {
	store stores.GroupStore
}

func RegisterGroupRoutes(r *gin.RouterGroup, store stores.GroupStore) {
	h := &GroupHandler{store: store}

	r.GET("", h.ListGroups)
	r.GET("/:id", h.GetGroup)
	r.POST("", h.CreateGroup)
	r.PUT("/:id", h.UpdateGroup)
	r.DELETE("/:id", h.DeleteGroup)

	r.GET("/:id/users", h.ListGroupUsers)
	r.POST("/:id/users/:userId", h.AddGroupUser)
	r.DELETE("/:id/users/:userId", h.RemoveGroupUser)

	r.GET("/:id/admins", h.ListGroupAdmins)
	r.POST("/:id/admins/:userId", h.AddGroupAdmin)
	r.DELETE("/:id/admins/:userId", h.RemoveGroupAdmin)

	r.GET("/:id/roles", h.ListGroupRoles)
	r.POST("/:id/roles/:roleId", h.AddGroupRole)
	r.DELETE("/:id/roles/:roleId", h.RemoveGroupRole)
}

func (h *GroupHandler) ListGroups(c *gin.Context) {
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

func (h *GroupHandler) GetGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_id",
			Message: "Invalid UUID format",
		}))
		return
	}

	group, err := h.store.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "Group not found",
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(group))
}

type CreateGroupRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	ImageURI    *string `json:"imageUri,omitempty"`
	Email       *string `json:"email,omitempty"`
	IsActive    bool    `json:"isActive"`
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	group := &models.Group{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		ImageURI:    req.ImageURI,
		Email:       req.Email,
		IsActive:    req.IsActive,
	}

	if err := h.store.Create(c.Request.Context(), group); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "create_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(group))
}

type UpdateGroupRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	ImageURI    *string `json:"imageUri,omitempty"`
	Email       *string `json:"email,omitempty"`
	IsActive    bool    `json:"isActive"`
}

func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_id",
			Message: "Invalid UUID format",
		}))
		return
	}

	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	group, err := h.store.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "Group not found",
		}))
		return
	}
	if group == nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "Group not found",
		}))
		return
	}

	group.Name = req.Name
	group.Description = req.Description
	group.ImageURI = req.ImageURI
	group.Email = req.Email
	group.IsActive = req.IsActive

	if err := h.store.Update(c.Request.Context(), group); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "update_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(group))
}

func (h *GroupHandler) DeleteGroup(c *gin.Context) {
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

func (h *GroupHandler) ListGroupUsers(c *gin.Context) {
	idStr := c.Param("id")
	groupID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid UUID format"}))
		return
	}

	users, err := h.store.ListUsers(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(users))
}

func (h *GroupHandler) AddGroupUser(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid group UUID format"}))
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid user UUID format"}))
		return
	}

	if err := h.store.AddUser(c.Request.Context(), groupID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "add_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse[any](nil))
}

func (h *GroupHandler) RemoveGroupUser(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid group UUID format"}))
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid user UUID format"}))
		return
	}

	if err := h.store.RemoveUser(c.Request.Context(), groupID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "remove_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}

func (h *GroupHandler) ListGroupAdmins(c *gin.Context) {
	idStr := c.Param("id")
	groupID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid UUID format"}))
		return
	}

	admins, err := h.store.ListAdmins(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(admins))
}

func (h *GroupHandler) AddGroupAdmin(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid group UUID format"}))
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid user UUID format"}))
		return
	}

	if err := h.store.AddAdmin(c.Request.Context(), groupID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "add_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse[any](nil))
}

func (h *GroupHandler) RemoveGroupAdmin(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid group UUID format"}))
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid user UUID format"}))
		return
	}

	if err := h.store.RemoveAdmin(c.Request.Context(), groupID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "remove_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}

func (h *GroupHandler) ListGroupRoles(c *gin.Context) {
	idStr := c.Param("id")
	groupID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid UUID format"}))
		return
	}

	roles, err := h.store.ListRoles(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(roles))
}

func (h *GroupHandler) AddGroupRole(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid group UUID format"}))
		return
	}

	roleID, err := uuid.Parse(c.Param("roleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid role UUID format"}))
		return
	}

	if err := h.store.AddRole(c.Request.Context(), groupID, roleID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "add_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse[any](nil))
}

func (h *GroupHandler) RemoveGroupRole(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid group UUID format"}))
		return
	}

	roleID, err := uuid.Parse(c.Param("roleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid role UUID format"}))
		return
	}

	if err := h.store.RemoveRole(c.Request.Context(), groupID, roleID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "remove_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse[any](nil))
}
