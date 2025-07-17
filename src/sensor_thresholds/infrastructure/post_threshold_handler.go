// src/sensor_thresholds/infrastructure/post_threshold_handler.go
package infrastructure

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/application"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/domain"
)

type PostThresholdHandler struct {
    useCase *application.PostThresholdUseCase
}

func NewPostThresholdHandler(useCase *application.PostThresholdUseCase) *PostThresholdHandler {
    return &PostThresholdHandler{useCase: useCase}
}

func (h *PostThresholdHandler) CreateThreshold(c *gin.Context) {
    var threshold domain.SensorThreshold
    if err := c.ShouldBindJSON(&threshold); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.useCase.CreateThreshold(c.Request.Context(), &threshold); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, threshold)
}