// src/alerts/infrastructure/alert_routes.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/application"
)

func InitAlertRoutes(router *gin.Engine) {
    repo := NewMySQLAlertRepository()

    // Crear casos de uso
    postUseCase := application.NewPostAlertUseCase(repo)
    getUseCase := application.NewGetAlertUseCase(repo)
    putUseCase := application.NewPutAlertUseCase(repo)
    deleteUseCase := application.NewDeleteAlertUseCase(repo)

    // Crear handlers
    postHandler := NewPostAlertHandler(postUseCase)
    getHandler := NewGetAlertHandler(getUseCase)
    putHandler := NewPutAlertHandler(putUseCase)
    deleteHandler := NewDeleteAlertHandler(deleteUseCase)

    // Crear controlador
    controller := NewAlertController(postHandler, getHandler, putHandler, deleteHandler)

    // Configurar rutas
    alertGroup := router.Group("/api/alerts")
    {
        alertGroup.POST("/", controller.CreateAlert)
        alertGroup.GET("/:id", controller.GetAlert)
        alertGroup.GET("/sensor/:sensor_id", controller.GetSensorAlerts)
        alertGroup.GET("/unsent", controller.GetUnsentAlerts)
        alertGroup.PUT("/:id/mark-sent", controller.MarkAlertAsSent)
        alertGroup.PUT("/:id", controller.UpdateAlert)
        alertGroup.DELETE("/:id", controller.DeleteAlert)
    }
}