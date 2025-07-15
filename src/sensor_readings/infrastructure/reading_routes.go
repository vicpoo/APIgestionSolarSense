//api/src/sensors/application/sensor_service.go

package infrastructure

import (
    "github.com/gin-gonic/gin"

    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/application"

)

func InitReadingRoutes(router *gin.Engine) {
    repo := NewMySQLReadingRepository()
    service := application.NewReadingService(repo)
    handlers := NewReadingHandlers(service)

    readingGroup := router.Group("/api/readings")
    {
        readingGroup.POST("/", handlers.CreateReading)
        readingGroup.GET("/sensor/:sensor_id", handlers.GetReadings)
        readingGroup.GET("/sensor/:sensor_id/latest", handlers.GetLatestReading)
    }
}