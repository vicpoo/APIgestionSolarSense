// src/notification_settings/infrastructure/settings_controller.go
package infrastructure

import "github.com/gin-gonic/gin"

type SettingsController struct {
    getHandler    *GetSettingsHandler
    postHandler   *PostSettingsHandler
    putHandler    *PutSettingsHandler
    deleteHandler *DeleteSettingsHandler
}

func NewSettingsController(
    getHandler *GetSettingsHandler,
    postHandler *PostSettingsHandler,
    putHandler *PutSettingsHandler,
    deleteHandler *DeleteSettingsHandler,
) *SettingsController {
    return &SettingsController{
        getHandler:    getHandler,
        postHandler:   postHandler,
        putHandler:    putHandler,
        deleteHandler: deleteHandler,
    }
}

func (c *SettingsController) GetSettings(ctx *gin.Context) {
    c.getHandler.GetSettings(ctx)
}

func (c *SettingsController) CreateSettings(ctx *gin.Context) {
    c.postHandler.CreateSettings(ctx)
}

func (c *SettingsController) UpdateSettings(ctx *gin.Context) {
    c.putHandler.UpdateSettings(ctx)
}

func (c *SettingsController) DeleteSettings(ctx *gin.Context) {
    c.deleteHandler.DeleteSettings(ctx)
}