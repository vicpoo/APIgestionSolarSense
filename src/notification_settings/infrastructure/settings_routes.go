//api\src\notification_settings\infrastructure\settings_handlers.go

package infrastructure

import (
    "github.com/gin-gonic/gin"

    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/application"

)

func InitSettingsRoutes(router *gin.Engine) {
    repo := NewMySQLSettingsRepository()
    service := application.NewSettingsService(repo)
    handlers := NewSettingsHandlers(service)

    settingsGroup := router.Group("/api/notification-settings")
    {
        settingsGroup.GET("/user/:user_id", handlers.GetSettings)
        settingsGroup.PUT("/user/:user_id", handlers.UpdateSettings)
    }
}