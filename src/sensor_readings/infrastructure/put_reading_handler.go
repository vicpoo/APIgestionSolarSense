// src/sensor_readings/infrastructure/put_reading_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/application"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/domain"
)

type PutReadingHandler struct {
    useCase *application.PutReadingUseCase
}

func NewPutReadingHandler(useCase *application.PutReadingUseCase) *PutReadingHandler {
    return &PutReadingHandler{useCase: useCase}
}

func (h *PutReadingHandler) UpdateReading(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading ID"})
        return
    }

    var reading domain.SensorReading
    if err := c.ShouldBindJSON(&reading); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Asegurar que el ID de la URL coincide con el del body
    reading.ID = id

    if err := h.useCase.UpdateReading(c.Request.Context(), &reading); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}