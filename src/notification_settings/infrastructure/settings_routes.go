// src/notification_settings/infrastructure/settings_routes.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/application"
)

func InitSettingsRoutes(router *gin.Engine) {
    repo := NewMySQLSettingsRepository()

    // Crear casos de uso
    getUseCase := application.NewGetSettingsUseCase(repo)
    postUseCase := application.NewPostSettingsUseCase(repo)
    putUseCase := application.NewPutSettingsUseCase(repo)
    deleteUseCase := application.NewDeleteSettingsUseCase(repo)

    // Crear handlers
    getHandler := NewGetSettingsHandler(getUseCase)
    postHandler := NewPostSettingsHandler(postUseCase)
    putHandler := NewPutSettingsHandler(putUseCase)
    deleteHandler := NewDeleteSettingsHandler(deleteUseCase)

    // Crear controlador
    controller := NewSettingsController(getHandler, postHandler, putHandler, deleteHandler)

    // Configurar rutas
    settingsGroup := router.Group("/api/notification-settings")
    {
        settingsGroup.GET("/user/:user_id", controller.GetSettings)
        settingsGroup.POST("/", controller.CreateSettings)
        settingsGroup.PUT("/user/:user_id", controller.UpdateSettings)
        settingsGroup.DELETE("/user/:user_id", controller.DeleteSettings)
    }
}