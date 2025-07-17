// src/notification_settings/infrastructure/put_settings_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/application"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/domain"
)

type PutSettingsHandler struct {
    useCase *application.PutSettingsUseCase
}

func NewPutSettingsHandler(useCase *application.PutSettingsUseCase) *PutSettingsHandler {
    return &PutSettingsHandler{useCase: useCase}
}

func (h *PutSettingsHandler) UpdateSettings(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var settings domain.NotificationSettings
    if err := c.ShouldBindJSON(&settings); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    settings.UserID = userID

    if err := h.useCase.UpdateSettings(c.Request.Context(), &settings); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}