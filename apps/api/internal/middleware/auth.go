package middleware

import (
	"net/http"

	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UserIDKey    = "userID"
	SessionIDKey = "sessionID"
)

// RequireAuth is a placeholder middleware that ensures a user is authenticated.
// In a real implementation, it would parse the JWT or session token from the Authorization header or cookie,
// validate it, and set the user ID in the context.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example: Extract token, validate, get userID
		// For now, if there is a mocked header, we use it for testing, else fail
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
				Code:    "unauthorized",
				Message: "Missing authentication token",
			}))
			c.Abort()
			return
		}

		// Temporary mock logic for testing:
		// If the header is "Bearer mock-user-id", we parse it.
		// In production, we'd verify the JWT and extract claims.
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token := authHeader[7:]
			userID, err := uuid.Parse(token)
			if err == nil {
				c.Set(UserIDKey, userID)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, routes.ErrorResponse(&routes.ApiError{
			Code:    "unauthorized",
			Message: "Invalid authentication token",
		}))
		c.Abort()
	}
}

// RequireClaim checks if the authenticated user has a specific claim/permission.
func RequireClaim(claimType, claimValue string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Validate user claims from JWT or database
		c.Next()
	}
}
