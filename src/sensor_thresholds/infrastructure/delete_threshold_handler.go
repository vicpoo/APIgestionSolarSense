// src/sensor_thresholds/infrastructure/delete_threshold_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/application"
)

type DeleteThresholdHandler struct {
    useCase *application.DeleteThresholdUseCase
}

func NewDeleteThresholdHandler(useCase *application.DeleteThresholdUseCase) *DeleteThresholdHandler {
    return &DeleteThresholdHandler{useCase: useCase}
}

func (h *DeleteThresholdHandler) DeleteThreshold(c *gin.Context) {
    sensorID, err := strconv.Atoi(c.Param("sensor_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }

    if err := h.useCase.DeleteThreshold(c.Request.Context(), sensorID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}