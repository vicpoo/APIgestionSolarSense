//api\src\sensor_thresholds\domain\threshold_repository.go

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
        thresholdGroup.PUT("/sensor/:sensor_id", handlers.SetThresholds)
    }
}