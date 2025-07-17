// src/sensors/infrastructure/sensor_controller.go
package infrastructure

import "github.com/gin-gonic/gin"

type SensorController struct {
    getHandler    *GetSensorHandler
    postHandler   *PostSensorHandler
    putHandler    *PutSensorHandler
    deleteHandler *DeleteSensorHandler
}

func NewSensorController(
    getHandler *GetSensorHandler,
    postHandler *PostSensorHandler,
    putHandler *PutSensorHandler,
    deleteHandler *DeleteSensorHandler,
) *SensorController {
    return &SensorController{
        getHandler:    getHandler,
        postHandler:   postHandler,
        putHandler:    putHandler,
        deleteHandler: deleteHandler,
    }
}

func (c *SensorController) CreateSensor(ctx *gin.Context) {
    c.postHandler.CreateSensor(ctx)
}

func (c *SensorController) GetSensor(ctx *gin.Context) {
    c.getHandler.GetSensor(ctx)
}

func (c *SensorController) GetUserSensors(ctx *gin.Context) {
    c.getHandler.GetUserSensors(ctx)
}

func (c *SensorController) UpdateSensor(ctx *gin.Context) {
    c.putHandler.UpdateSensor(ctx)
}

func (c *SensorController) DeleteSensor(ctx *gin.Context) {
    c.deleteHandler.DeleteSensor(ctx)
}