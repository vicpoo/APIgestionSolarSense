// src/sensor_readings/infrastructure/reading_controller.go
package infrastructure

import "github.com/gin-gonic/gin"

type ReadingController struct {
    getHandler    *GetReadingHandler
    postHandler   *PostReadingHandler
    putHandler    *PutReadingHandler
    deleteHandler *DeleteReadingHandler
}

func NewReadingController(
    getHandler *GetReadingHandler,
    postHandler *PostReadingHandler,
    putHandler *PutReadingHandler,
    deleteHandler *DeleteReadingHandler,
) *ReadingController {
    return &ReadingController{
        getHandler:    getHandler,
        postHandler:   postHandler,
        putHandler:    putHandler,
        deleteHandler: deleteHandler,
    }
}

func (c *ReadingController) CreateReading(ctx *gin.Context) {
    c.postHandler.CreateReading(ctx)
}

func (c *ReadingController) GetReadings(ctx *gin.Context) {
    c.getHandler.GetReadings(ctx)
}

func (c *ReadingController) GetLatestReading(ctx *gin.Context) {
    c.getHandler.GetLatestReading(ctx)
}

func (c *ReadingController) UpdateReading(ctx *gin.Context) {
    c.putHandler.UpdateReading(ctx)
}

func (c *ReadingController) DeleteReading(ctx *gin.Context) {
    c.deleteHandler.DeleteReading(ctx)
}