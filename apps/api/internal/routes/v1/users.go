package v1

import (
	"net/http"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/middleware"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/frostyeti/hyprship/apps/api/internal/svc/identity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	svc          identity.IdentityService
	store        stores.UserStore
	sessionStore stores.UserSessionStore
}

func RegisterUserRoutes(r *gin.RouterGroup, svc identity.IdentityService, store stores.UserStore, sessionStore stores.UserSessionStore) {
	h := &UserHandler{svc: svc, store: store, sessionStore: sessionStore}

	// Unauthenticated / Self Routes (Require auth, but not necessarily admin roles)
	meGroup := r.Group("/me")
	meGroup.Use(middleware.RequireAuth())
	meGroup.GET("", h.GetMe)
	meGroup.PUT("", h.UpdateMe)

	// Self Session Management
	meGroup.GET("/sessions", h.ListMySessions)
	meGroup.DELETE("/sessions/:sessionId", h.DeleteMySession)

	// Administrative Routes (Should require users:read or users:write)
	adminGroup := r.Group("")
	// adminGroup.Use(middleware.RequireAuth(), middleware.RequireClaim("permission", "users:read"))
	adminGroup.GET("", h.ListUsers)
	adminGroup.GET("/:id", h.GetUser)

	// adminGroupWrite.Use(middleware.RequireAuth(), middleware.RequireClaim("permission", "users:write"))
	adminGroup.POST("", h.CreateUser)
	adminGroup.PUT("/:id", h.UpdateUser)
	adminGroup.DELETE("/:id", h.DeleteUser)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
			Code:    "unauthorized",
			Message: "User not authenticated",
		}))
		return
	}

	user, err := h.store.Get(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "User not found",
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(user))
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
			Code:    "unauthorized",
			Message: "User not authenticated",
		}))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	user, err := h.store.Get(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "User not found",
		}))
		return
	}

	if req.Name != nil {
		user.Name = req.Name
	}
	if req.Phone != nil {
		user.PrimaryPhone = req.Phone
	}

	if err := h.store.Update(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "update_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(user))
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	// Simple list options from query params
	// TODO: fully parse options
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

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_id",
			Message: "Invalid UUID format",
		}))
		return
	}

	user, err := h.store.Get(c.Request.Context(), id)
	if err != nil {
		// TODO: distinguish between not found and server error
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "User not found",
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(user))
}

type CreateUserRequest struct {
	Email string  `json:"email" binding:"required,email"`
	Name  *string `json:"name,omitempty"`
	Phone *string `json:"phone,omitempty"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	user := &models.User{
		ID:           uuid.New(),
		PrimaryEmail: &req.Email,
		Name:         req.Name,
		PrimaryPhone: req.Phone,
	}

	if err := h.store.Create(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "create_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(user))
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	IsBanned *bool   `json:"isBanned,omitempty"`
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_id",
			Message: "Invalid UUID format",
		}))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	user, err := h.store.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "User not found",
		}))
		return
	}

	if req.Name != nil {
		user.Name = req.Name
	}
	if req.Phone != nil {
		user.PrimaryPhone = req.Phone
	}
	if req.IsBanned != nil {
		user.IsBanned = *req.IsBanned
	}

	if err := h.store.Update(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "update_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(user))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
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

func (h *UserHandler) ListMySessions(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
			Code:    "unauthorized",
			Message: "User not authenticated",
		}))
		return
	}

	sessions, err := h.sessionStore.ListByUser(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "list_failed",
			Message: "Failed to list sessions",
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(sessions))
}

func (h *UserHandler) DeleteMySession(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
			Code:    "unauthorized",
			Message: "User not authenticated",
		}))
		return
	}

	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_id",
			Message: "Invalid UUID format",
		}))
		return
	}

	// Verify the session actually belongs to this user before deleting
	session, err := h.sessionStore.Get(c.Request.Context(), sessionID)
	if err != nil || session.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{
			Code:    "not_found",
			Message: "Session not found",
		}))
		return
	}

	if err := h.svc.Logout(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "delete_failed",
			Message: "Failed to revoke session",
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse("ok"))
}
