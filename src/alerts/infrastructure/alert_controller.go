// src/alerts/infrastructure/alert_controller.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    
)

type AlertController struct {
    postHandler   *PostAlertHandler
    getHandler    *GetAlertHandler
    putHandler    *PutAlertHandler
    deleteHandler *DeleteAlertHandler
}

func NewAlertController(
    postHandler *PostAlertHandler,
    getHandler *GetAlertHandler,
    putHandler *PutAlertHandler,
    deleteHandler *DeleteAlertHandler,
) *AlertController {
    return &AlertController{
        postHandler:   postHandler,
        getHandler:    getHandler,
        putHandler:    putHandler,
        deleteHandler: deleteHandler,
    }
}
func (c *AlertController) CreateAlert(ctx *gin.Context) {
    c.postHandler.CreateAlert(ctx)
}

func (c *AlertController) GetAlert(ctx *gin.Context) {
    c.getHandler.GetAlert(ctx)
}

func (c *AlertController) GetSensorAlerts(ctx *gin.Context) {
    c.getHandler.GetSensorAlerts(ctx)
}

func (c *AlertController) GetUnsentAlerts(ctx *gin.Context) {
    c.getHandler.GetUnsentAlerts(ctx)
}

func (c *AlertController) MarkAlertAsSent(ctx *gin.Context) {
    c.postHandler.MarkAlertAsSent(ctx)
}

func (c *AlertController) UpdateAlert(ctx *gin.Context) {
    c.putHandler.UpdateAlert(ctx)
}

func (c *AlertController) DeleteAlert(ctx *gin.Context) {
    c.deleteHandler.DeleteAlert(ctx)
}