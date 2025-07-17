// src/sensor_thresholds/infrastructure/put_threshold_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/application"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/domain"
)

type PutThresholdHandler struct {
    useCase *application.PutThresholdUseCase
}

func NewPutThresholdHandler(useCase *application.PutThresholdUseCase) *PutThresholdHandler {
    return &PutThresholdHandler{useCase: useCase}
}

func (h *PutThresholdHandler) UpdateThreshold(c *gin.Context) {
    sensorID, err := strconv.Atoi(c.Param("sensor_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }

    var threshold domain.SensorThreshold
    if err := c.ShouldBindJSON(&threshold); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    threshold.SensorID = sensorID

    if err := h.useCase.UpdateThreshold(c.Request.Context(), &threshold); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}