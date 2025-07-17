// src/sessions/infrastructure/session_routes.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/application"
)

func InitSessionRoutes(router *gin.Engine) {
    repo := NewMySQLSessionRepository()

    // Crear casos de uso
    getUseCase := application.NewGetSessionUseCase(repo)
    postUseCase := application.NewPostSessionUseCase(repo)
    putUseCase := application.NewPutSessionUseCase(repo)
    deleteUseCase := application.NewDeleteSessionUseCase(repo)

    // Crear handlers
    getHandler := NewGetSessionHandler(getUseCase)
    postHandler := NewPostSessionHandler(postUseCase)
    putHandler := NewPutSessionHandler(putUseCase)
    deleteHandler := NewDeleteSessionHandler(deleteUseCase)

    // Crear controlador
    controller := NewSessionController(getHandler, postHandler, putHandler, deleteHandler)

    // Configurar rutas
    sessionGroup := router.Group("/api/sessions")
    {
        sessionGroup.POST("/", controller.CreateSession)
        sessionGroup.GET("/validate", controller.ValidateSession)
        sessionGroup.PUT("/refresh", controller.RefreshSession)
        sessionGroup.POST("/logout", controller.InvalidateSession)
        sessionGroup.DELETE("/:id", controller.DeleteSession)
    }
}