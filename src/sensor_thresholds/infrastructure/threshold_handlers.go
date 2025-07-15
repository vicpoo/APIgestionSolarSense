//api\src\sensor_thresholds\domain\threshold_repository.go

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
    sensorID, _ := strconv.Atoi(c.Param("sensor_id"))
    thresholds, err := h.service.GetThresholds(c.Request.Context(), sensorID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, thresholds)
}

func (h *ThresholdHandlers) SetThresholds(c *gin.Context) {
    sensorID, _ := strconv.Atoi(c.Param("sensor_id"))
    
    var thresholds domain.SensorThreshold
    if err := c.ShouldBindJSON(&thresholds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    thresholds.SensorID = sensorID
    
    if err := h.service.SetThresholds(c.Request.Context(), &thresholds); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}