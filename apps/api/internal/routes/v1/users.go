package v1

import (
	"net/http"
	"strconv"

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
	apiKeySvc    identity.APIKeyService
	passkeySvc   identity.PasskeyService
	mfaSvc       identity.MfaService
}

func RegisterUserRoutes(r *gin.RouterGroup, svc identity.IdentityService, store stores.UserStore, sessionStore stores.UserSessionStore, apiKeySvc identity.APIKeyService, passkeySvc identity.PasskeyService, mfaSvc identity.MfaService) {
	h := &UserHandler{svc: svc, store: store, sessionStore: sessionStore, apiKeySvc: apiKeySvc, passkeySvc: passkeySvc, mfaSvc: mfaSvc}

	// Unauthenticated / Self Routes (Require auth, but not necessarily admin roles)
	meGroup := r.Group("/me")
	meGroup.Use(middleware.RequireAuth())
	meGroup.GET("", h.GetMe)
	meGroup.PUT("", h.UpdateMe)

	// Self Session Management
	meGroup.GET("/sessions", h.ListMySessions)
	meGroup.DELETE("/sessions/:sessionId", h.DeleteMySession)

	// MFA Management
	meGroup.POST("/mfa/setup", h.SetupMfa)
	meGroup.POST("/mfa/verify", h.VerifyAndEnableMfa)
	meGroup.DELETE("/mfa", h.DisableMfa)

	// Administrative Routes (Should require users:read or users:write)
	adminGroup := r.Group("")
	// adminGroup.Use(middleware.RequireAuth(), middleware.RequireClaim("permission", "users:read"))
	adminGroup.GET("", h.ListUsers)
	adminGroup.GET("/:id", h.GetUser)

	// adminGroupWrite.Use(middleware.RequireAuth(), middleware.RequireClaim("permission", "users:write"))
	adminGroup.POST("", h.CreateUser)
	adminGroup.PUT("/:id", h.UpdateUser)
	adminGroup.DELETE("/:id", h.DeleteUser)

	adminGroup.GET("/:id/claims", h.ListUserClaims)
	adminGroup.POST("/:id/claims", h.AddUserClaim)
	adminGroup.DELETE("/:id/claims/:claimId", h.RemoveUserClaim)
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

func (h *UserHandler) ListUserClaims(c *gin.Context) {
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

type AddUserClaimRequest struct {
	Type  string `json:"type" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func (h *UserHandler) AddUserClaim(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_id", Message: "Invalid UUID format"}))
		return
	}

	var req AddUserClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	claim := &models.UserClaim{
		UserID: id,
		Type:   req.Type,
		Value:  req.Value,
	}

	if err := h.store.AddClaim(c.Request.Context(), claim); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "add_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(claim))
}

func (h *UserHandler) RemoveUserClaim(c *gin.Context) {
	// The id parameter is the user ID, but RemoveClaim takes a claim ID
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

func (h *UserHandler) SetupMfa(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "unauthorized", Message: "User not authenticated"}))
		return
	}

	user, err := h.store.Get(c.Request.Context(), userID.(uuid.UUID))
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, routes.ErrorResponse(&routes.ApiError{Code: "not_found", Message: "User not found"}))
		return
	}

	secret, url, err := h.mfaSvc.GenerateTotpSecret(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "setup_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(map[string]string{
		"secret": secret,
		"url":    url,
	}))
}

type VerifyMfaRequest struct {
	Code string `json:"code" binding:"required"`
}

func (h *UserHandler) VerifyAndEnableMfa(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "unauthorized", Message: "User not authenticated"}))
		return
	}

	var req VerifyMfaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "validation_failed", Message: err.Error()}))
		return
	}

	if err := h.mfaSvc.VerifyAndEnableTotp(c.Request.Context(), userID.(uuid.UUID), req.Code); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "verification_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse("ok"))
}

func (h *UserHandler) DisableMfa(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "unauthorized", Message: "User not authenticated"}))
		return
	}

	if err := h.mfaSvc.DisableTotp(c.Request.Context(), userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "disable_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse("ok"))
}

func (h *UserHandler) ListMyPasskeys(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "unauthorized", Message: "User not authenticated"}))
		return
	}

	passkeys, err := h.passkeySvc.ListPasskeys(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "list_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse(passkeys))
}

func (h *UserHandler) StartPasskeyRegistration(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "unauthorized", Message: "User not authenticated"}))
		return
	}

	user, err := h.store.Get(c.Request.Context(), userID.(uuid.UUID))
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "unauthorized", Message: "User not found"}))
		return
	}

	creationData, sessionData, err := h.passkeySvc.BeginRegistration(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "registration_failed", Message: err.Error()}))
		return
	}

	sessionID := uuid.New().String()
	identity.SaveSessionData(sessionID, *sessionData)

	c.JSON(http.StatusOK, routes.SuccessResponse(map[string]any{
		"creation":  creationData,
		"sessionId": sessionID,
	}))
}

func (h *UserHandler) FinishPasskeyRegistration(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "unauthorized", Message: "User not authenticated"}))
		return
	}

	sessionID := c.Query("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "missing_session_id", Message: "Missing sessionId query parameter"}))
		return
	}

	sessionData, ok := identity.GetSessionData(sessionID)
	if !ok {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_session", Message: "Registration session expired or invalid"}))
		return
	}

	user, err := h.store.Get(c.Request.Context(), userID.(uuid.UUID))
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "unauthorized", Message: "User not found"}))
		return
	}

	credential, err := h.passkeySvc.FinishRegistration(c.Request.Context(), user, sessionData, c.Request)
	if err != nil || credential == nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "registration_failed", Message: err.Error()}))
		return
	}

	identity.DeleteSessionData(sessionID)

	c.JSON(http.StatusCreated, routes.SuccessResponse("ok"))
}

func (h *UserHandler) DeleteMyPasskey(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{Code: "unauthorized", Message: "User not authenticated"}))
		return
	}

	passkeyIDStr := c.Param("passkeyId")
	passkeyID, err := uuid.Parse(passkeyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{Code: "invalid_parameter", Message: "Invalid passkey ID"}))
		return
	}

	if err := h.passkeySvc.DeletePasskey(c.Request.Context(), userID.(uuid.UUID), passkeyID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "delete_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse("ok"))
}
