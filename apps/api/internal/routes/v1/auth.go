package v1

import (
	"errors"
	"net/http"

	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/frostyeti/hyprship/apps/api/internal/svc/identity"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc identity.IdentityService
}

func RegisterAuthRoutes(r *gin.RouterGroup, svc identity.IdentityService) {
	h := &AuthHandler{svc: svc}
	
	r.POST("/login", h.Login)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()
	
	session, err := h.svc.LoginWithPassword(c.Request.Context(), req.Email, req.Password, &ip, &userAgent)
	if err != nil {
		code := "invalid_credentials"
		status := http.StatusUnauthorized
		if errors.Is(err, identity.ErrAccountLocked) {
			code = "account_locked"
		} else if errors.Is(err, identity.ErrMaxConcurrentSessions) {
			code = "max_sessions_reached"
			status = http.StatusForbidden
		}

		c.JSON(status, routes.ErrorResponse(&routes.ApiError{
			Code:    code,
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(session))
}
