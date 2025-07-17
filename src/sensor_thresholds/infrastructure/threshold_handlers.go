// api/src/sensor_thresholds/infrastructure/threshold_handlers.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/application"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/domain"
)

type ThresholdHandlers struct {
    service *application.ThresholdService
}

func NewThresholdHandlers(service *application.ThresholdService) *ThresholdHandlers {
    return &ThresholdHandlers{service: service}
}

func (h *ThresholdHandlers) GetThresholds(c *gin.Context) {
    sensorID, err := strconv.Atoi(c.Param("sensor_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }

    thresholds, err := h.service.GetThresholds(c.Request.Context(), sensorID)
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

func (h *ThresholdHandlers) CreateThreshold(c *gin.Context) {
    var threshold domain.SensorThreshold
    if err := c.ShouldBindJSON(&threshold); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.CreateThreshold(c.Request.Context(), &threshold); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, threshold)
}

func (h *ThresholdHandlers) UpdateThreshold(c *gin.Context) {
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

    if err := h.service.UpdateThreshold(c.Request.Context(), &threshold); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

func (h *ThresholdHandlers) DeleteThreshold(c *gin.Context) {
    sensorID, err := strconv.Atoi(c.Param("sensor_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }

    if err := h.service.DeleteThreshold(c.Request.Context(), sensorID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}