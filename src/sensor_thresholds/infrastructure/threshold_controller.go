// src/sensor_thresholds/infrastructure/threshold_controller.go
package infrastructure

import "github.com/gin-gonic/gin"

type ThresholdController struct {
    getHandler    *GetThresholdHandler
    postHandler   *PostThresholdHandler
    putHandler    *PutThresholdHandler
    deleteHandler *DeleteThresholdHandler
}

func NewThresholdController(
    getHandler *GetThresholdHandler,
    postHandler *PostThresholdHandler,
    putHandler *PutThresholdHandler,
    deleteHandler *DeleteThresholdHandler,
) *ThresholdController {
    return &ThresholdController{
        getHandler:    getHandler,
        postHandler:   postHandler,
        putHandler:    putHandler,
        deleteHandler: deleteHandler,
    }
}

func (c *ThresholdController) GetThresholds(ctx *gin.Context) {
    c.getHandler.GetThresholds(ctx)
}

func (c *ThresholdController) CreateThreshold(ctx *gin.Context) {
    c.postHandler.CreateThreshold(ctx)
}

func (c *ThresholdController) UpdateThreshold(ctx *gin.Context) {
    c.putHandler.UpdateThreshold(ctx)
}

func (c *ThresholdController) DeleteThreshold(ctx *gin.Context) {
    c.deleteHandler.DeleteThreshold(ctx)
}