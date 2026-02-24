with open("apps/api/internal/routes/v1/users.go", "r") as f:
    content = f.read()

import re

# Replace Passkey handlers
impl = """func (h *UserHandler) ListMyPasskeys(c *gin.Context) {
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

	if err := h.passkeySvc.DeletePasskey(c.Request.Context(), passkeyID); err != nil {
		c.JSON(http.StatusInternalServerError, routes.ErrorResponse(&routes.ApiError{Code: "delete_failed", Message: err.Error()}))
		return
	}

	c.JSON(http.StatusOK, routes.SuccessResponse("ok"))
}"""

content = re.sub(
    r'func \(h \*UserHandler\) ListMyPasskeys\(c \*gin\.Context\) \{\n.*?func \(h \*UserHandler\) DeleteMyPasskey\(c \*gin\.Context\) \{\n\tc\.JSON\(http\.StatusNotImplemented, routes\.ErrorResponse\(&routes\.ApiError\{Code: "not_implemented"\}\)\)\n\}',
    impl,
    content,
    flags=re.DOTALL
)

with open("apps/api/internal/routes/v1/users.go", "w") as f:
    f.write(content)
