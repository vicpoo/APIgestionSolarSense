// src/alerts/infrastructure/delete_alert_handler.go
package infrastructure

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/application"
)

type DeleteAlertHandler struct {
    useCase *application.DeleteAlertUseCase
}

func NewDeleteAlertHandler(useCase *application.DeleteAlertUseCase) *DeleteAlertHandler {
    return &DeleteAlertHandler{useCase: useCase}
}

func (h *DeleteAlertHandler) DeleteAlert(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
        return
    }

    if err := h.useCase.DeleteAlert(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}