package middleware

import (
	"net/http"
	"strings"

	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UserIDKey    = "userID"
	SessionIDKey = "sessionID"
	RolesKey     = "roles"
	GroupsKey    = "groups"
	ClaimsKey    = "claims"
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

// ClaimsMiddleware injects the authenticated user's groups, roles, and claims into the context.
func ClaimsMiddleware(userStore stores.UserStore, groupStore stores.GroupStore, roleStore stores.RoleStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get(UserIDKey)
		if !exists {
			c.Next()
			return
		}
		userID := userIDVal.(uuid.UUID)

		var roleIDs []string
		var groupIDs []string
		var claims []string

		ctx := c.Request.Context()

		// Fetch Groups for user
		if groupStore != nil {
			groups, err := groupStore.ListGroupsByUserID(ctx, userID)
			if err == nil {
				for _, g := range groups {
					groupIDs = append(groupIDs, g.GroupID.String())
					// Fetch roles for this group
					groupRoles, err := groupStore.ListRoles(ctx, g.GroupID)
					if err == nil {
						for _, gr := range groupRoles {
							roleIDs = append(roleIDs, gr.RoleID.String())
						}
					}
				}
			}
		}

		// Fetch User Claims
		if userStore != nil {
			userClaims, err := userStore.ListClaims(ctx, userID)
			if err == nil {
				for _, uc := range userClaims {
					claims = append(claims, uc.Type+":"+uc.Value)
				}
			}
		}

		// Fetch Role Claims for each collected role
		if roleStore != nil {
			for _, rIDStr := range roleIDs {
				rID, err := uuid.Parse(rIDStr)
				if err == nil {
					roleClaims, err := roleStore.ListClaims(ctx, rID)
					if err == nil {
						for _, rc := range roleClaims {
							claims = append(claims, rc.Type+":"+rc.Value)
						}
					}
				}
			}
		}

		c.Set(GroupsKey, groupIDs)
		c.Set(RolesKey, roleIDs)
		c.Set(ClaimsKey, claims)

		c.Next()
	}
}

// RequireClaim checks if the authenticated user has a specific claim/permission.
func RequireClaim(claimType, claimValue string) gin.HandlerFunc {
	expectedClaim := claimType + ":" + claimValue
	return func(c *gin.Context) {
		claimsVal, exists := c.Get(ClaimsKey)
		if !exists {
			c.JSON(http.StatusForbidden, routes.ErrorResponse(&routes.ApiError{
				Code:    "forbidden",
				Message: "Missing required permissions",
			}))
			c.Abort()
			return
		}

		claims, ok := claimsVal.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, routes.ErrorResponse(&routes.ApiError{
				Code:    "forbidden",
				Message: "Missing required permissions",
			}))
			c.Abort()
			return
		}

		hasClaim := false
		for _, claim := range claims {
			if strings.EqualFold(claim, expectedClaim) {
				hasClaim = true
				break
			}
		}

		if !hasClaim {
			c.JSON(http.StatusForbidden, routes.ErrorResponse(&routes.ApiError{
				Code:    "forbidden",
				Message: "Missing required permissions",
			}))
			c.Abort()
			return
		}

		c.Next()
	}
}
