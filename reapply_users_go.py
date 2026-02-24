import re

with open("apps/api/internal/routes/v1/users.go", "r") as f:
    content = f.read()

# 1. Add PasskeyService to UserHandler
content = re.sub(
    r'(apiKeyStore  stores.UserAPIKeyStore)',
    r'\1\n\tpasskeySvc   identity.PasskeyService',
    content
)

# 2. Add APIKeyService and PasskeyService to RegisterUserRoutes
content = re.sub(
    r'(apiKeyStore stores.UserAPIKeyStore)',
    r'apiKeySvc identity.APIKeyService, passkeySvc identity.PasskeyService',
    content
)

# 3. Add to struct initialization
content = re.sub(
    r'(apiKeyStore: apiKeyStore)',
    r'apiKeySvc: apiKeySvc, passkeySvc: passkeySvc',
    content
)
content = content.replace("apiKeyStore  stores.UserAPIKeyStore", "apiKeySvc    identity.APIKeyService\n\tpasskeySvc   identity.PasskeyService")

# 4. Add route registrations
route_registration = """	// Self API Key Management
	meGroup.GET("/api-keys", h.ListMyAPIKeys)
	meGroup.POST("/api-keys", h.CreateMyAPIKey)
	meGroup.DELETE("/api-keys/:keyId", h.DeleteMyAPIKey)

	// Self Passkey Management
	meGroup.GET("/passkeys", h.ListMyPasskeys)
	meGroup.POST("/passkeys/register/start", h.StartPasskeyRegistration)
	meGroup.POST("/passkeys/register/finish", h.FinishPasskeyRegistration)
	meGroup.DELETE("/passkeys/:passkeyId", h.DeleteMyPasskey)"""

content = content.replace("""	// Self API Key Management
	meGroup.GET("/api-keys", h.ListMyAPIKeys)
	meGroup.POST("/api-keys", h.CreateMyAPIKey)
	meGroup.DELETE("/api-keys/:keyId", h.DeleteMyAPIKey)""", route_registration)

# 5. Fix ListMyAPIKeys
content = re.sub(
    r'keys, err := h.apiKeyStore.ListByUser\(c.Request.Context\(\), userID\.\(uuid\.UUID\)\)',
    r'keys, err := h.apiKeySvc.ListAPIKeys(c.Request.Context(), userID.(uuid.UUID))',
    content
)

# 6. Fix CreateMyAPIKey
old_create = """type CreateAPIKeyRequest struct {
	Name    string  `json:"name" binding:"required"`
	Comment *string `json:"comment,omitempty"`
}

func (h *UserHandler) CreateMyAPIKey(c *gin.Context) {
	_, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
			Code:    "unauthorized",
			Message: "User not authenticated",
		}))
		return
	}

	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	// This logic would ideally be in a service layer as we need to generate secure keys and hashes
	// For now, it represents the API surface boundary.
	c.JSON(http.StatusNotImplemented, routes.ErrorResponse(&routes.ApiError{
		Code: "not_implemented",
	}))
}"""

new_create = """type CreateAPIKeyRequest struct {
	Name        string  `json:"name" binding:"required"`
	Comment     *string `json:"comment,omitempty"`
	ExpiresDays *int    `json:"expiresDays,omitempty"`
}

func (h *UserHandler) CreateMyAPIKey(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
			Code:    "unauthorized",
			Message: "User not authenticated",
		}))
		return
	}

	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "validation_failed",
			Message: err.Error(),
		}))
		return
	}

	var expiresAt *time.Time
	if req.ExpiresDays != nil && *req.ExpiresDays > 0 {
		t := time.Now().AddDate(0, 0, *req.ExpiresDays)
		expiresAt = &t
	}

	apiKey, rawKey, err := h.apiKeySvc.CreateAPIKey(c.Request.Context(), userID.(uuid.UUID), req.Name, req.Comment, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "create_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusCreated, routes.SuccessResponse(map[string]any{
		"apiKey": apiKey,
		"key":    rawKey,
	}))
}"""
content = content.replace(old_create, new_create)

# 7. Fix DeleteMyAPIKey
old_delete = """func (h *UserHandler) DeleteMyAPIKey(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
			Code:    "unauthorized",
			Message: "User not authenticated",
		}))
		return
	}

	// Implementation would parse the ID and delete it ensuring the key belongs to the current user.
	c.JSON(http.StatusNotImplemented, routes.ErrorResponse(&routes.ApiError{
		Code: "not_implemented",
	}))
}"""

new_delete = """func (h *UserHandler) DeleteMyAPIKey(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
			Code:    "unauthorized",
			Message: "User not authenticated",
		}))
		return
	}

	keyIDStr := c.Param("keyId")
	keyID, err := strconv.Atoi(keyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.ErrorResponse(&routes.ApiError{
			Code:    "invalid_parameter",
			Message: "Invalid key ID",
		}))
		return
	}

	if err := h.apiKeySvc.RevokeAPIKey(c.Request.Context(), userID.(uuid.UUID), int32(keyID)); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{
			Code:    "revoke_failed",
			Message: err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse("ok"))
}"""
content = content.replace(old_delete, new_delete)

# 8. Add passkey handlers
handlers = """

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
}"""

content += handlers

# Ensure imports
if '"strconv"' not in content:
    content = content.replace('"net/http"', '"net/http"\n\t"strconv"\n\t"time"')

with open("apps/api/internal/routes/v1/users.go", "w") as f:
    f.write(content)

