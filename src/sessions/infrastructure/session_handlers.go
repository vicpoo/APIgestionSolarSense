//api/src/sessions/infrastructure/session_handlers.go

package infrastructure

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/application"
   
)

type SessionHandlers struct {
    service *application.SessionService
}

func NewSessionHandlers(service *application.SessionService) *SessionHandlers {
    return &SessionHandlers{service: service}
}

func (h *SessionHandlers) ValidateSession(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing session token"})
        return
    }

    session, err := h.service.ValidateSession(c.Request.Context(), token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
        return
    }

    c.JSON(http.StatusOK, session)
}

func (h *SessionHandlers) Logout(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session token"})
        return
    }

    if err := h.service.InvalidateSession(c.Request.Context(), token); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}