//go:build e2e
// +build e2e

package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	v1 "github.com/frostyeti/hyprship/apps/api/internal/routes/v1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	v1.RegisterRoutes(apiV1)
	return r
}

func TestE2E_Ping(t *testing.T) {
	r := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp routes.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
	assert.Equal(t, "pong", resp.Value)
}

func TestE2E_Healthz(t *testing.T) {
	r := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/healthz", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestE2E_Sample(t *testing.T) {
	r := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/sample", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp routes.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)

	val, ok := resp.Value.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "sample", val["name"])
}
