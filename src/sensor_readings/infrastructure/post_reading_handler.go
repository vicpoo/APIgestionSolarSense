// src/sensor_readings/infrastructure/post_reading_handler.go
package infrastructure

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/application"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/domain"
)

type PostReadingHandler struct {
    useCase *application.PostReadingUseCase
}

func NewPostReadingHandler(useCase *application.PostReadingUseCase) *PostReadingHandler {
    return &PostReadingHandler{useCase: useCase}
}

func (h *PostReadingHandler) CreateReading(c *gin.Context) {
    var reading domain.SensorReading
    if err := c.ShouldBindJSON(&reading); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.useCase.CreateReading(c.Request.Context(), &reading); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, reading)
}