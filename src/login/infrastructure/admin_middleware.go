// src/login/infrastructure/admin_middleware.go
package infrastructure

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/core"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            return
        }

        claims, err := ValidateJWTToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
            return
        }

        c.Set("jwtClaims", claims)
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
        claims, exists := c.Get("jwtClaims")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
            return
        }

        jwtClaims, ok := claims.(*core.JWTClaims)
        if !ok || !jwtClaims.IsAdmin {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
            return
        }
        c.Next()
    }
}