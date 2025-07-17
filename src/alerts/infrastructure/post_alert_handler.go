// src/alerts/infrastructure/post_alert_handler.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/alerts/application"
	"github.com/vicpoo/apigestion-solar-go/src/alerts/domain"
)

type PostAlertHandler struct {
    useCase *application.PostAlertUseCase
}

func NewPostAlertHandler(useCase *application.PostAlertUseCase) *PostAlertHandler {
    return &PostAlertHandler{useCase: useCase}
}

func (h *PostAlertHandler) CreateAlert(c *gin.Context) {
    var alert domain.Alert
    if err := c.ShouldBindJSON(&alert); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.useCase.CreateAlert(c.Request.Context(), &alert); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, alert)
}

func (h *PostAlertHandler) MarkAlertAsSent(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.useCase.MarkAlertAsSent(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}