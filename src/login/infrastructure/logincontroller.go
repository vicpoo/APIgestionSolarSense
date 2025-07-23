// api/src/login/infrastructure/logincontroller.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/login/application"

)

type LoginController struct {
    authHandlers   *application.AuthHandlers
    getAuthUseCase *application.GetAuthUseCase
    updateUseCase  *application.UpdateAuthUseCase
    deleteUseCase  *application.DeleteAuthUseCase
    getAuthHandler *GetAuthHandler // Nuevo
}

func NewLoginController(
    authHandlers *application.AuthHandlers,
    getUseCase *application.GetAuthUseCase,
    updateUseCase *application.UpdateAuthUseCase,
    deleteUseCase *application.DeleteAuthUseCase,
    getAuthHandler *GetAuthHandler, // Nuevo
) *LoginController {
    return &LoginController{
        authHandlers:   authHandlers,
        getAuthUseCase: getUseCase,
        updateUseCase:  updateUseCase,
        deleteUseCase:  deleteUseCase,
        getAuthHandler: getAuthHandler, // Nuevo
    }
}


// Handler para GET /api/auth/email
func (c *LoginController) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email parameter is required"})
		return
	}

	user, err := c.getAuthUseCase.GetUserByEmail(ctx.Request.Context(), email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Handler para PUT /api/auth/email
func (c *LoginController) UpdateUserEmail(ctx *gin.Context) {
	var updateRequest struct {
		CurrentEmail string `json:"current_email" binding:"required,email"`
		NewEmail     string `json:"new_email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	err := c.updateUseCase.UpdateUserEmail(ctx.Request.Context(), updateRequest.CurrentEmail, updateRequest.NewEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "email updated successfully"})
}

// Handler para DELETE /api/auth/email
func (c *LoginController) DeleteUserByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email parameter is required"})
		return
	}

	err := c.deleteUseCase.DeleteUserByEmail(ctx.Request.Context(), email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// Handlers existentes (registro, login, google auth)
func (c *LoginController) RegisterEmail(ctx *gin.Context) {
	c.authHandlers.RegisterEmail(ctx)
}

func (c *LoginController) LoginEmail(ctx *gin.Context) {
	c.authHandlers.LoginEmail(ctx)
}

func (c *LoginController) GoogleAuth(ctx *gin.Context) {
	c.authHandlers.GoogleAuth(ctx)
}

func (c *LoginController) GetAllUsers(ctx *gin.Context) {
    c.getAuthHandler.GetAllUsers(ctx)
}

func (c *LoginController) GetUserByID(ctx *gin.Context) {
    c.getAuthHandler.GetUserByID(ctx)
}