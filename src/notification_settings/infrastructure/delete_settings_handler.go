// src/notification_settings/infrastructure/delete_settings_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/application"
)

type DeleteSettingsHandler struct {
    useCase *application.DeleteSettingsUseCase
}

func NewDeleteSettingsHandler(useCase *application.DeleteSettingsUseCase) *DeleteSettingsHandler {
    return &DeleteSettingsHandler{useCase: useCase}
}

func (h *DeleteSettingsHandler) DeleteSettings(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    if err := h.useCase.DeleteSettings(c.Request.Context(), userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}