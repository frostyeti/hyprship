package v1

import (
	"net/http"

	"github.com/frostyeti/hyprship/apps/api/internal/middleware"
	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/frostyeti/hyprship/apps/api/internal/svc/identity"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, identitySvc identity.IdentityService, userStore stores.UserStore, roleStore stores.RoleStore, groupStore stores.GroupStore, projectStore stores.ProjectStore, sessionStore stores.UserSessionStore, apiKeySvc identity.APIKeyService, passkeySvc identity.PasskeyService, mfaSvc identity.MfaService, verificationSvc identity.VerificationService, importExportSvc identity.ImportExportService) {
	r.GET("/ping", Ping)
	r.GET("/healthz", Healthz)
	r.GET("/sample", Sample)

	authGroup := r.Group("/auth")
	authGroup.Use(middleware.RateLimit(5, 10))

	RegisterAuthRoutes(authGroup, identitySvc, userStore, passkeySvc, mfaSvc)

	usersGroup := r.Group("/users")
	RegisterUserRoutes(usersGroup, identitySvc, userStore, sessionStore, apiKeySvc, passkeySvc, mfaSvc, verificationSvc, importExportSvc)

	rolesGroup := r.Group("/roles")
	RegisterRoleRoutes(rolesGroup, identitySvc, roleStore, importExportSvc)

	groupsGroup := r.Group("/groups")
	RegisterGroupRoutes(groupsGroup, groupStore)

	projectsGroup := r.Group("/projects")
	RegisterProjectRoutes(projectsGroup, projectStore)
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, routes.SuccessResponse("pong"))
}

func Healthz(c *gin.Context) {
	// standard health check without envelope
	c.String(http.StatusOK, "OK")
}

func Sample(c *gin.Context) {
	c.JSON(http.StatusOK, routes.SuccessResponse(map[string]any{
		"name": "sample",
	}))
}
