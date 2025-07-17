// src/sessions/infrastructure/post_session_handler.go
package infrastructure

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/application"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/domain"
)

type PostSessionHandler struct {
    useCase *application.PostSessionUseCase
}

func NewPostSessionHandler(useCase *application.PostSessionUseCase) *PostSessionHandler {
    return &PostSessionHandler{useCase: useCase}
}

func (h *PostSessionHandler) CreateSession(c *gin.Context) {
    var session domain.Session
    if err := c.ShouldBindJSON(&session); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.useCase.CreateSession(c.Request.Context(), &session); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, session)
}