// src/alerts/infrastructure/put_alert_handler.go
package infrastructure

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/application"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/domain"
)

type PutAlertHandler struct {
    useCase *application.PutAlertUseCase
}

func NewPutAlertHandler(useCase *application.PutAlertUseCase) *PutAlertHandler {
    return &PutAlertHandler{useCase: useCase}
}

func (h *PutAlertHandler) UpdateAlert(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
        return
    }

    var alert domain.Alert
    if err := c.ShouldBindJSON(&alert); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Asegurarnos que el ID de la URL coincide con el ID del body
    alert.ID = id

    if err := h.useCase.UpdateAlert(c.Request.Context(), &alert); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, alert)
}