// src/sensors/infrastructure/put_sensor_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/application"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/domain"
)

type PutSensorHandler struct {
    useCase *application.PutSensorUseCase
}

func NewPutSensorHandler(useCase *application.PutSensorUseCase) *PutSensorHandler {
    return &PutSensorHandler{useCase: useCase}
}

func (h *PutSensorHandler) UpdateSensor(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }

    var sensor domain.Sensor
    if err := c.ShouldBindJSON(&sensor); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    sensor.ID = id

    if err := h.useCase.UpdateSensor(c.Request.Context(), &sensor); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, sensor)
}