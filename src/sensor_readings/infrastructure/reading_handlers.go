//api/src/sensors/application/sensor_service.go

package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/application"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/domain"
)

type ReadingHandlers struct {
    service *application.ReadingService
}

func NewReadingHandlers(service *application.ReadingService) *ReadingHandlers {
    return &ReadingHandlers{service: service}
}

func (h *ReadingHandlers) CreateReading(c *gin.Context) {
    var reading domain.SensorReading
    if err := c.ShouldBindJSON(&reading); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.CreateReading(c.Request.Context(), &reading); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, reading)
}

func (h *ReadingHandlers) GetReadings(c *gin.Context) {
    sensorID, _ := strconv.Atoi(c.Param("sensor_id"))
    limit := 100
    if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
        limit = l
    }
    
    readings, err := h.service.GetReadings(c.Request.Context(), sensorID, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, readings)
}

func (h *ReadingHandlers) GetLatestReading(c *gin.Context) {
    sensorID, _ := strconv.Atoi(c.Param("sensor_id"))
    reading, err := h.service.GetLatestReading(c.Request.Context(), sensorID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, reading)
}