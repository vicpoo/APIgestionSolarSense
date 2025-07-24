package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el email del usuario autenticado (ajusta según tu sistema de autenticación)
		userEmail, exists := c.Get("userEmail")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		// Verificar si es el admin
		if userEmail != "admin@integrador.com" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin privileges required"})
			return
		}

		c.Next()
	}
}