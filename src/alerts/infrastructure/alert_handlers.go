//src/alerts/infrastructure/api/alert_handlers.go

package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/application"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/domain"
)

type AlertHandlers struct {
    service *application.AlertService
}

func NewAlertHandlers(service *application.AlertService) *AlertHandlers {
    return &AlertHandlers{service: service}
}

func (h *AlertHandlers) CreateAlert(c *gin.Context) {
    var alert domain.Alert
    if err := c.ShouldBindJSON(&alert); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.CreateAlert(c.Request.Context(), &alert); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, alert)
}

func (h *AlertHandlers) GetAlert(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    alert, err := h.service.GetAlert(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
        return
    }
    c.JSON(http.StatusOK, alert)
}

func (h *AlertHandlers) GetSensorAlerts(c *gin.Context) {
    sensorID, _ := strconv.Atoi(c.Param("sensor_id"))
    limit := 100
    if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
        limit = l
    }
    
    alerts, err := h.service.GetSensorAlerts(c.Request.Context(), sensorID, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, alerts)
}

func (h *AlertHandlers) GetUnsentAlerts(c *gin.Context) {
    alerts, err := h.service.GetUnsentAlerts(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, alerts)
}

func (h *AlertHandlers) MarkAlertAsSent(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.service.MarkAlertAsSent(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}