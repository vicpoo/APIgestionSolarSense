// src/sensor_readings/infrastructure/get_reading_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/application"
)

type GetReadingHandler struct {
    useCase *application.GetReadingUseCase
}

func NewGetReadingHandler(useCase *application.GetReadingUseCase) *GetReadingHandler {
    return &GetReadingHandler{useCase: useCase}
}

func (h *GetReadingHandler) GetReadings(c *gin.Context) {
    sensorID, err := strconv.Atoi(c.Param("sensor_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }

    limit := 100
    if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
        limit = l
    }

    readings, err := h.useCase.GetReadings(c.Request.Context(), sensorID, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, readings)
}

func (h *GetReadingHandler) GetLatestReading(c *gin.Context) {
    sensorID, err := strconv.Atoi(c.Param("sensor_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }

    reading, err := h.useCase.GetLatestReading(c.Request.Context(), sensorID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, reading)
}