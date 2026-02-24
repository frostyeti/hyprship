with open("apps/api/internal/routes/v1/auth.go", "r") as f:
    content = f.read()

import re

# Inject UserStore
content = re.sub(
    r'type AuthHandler struct \{\n\tsvc        identity\.IdentityService\n\tpasskeySvc identity\.PasskeyService\n\}',
    r'type AuthHandler struct {\n\tsvc        identity.IdentityService\n\tuserStore  stores.UserStore\n\tpasskeySvc identity.PasskeyService\n}',
    content
)

content = re.sub(
    r'func RegisterAuthRoutes\(r \*gin\.RouterGroup, svc identity\.IdentityService, passkeySvc identity\.PasskeyService\) \{',
    r'func RegisterAuthRoutes(r *gin.RouterGroup, svc identity.IdentityService, userStore stores.UserStore, passkeySvc identity.PasskeyService) {',
    content
)

content = re.sub(
    r'h := &AuthHandler\{svc: svc, passkeySvc: passkeySvc\}',
    r'h := &AuthHandler{svc: svc, userStore: userStore, passkeySvc: passkeySvc}',
    content
)

# Replace Start and Finish Login
impl = """type StartPasskeyLoginRequest struct {
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
}"""

content = re.sub(
    r'func \(h \*AuthHandler\) StartPasskeyLogin\(c \*gin\.Context\) \{\n\tc\.JSON\(http\.StatusNotImplemented.*\}\n\nfunc \(h \*AuthHandler\) FinishPasskeyLogin\(c \*gin\.Context\) \{\n\tc\.JSON\(http\.StatusNotImplemented.*\}',
    impl,
    content,
    flags=re.DOTALL
)

with open("apps/api/internal/routes/v1/auth.go", "w") as f:
    f.write(content)
