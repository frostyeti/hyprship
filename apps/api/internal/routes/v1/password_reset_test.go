package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestForgotPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockSvc := &MockIdentityService{
		ForgotPasswordFunc: func(ctx context.Context, email string) error {
			assert.Equal(t, "user@example.com", email)
			return nil
		},
	}

	RegisterAuthRoutes(r.Group("/api/v1/auth"), mockSvc, nil, nil, nil)

	body, _ := json.Marshal(ForgotPasswordRequest{
		Email: "user@example.com",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/forgot-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response routes.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)
	assert.Equal(t, "ok", response.Value)
}

func TestResetPassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockSvc := &MockIdentityService{
		ResetPasswordFunc: func(ctx context.Context, email, token, newPassword string) error {
			assert.Equal(t, "user@example.com", email)
			assert.Equal(t, "securetoken123", token)
			assert.Equal(t, "NewPassword123!", newPassword)
			return nil
		},
	}

	RegisterAuthRoutes(r.Group("/api/v1/auth"), mockSvc, nil, nil, nil)

	body, _ := json.Marshal(ResetPasswordRequest{
		Email:       "user@example.com",
		Token:       "securetoken123",
		NewPassword: "NewPassword123!",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/reset-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response routes.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Ok)
	assert.Equal(t, "ok", response.Value)
}

func TestResetPassword_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockSvc := &MockIdentityService{
		ResetPasswordFunc: func(ctx context.Context, email, token, newPassword string) error {
			return errors.New("invalid token")
		},
	}

	RegisterAuthRoutes(r.Group("/api/v1/auth"), mockSvc, nil, nil, nil)

	body, _ := json.Marshal(ResetPasswordRequest{
		Email:       "user@example.com",
		Token:       "invalidtoken",
		NewPassword: "NewPassword123!",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/reset-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response routes.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Ok)
}
