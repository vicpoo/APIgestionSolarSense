// src/sensor_thresholds/infrastructure/get_threshold_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/application"
)

type GetThresholdHandler struct {
    useCase *application.GetThresholdUseCase
}

func NewGetThresholdHandler(useCase *application.GetThresholdUseCase) *GetThresholdHandler {
    return &GetThresholdHandler{useCase: useCase}
}

func (h *GetThresholdHandler) GetThresholds(c *gin.Context) {
    sensorID, err := strconv.Atoi(c.Param("sensor_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }

    thresholds, err := h.useCase.GetThresholds(c.Request.Context(), sensorID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if thresholds == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Thresholds not found"})
        return
    }

    c.JSON(http.StatusOK, thresholds)
}