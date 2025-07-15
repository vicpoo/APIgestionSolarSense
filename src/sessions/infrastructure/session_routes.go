//api/src/sessions/infrastructure/session_routes.go

package infrastructure

import (
    "github.com/gin-gonic/gin"

    "github.com/vicpoo/apigestion-solar-go/src/sessions/application"
  
)

func InitSessionRoutes(router *gin.Engine) {
    repo := NewMySQLSessionRepository()
    service := application.NewSessionService(repo)
    handlers := NewSessionHandlers(service)

    sessionGroup := router.Group("/api/sessions")
    {
        sessionGroup.GET("/validate", handlers.ValidateSession)
        sessionGroup.POST("/logout", handlers.Logout)
    }
}