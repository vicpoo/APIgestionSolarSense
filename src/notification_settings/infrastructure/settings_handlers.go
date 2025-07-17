//api\src\notification_settings/infrastructure/settings_handlers.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/application"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/domain"
)

type SettingsHandlers struct {
    service *application.SettingsService
}

func NewSettingsHandlers(service *application.SettingsService) *SettingsHandlers {
    return &SettingsHandlers{service: service}
}

func (h *SettingsHandlers) GetSettings(c *gin.Context) {
    userID, _ := strconv.Atoi(c.Param("user_id"))
    settings, err := h.service.GetSettings(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, settings)
}

func (h *SettingsHandlers) UpdateSettings(c *gin.Context) {
    userID, _ := strconv.Atoi(c.Param("user_id"))
    
    var settings domain.NotificationSettings
    if err := c.ShouldBindJSON(&settings); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    settings.UserID = userID
    
    if err := h.service.UpdateSettings(c.Request.Context(), &settings); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}