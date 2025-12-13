package middleware

import (
	"net/http"
	"strings"

	"ingredient-recognition-backend/internal/domain"
	"ingredient-recognition-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and adds user to context
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		user, err := authService.GetUserFromToken(c.Request.Context(), token)
		if err != nil {
			if err == domain.ErrUnauthorized {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			}
			c.Abort()
			return
		}

		// Store user in context
		c.Set("user", user)
		c.Set("userID", user.Id)

		c.Next()
	}
}

// GetUserFromContext extracts user from context
func GetUserFromContext(c *gin.Context) (*domain.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, domain.ErrUnauthorized
	}

	u, ok := user.(*domain.User)
	if !ok {
		return nil, domain.ErrUnauthorized
	}

	return u, nil
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", domain.ErrUnauthorized
	}

	id, ok := userID.(string)
	if !ok {
		return "", domain.ErrUnauthorized
	}

	return id, nil
}
