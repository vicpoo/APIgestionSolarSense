//src/alerts/infrastructure/api/alert_routes.go

package infrastructure

import (
    "github.com/gin-gonic/gin"

    "github.com/vicpoo/apigestion-solar-go/src/alerts/application"
   
)

func InitAlertRoutes(router *gin.Engine) {
    repo := NewMySQLAlertRepository()
    service := application.NewAlertService(repo)
    handlers := NewAlertHandlers(service)

    alertGroup := router.Group("/api/alerts")
    {
        alertGroup.POST("/", handlers.CreateAlert)
        alertGroup.GET("/:id", handlers.GetAlert)
        alertGroup.GET("/sensor/:sensor_id", handlers.GetSensorAlerts)
        alertGroup.GET("/unsent", handlers.GetUnsentAlerts)
        alertGroup.PUT("/:id/mark-sent", handlers.MarkAlertAsSent)
    }
}