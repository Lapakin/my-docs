package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lapotkin/file-storage/internal/adapter/http/rest"
	"github.com/lapotkin/file-storage/internal/auth/jwt"
	"github.com/lapotkin/file-storage/pkg/models"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

const (
	UserIDKey       = "user_id"
	UserEmailKey    = "user_email"
	UserUsernameKey = "user_username"
	UserRoleKey     = "user_role"
	ClaimsKey       = "claims"
)

// Auth validates JWT token from Authorization header
func Auth(tokenManager *jwt.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			rest.RespondError(c, http.StatusUnauthorized, rest.ErrUnauthorized)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			rest.RespondError(c, http.StatusUnauthorized, rest.ErrInvalidAuthHeader)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, BearerPrefix)
		claims, err := tokenManager.ValidateToken(token)
		if err != nil {
			rest.RespondError(c, http.StatusUnauthorized, rest.ErrInvalidToken)
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserUsernameKey, claims.Username)
		c.Set(UserRoleKey, claims.Role)
		c.Set(ClaimsKey, claims)

		c.Next()
	}
}

// Admin ensures the user has admin role
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(UserRoleKey)
		if !exists {
			rest.RespondError(c, http.StatusUnauthorized, rest.ErrUnauthorized)
			c.Abort()
			return
		}

		if role != models.RoleAdmin {
			rest.RespondError(c, http.StatusForbidden, rest.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

// CookieAuthMiddleware checks for authentication token in cookies
func CookieAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
