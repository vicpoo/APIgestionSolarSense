// src/sessions/infrastructure/session_controller.go
package infrastructure

import "github.com/gin-gonic/gin"

type SessionController struct {
    getHandler    *GetSessionHandler
    postHandler   *PostSessionHandler
    putHandler    *PutSessionHandler
    deleteHandler *DeleteSessionHandler
}

func NewSessionController(
    getHandler *GetSessionHandler,
    postHandler *PostSessionHandler,
    putHandler *PutSessionHandler,
    deleteHandler *DeleteSessionHandler,
) *SessionController {
    return &SessionController{
        getHandler:    getHandler,
        postHandler:   postHandler,
        putHandler:    putHandler,
        deleteHandler: deleteHandler,
    }
}

func (c *SessionController) ValidateSession(ctx *gin.Context) {
    c.getHandler.ValidateSession(ctx)
}

func (c *SessionController) CreateSession(ctx *gin.Context) {
    c.postHandler.CreateSession(ctx)
}

func (c *SessionController) RefreshSession(ctx *gin.Context) {
    c.putHandler.RefreshSession(ctx)
}

func (c *SessionController) InvalidateSession(ctx *gin.Context) {
    c.deleteHandler.InvalidateSession(ctx)
}

func (c *SessionController) DeleteSession(ctx *gin.Context) {
    c.deleteHandler.DeleteSession(ctx)
}