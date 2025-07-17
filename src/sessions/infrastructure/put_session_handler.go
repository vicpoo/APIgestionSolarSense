// src/sessions/infrastructure/put_session_handler.go
package infrastructure

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/application"
)

type PutSessionHandler struct {
    useCase *application.PutSessionUseCase
}

func NewPutSessionHandler(useCase *application.PutSessionUseCase) *PutSessionHandler {
    return &PutSessionHandler{useCase: useCase}
}

func (h *PutSessionHandler) RefreshSession(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session token"})
        return
    }

    // Extender la sesión por 1 hora más
    newExpiry := time.Now().Add(time.Hour * 1)
    
    if err := h.useCase.RefreshSession(c.Request.Context(), token, newExpiry); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}