// src/notification_settings/infrastructure/get_settings_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/application"
)

type GetSettingsHandler struct {
    useCase *application.GetSettingsUseCase
}

func NewGetSettingsHandler(useCase *application.GetSettingsUseCase) *GetSettingsHandler {
    return &GetSettingsHandler{useCase: useCase}
}

func (h *GetSettingsHandler) GetSettings(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    settings, err := h.useCase.GetSettings(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if settings == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Settings not found"})
        return
    }

    c.JSON(http.StatusOK, settings)
}