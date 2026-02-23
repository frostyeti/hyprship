package v1

import (
	"net/http"

	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/ping", Ping)
	r.GET("/healthz", Healthz)
	r.GET("/sample", Sample)
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
