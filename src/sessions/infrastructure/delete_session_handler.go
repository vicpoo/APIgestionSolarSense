// src/sessions/infrastructure/delete_session_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/application"
)

type DeleteSessionHandler struct {
    useCase *application.DeleteSessionUseCase
}

func NewDeleteSessionHandler(useCase *application.DeleteSessionUseCase) *DeleteSessionHandler {
    return &DeleteSessionHandler{useCase: useCase}
}

func (h *DeleteSessionHandler) InvalidateSession(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session token"})
        return
    }

    if err := h.useCase.InvalidateSession(c.Request.Context(), token); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

func (h *DeleteSessionHandler) DeleteSession(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
        return
    }

    if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}