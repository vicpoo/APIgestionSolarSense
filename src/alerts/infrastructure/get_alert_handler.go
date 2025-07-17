// src/alerts/infrastructure/get_alert_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/application"
)

type GetAlertHandler struct {
    useCase *application.GetAlertUseCase
}

func NewGetAlertHandler(useCase *application.GetAlertUseCase) *GetAlertHandler {
    return &GetAlertHandler{useCase: useCase}
}

func (h *GetAlertHandler) GetAlert(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    alert, err := h.useCase.GetAlert(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
        return
    }
    c.JSON(http.StatusOK, alert)
}

func (h *GetAlertHandler) GetSensorAlerts(c *gin.Context) {
    sensorID, _ := strconv.Atoi(c.Param("sensor_id"))
    limit := 100
    if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
        limit = l
    }
    
    alerts, err := h.useCase.GetSensorAlerts(c.Request.Context(), sensorID, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, alerts)
}

func (h *GetAlertHandler) GetUnsentAlerts(c *gin.Context) {
    alerts, err := h.useCase.GetUnsentAlerts(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, alerts)
}