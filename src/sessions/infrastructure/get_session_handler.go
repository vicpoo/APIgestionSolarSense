// src/sessions/infrastructure/get_session_handler.go
package infrastructure

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/application"
)

type GetSessionHandler struct {
    useCase *application.GetSessionUseCase
}

func NewGetSessionHandler(useCase *application.GetSessionUseCase) *GetSessionHandler {
    return &GetSessionHandler{useCase: useCase}
}

func (h *GetSessionHandler) ValidateSession(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing session token"})
        return
    }

    session, err := h.useCase.ValidateSession(c.Request.Context(), token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
        return
    }

    c.JSON(http.StatusOK, session)
}