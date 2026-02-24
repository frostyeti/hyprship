package v1

import (
	"errors"
	"net/http"

	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/frostyeti/hyprship/apps/api/internal/svc/identity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	svc        identity.IdentityService
	userStore  stores.UserStore
	passkeySvc identity.PasskeyService
	mfaSvc     identity.MfaService
}

func RegisterAuthRoutes(r *gin.RouterGroup, svc identity.IdentityService, userStore stores.UserStore, passkeySvc identity.PasskeyService, mfaSvc identity.MfaService) {
	h := &AuthHandler{svc: svc, userStore: userStore, passkeySvc: passkeySvc, mfaSvc: mfaSvc}

	r.POST("/login", h.Login)
	r.POST("/login/mfa", h.LoginMfa)

	// Passkey login
	r.POST("/login/passkey/start", h.StartPasskeyLogin)
	r.POST("/login/passkey/finish", h.FinishPasskeyLogin)
	r.POST("/logout", h.Logout)
	r.POST("/forgot-password", h.ForgotPassword)
	r.POST("/reset-password", h.ResetPassword)
	r.POST("/forgot-email", h.ForgotEmail)
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

	session, mfaToken, err := h.svc.LoginWithPassword(c.Request.Context(), req.Email, req.Password, &ip, &userAgent)
	if err != nil {
		if errors.Is(err, identity.ErrMfaRequired) {
			c.JSON(http.StatusOK, routes.SuccessResponse(map[string]any{
				"mfaRequired": true,
				"mfaToken":    mfaToken,
			}))
			return
		}

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

type LoginMfaRequest struct {
	MfaToken string `json:"mfaToken" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

func (h *AuthHandler) LoginMfa(c *gin.Context) {
	var req LoginMfaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	userID, ok := identity.GetMfaSessionData(req.MfaToken)
	if !ok {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "invalid_mfa_token", Message: "MFA session expired or invalid"}))
		return
	}

	// Verify the MFA Code
	valid, err := h.mfaSvc.VerifyTotp(c.Request.Context(), userID, req.Code)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "invalid_mfa_code", Message: "Invalid MFA Code"}))
		return
	}

	identity.DeleteMfaSessionData(req.MfaToken)

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	session, err := h.svc.CreateSession(c.Request.Context(), userID, &ip, &userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "session_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(session))
}

type LogoutRequest struct {
	SessionID string `json:"sessionId" binding:"required,uuid"`
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	sessionUUID, _ := uuid.Parse(req.SessionID)

	if err := h.svc.Logout(c.Request.Context(), sessionUUID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "logout_failed",
			Message: "Failed to revoke session",
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse("ok"))
}

type ForgotEmailRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (h *AuthHandler) ForgotEmail(c *gin.Context) {
	var req ForgotEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	// TODO: implement service logic to send SMS with email
	// h.svc.ForgotEmail(c.Request.Context(), req.Phone)

	c.JSON(http.StatusOK, routes.SuccessResponse("ok"))
}

type StartPasskeyLoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) StartPasskeyLogin(c *gin.Context) {
	var req StartPasskeyLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	user, err := h.userStore.GetByEmail(c.Request.Context(), req.Email)
	if err != nil || user == nil {
		// Do not reveal user existence
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_credentials",
			Message: "Invalid credentials",
		}))
		return
	}

	assertion, sessionData, err := h.passkeySvc.BeginLogin(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "login_failed",
			Message: err.Error(),
		}))
		return
	}

	sessionID := uuid.New().String()
	identity.SaveSessionData(sessionID, *sessionData)

	c.JSON(http.StatusOK, routes.SuccessResponse(map[string]any{
		"assertion": assertion,
		"sessionId": sessionID,
	}))
}

func (h *AuthHandler) FinishPasskeyLogin(c *gin.Context) {
	sessionID := c.Query("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "missing_session_id",
			Message: "Missing sessionId query parameter",
		}))
		return
	}

	sessionData, ok := identity.GetSessionData(sessionID)
	if !ok {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_session",
			Message: "Login session expired or invalid",
		}))
		return
	}

	userID, err := uuid.Parse(string(sessionData.UserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "internal_error", Message: err.Error()}))
		return
	}

	user, err := h.userStore.Get(c.Request.Context(), userID)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "invalid_credentials", Message: "Invalid credentials"}))
		return
	}

	credential, err := h.passkeySvc.FinishLogin(c.Request.Context(), user, sessionData, c.Request)
	if err != nil || credential == nil {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
			Code:    "login_failed",
			Message: err.Error(),
		}))
		return
	}

	identity.DeleteSessionData(sessionID)

	var ipAddress *string
	if ip := c.ClientIP(); ip != "" {
		ipAddress = &ip
	}
	var userAgent *string
	if ua := c.Request.UserAgent(); ua != "" {
		userAgent = &ua
	}

	userSession, err := h.svc.CreateSession(c.Request.Context(), user.ID, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "session_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(userSession))
}
