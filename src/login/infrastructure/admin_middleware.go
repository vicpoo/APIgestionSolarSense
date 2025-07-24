//src/login/infrastructure/admin_middleware.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}
		
		// Implementar lógica real de validación de token
		c.Set("userID", 123)
		c.Set("userEmail", "user@example.com")
		c.Set("authType", "email")
		c.Set("isAdmin", false)
		
		c.Next()
	}
}

func EmailUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authType, exists := c.Get("authType")
		if !exists || authType != "email" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "This endpoint is only for email users"})
			return
		}
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			return
		}
		c.Next()
	}
}