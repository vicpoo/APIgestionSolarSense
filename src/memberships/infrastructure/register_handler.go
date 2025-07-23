// src/memberships/infrastructure/register_handler.go
package infrastructure

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/application"
)

type RegisterHandler struct {
    useCase *application.RegisterUseCase
}

func NewRegisterHandler(useCase *application.RegisterUseCase) *RegisterHandler {
    return &RegisterHandler{useCase: useCase}
}

func (h *RegisterHandler) RegisterUser(c *gin.Context) {
    var request struct {
        Email    string `json:"email" binding:"required,email"`
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, err := h.useCase.RegisterUser(c.Request.Context(), request.Email, request.Username, request.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "user_id": userID,
        "message": "User registered successfully",
    })
}