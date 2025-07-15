//api/src/sensors/application/sensor_service.go

package infrastructure

import (
    "github.com/gin-gonic/gin"
   
    "github.com/vicpoo/apigestion-solar-go/src/sensors/application"
  
)

func InitSensorRoutes(router *gin.Engine) {
    repo := NewMySQLSensorRepository()
    service := application.NewSensorService(repo)
    handlers := NewSensorHandlers(service)

    sensorGroup := router.Group("/api/sensors")
    {
        sensorGroup.POST("/", handlers.CreateSensor)
        sensorGroup.GET("/:id", handlers.GetSensor)
        sensorGroup.GET("/user/:user_id", handlers.GetUserSensors)
        sensorGroup.PUT("/:id", handlers.UpdateSensor)
        sensorGroup.DELETE("/:id", handlers.DeleteSensor)
    }
}