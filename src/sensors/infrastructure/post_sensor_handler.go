// src/sensors/infrastructure/post_sensor_handler.go
package infrastructure

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/application"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/domain"
)

type PostSensorHandler struct {
    useCase *application.PostSensorUseCase
}

func NewPostSensorHandler(useCase *application.PostSensorUseCase) *PostSensorHandler {
    return &PostSensorHandler{useCase: useCase}
}

func (h *PostSensorHandler) CreateSensor(c *gin.Context) {
    var sensor domain.Sensor
    if err := c.ShouldBindJSON(&sensor); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.useCase.CreateSensor(c.Request.Context(), &sensor); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, sensor)
}