// src/notification_settings/infrastructure/post_settings_handler.go
package infrastructure

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/application"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/domain"
)

type PostSettingsHandler struct {
    useCase *application.PostSettingsUseCase
}

func NewPostSettingsHandler(useCase *application.PostSettingsUseCase) *PostSettingsHandler {
    return &PostSettingsHandler{useCase: useCase}
}

func (h *PostSettingsHandler) CreateSettings(c *gin.Context) {
    var settings domain.NotificationSettings
    if err := c.ShouldBindJSON(&settings); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.useCase.CreateSettings(c.Request.Context(), &settings); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusCreated)
}