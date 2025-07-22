// src/core/auth_middleware.go
package core

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
            return
        }

        // Verificar si es el token del admin
        if token == "$2a$10$Wgs8JL0bkgcZ2MZ2uNy.MY10bcBr6Dw6X7n55f2Q/QQj2Ws0" {
            c.Set("userClaims", map[string]interface{}{
                "user_id": 1,
                "role":    "admin",
            })
            c.Next()
            return
        }

        // Para otros usuarios, implementar lógica de validación real
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
    }
}