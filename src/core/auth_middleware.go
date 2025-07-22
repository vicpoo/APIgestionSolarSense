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

        // Para otros usuarios, establecer claims básicos
        // (En un sistema real, validarías el token JWT aquí)
        c.Set("userClaims", map[string]interface{}{
            "user_id": 2, // Esto sería extraído del token válido
            "role":    "user",
        })
        c.Next()
    }
}

func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, exists := c.Get("userClaims")
        if !exists {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authentication required"})
            return
        }
        
        claimsMap, ok := claims.(map[string]interface{})
        if !ok {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid user claims"})
            return
        }
        
        role, ok := claimsMap["role"].(string)
        if !ok || role != "admin" {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
            return
        }
        
        c.Next()
    }
}