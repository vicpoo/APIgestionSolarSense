// src/sensor_readings/infrastructure/reading_routes.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/application"
)

func InitReadingRoutes(router *gin.Engine) {
    repo := NewMySQLReadingRepository()

    // Crear casos de uso
    getUseCase := application.NewGetReadingUseCase(repo)
    postUseCase := application.NewPostReadingUseCase(repo)
    putUseCase := application.NewPutReadingUseCase(repo)
    deleteUseCase := application.NewDeleteReadingUseCase(repo)

    // Crear handlers
    getHandler := NewGetReadingHandler(getUseCase)
    postHandler := NewPostReadingHandler(postUseCase)
    putHandler := NewPutReadingHandler(putUseCase)
    deleteHandler := NewDeleteReadingHandler(deleteUseCase)

    // Crear controlador
    controller := NewReadingController(getHandler, postHandler, putHandler, deleteHandler)

    // Configurar rutas
    readingGroup := router.Group("/api/readings")
    {
        readingGroup.POST("/", controller.CreateReading)
        readingGroup.GET("/sensor/:sensor_id", controller.GetReadings)
        readingGroup.GET("/sensor/:sensor_id/latest", controller.GetLatestReading)
        readingGroup.PUT("/:id", controller.UpdateReading)
        readingGroup.DELETE("/:id", controller.DeleteReading)
    }
}