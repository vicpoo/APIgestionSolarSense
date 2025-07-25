// src/memberships/infrastructure/register_handler.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/memberships/application"
	"github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
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
        Username string `json:"username" binding:"required,min=3"`
        Password string `json:"password" binding:"required,min=8"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    userID, err := h.useCase.RegisterUser(c.Request.Context(), request.Email, request.Username, request.Password)
    if err != nil {
        status := http.StatusInternalServerError
        if err == domain.ErrEmailAlreadyExists {
            status = http.StatusConflict
        }
        c.JSON(status, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "user_id": userID,
        "message": "User registered successfully",
    })
}