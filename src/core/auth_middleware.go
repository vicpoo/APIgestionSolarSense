//src/core/auth_middleware.go
package core

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

const JwtSecret = "d3c8f2b9e7a14b0f932c0d7d9a7e4f5d6c1a2e8b9f3c4a6e7b1d0f4c9a5e6b7"

type JWTClaims struct {
    UserID   int64  `json:"user_id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    AuthType string `json:"auth_type"`
    IsAdmin  bool   `json:"is_admin"`
    jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
            return
        }

        token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(JwtSecret), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        claims, ok := token.Claims.(*JWTClaims)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            return
        }

        c.Set("userID", claims.UserID)
        c.Set("userEmail", claims.Email)
        c.Set("username", claims.Username)
        c.Set("authType", claims.AuthType)
        c.Set("isAdmin", claims.IsAdmin)
        
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

