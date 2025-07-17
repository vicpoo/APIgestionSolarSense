// api/src/sensor_thresholds/infrastructure/threshold_routes.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/application"
)

func InitThresholdRoutes(router *gin.Engine) {
    repo := NewMySQLThresholdRepository()
    service := application.NewThresholdService(repo)
    handlers := NewThresholdHandlers(service)

    thresholdGroup := router.Group("/api/thresholds")
    {
        thresholdGroup.GET("/sensor/:sensor_id", handlers.GetThresholds)
        thresholdGroup.POST("/", handlers.CreateThreshold)
        thresholdGroup.PUT("/sensor/:sensor_id", handlers.UpdateThreshold)
        thresholdGroup.DELETE("/sensor/:sensor_id", handlers.DeleteThreshold)
    }
}