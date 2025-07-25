// src/login/infrastructure/admin_middleware.go
package infrastructure

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/core"
	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := domain.ValidateJWTToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        // Verificar que el usuario aún existe
        if _, err := c.MustGet("authRepository").(domain.AuthRepository).GetUserByID(c.Request.Context(), claims.UserID); err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User no longer exists"})
            return
        }

        c.Set("userClaims", claims)
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