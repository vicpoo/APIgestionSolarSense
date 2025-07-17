// src/sensors/infrastructure/sensor_routes.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/application"
)

func InitSensorRoutes(router *gin.Engine) {
    repo := NewMySQLSensorRepository()

    // Crear casos de uso
    getUseCase := application.NewGetSensorUseCase(repo)
    postUseCase := application.NewPostSensorUseCase(repo)
    putUseCase := application.NewPutSensorUseCase(repo)
    deleteUseCase := application.NewDeleteSensorUseCase(repo)

    // Crear handlers
    getHandler := NewGetSensorHandler(getUseCase)
    postHandler := NewPostSensorHandler(postUseCase)
    putHandler := NewPutSensorHandler(putUseCase)
    deleteHandler := NewDeleteSensorHandler(deleteUseCase)

    // Crear controlador
    controller := NewSensorController(getHandler, postHandler, putHandler, deleteHandler)

    // Configurar rutas
    sensorGroup := router.Group("/api/sensors")
    {
        sensorGroup.POST("/", controller.CreateSensor)
        sensorGroup.GET("/:id", controller.GetSensor)
        sensorGroup.GET("/user/:user_id", controller.GetUserSensors)
        sensorGroup.PUT("/:id", controller.UpdateSensor)
        sensorGroup.DELETE("/:id", controller.DeleteSensor)
    }
}