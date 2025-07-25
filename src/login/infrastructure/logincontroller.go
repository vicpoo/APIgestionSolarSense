// api/src/login/infrastructure/logincontroller.go
package infrastructure

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/login/application"
	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
)

type LoginController struct {
	authHandlers   *application.AuthHandlers
	getUseCase     *application.GetAuthUseCase
	updateUseCase  *application.UpdateAuthUseCase
	deleteUseCase  *application.DeleteAuthUseCase
	getAuthHandler *GetAuthHandler
}

func NewLoginController(
	authHandlers *application.AuthHandlers,
	getUseCase *application.GetAuthUseCase,
	updateUseCase *application.UpdateAuthUseCase,
	deleteUseCase *application.DeleteAuthUseCase,
	getAuthHandler *GetAuthHandler,
) *LoginController {
	return &LoginController{
		authHandlers:   authHandlers,
		getUseCase:     getUseCase,
		updateUseCase:  updateUseCase,
		deleteUseCase:  deleteUseCase,
		getAuthHandler: getAuthHandler,
	}
}

func (c *LoginController) GetCurrentUser(ctx *gin.Context) {
    // Obtener claims del token
    claims, exists := ctx.Get("jwtClaims")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    jwtClaims, ok := claims.(*JWTClaims)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims"})
        return
    }

    // Obtener datos completos del usuario desde la base de datos
    user, err := c.getUseCase.GetUserByID(ctx.Request.Context(), jwtClaims.UserID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Devolver todos los datos del usuario
    ctx.JSON(http.StatusOK, gin.H{
        "user_id":    user.ID,
        "email":      user.Email,
        "username":   user.Username,
        "auth_type":  user.AuthType,
        "is_active":  user.IsActive,
        "last_login": user.LastLogin,
        "created_at": user.CreatedAt,
    })
}

func (c *LoginController) UpdateUserEmail(ctx *gin.Context) {
	userEmail, exists := ctx.Get("userEmail")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var updateRequest struct {
		NewEmail string `json:"new_email" binding:"required,email"`
	}
	
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	
	err := c.updateUseCase.UpdateUserEmail(ctx.Request.Context(), userEmail.(string), updateRequest.NewEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "Email updated successfully"})
}

func (c *LoginController) UpdatePassword(ctx *gin.Context) {
	userEmail, exists := ctx.Get("userEmail")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var updateRequest struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}
	
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	
	user, err := c.getUseCase.GetUserByEmail(ctx.Request.Context(), userEmail.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	err = c.updateUseCase.UpdatePassword(ctx.Request.Context(), user.ID, updateRequest.CurrentPassword, updateRequest.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func (c *LoginController) LoginEmail(ctx *gin.Context) {
    var creds domain.UserCredentials
    if err := ctx.ShouldBindJSON(&creds); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    // Autenticar al usuario - ahora usamos el response directamente en la respuesta
    authResponse, err := c.authHandlers.LoginEmail(ctx)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Obtener datos completos del usuario
    user, err := c.getUseCase.GetUserByEmail(ctx.Request.Context(), creds.Email)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get user data"})
        return
    }

    // Generar token con duración de 48 horas
    token, err := GenerateJWTToken(user)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    // Establecer la sesión
    session := sessions.Default(ctx)
    session.Set("userID", user.ID)
    if err := session.Save(); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
        return
    }

    // Mostrar en consola
    fmt.Printf("Usuario logueado - ID: %d, Email: %s\n", user.ID, user.Email)

    // Combinamos la respuesta de authHandlers con nuestros datos adicionales
    ctx.JSON(http.StatusOK, gin.H{
        "success":    authResponse.Success,
        "message":    authResponse.Message,
        "token":      token,
        "user_id":    user.ID,
        "email":      user.Email,
        "username":   user.Username,
        "auth_type":  user.AuthType,
        "expires_in": 48 * 3600, // 48 horas en segundos
    })
}
func (c *LoginController) RegisterEmail(ctx *gin.Context) {
	c.authHandlers.RegisterEmail(ctx)
}

func (c *LoginController) GoogleAuth(ctx *gin.Context) {
    response, err := c.authHandlers.GoogleAuth(ctx)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    var request struct {
        IDToken string `json:"idToken"`
    }
    if err := ctx.ShouldBindJSON(&request); err == nil {
        userData, err := decodeTokenWithoutVerification(request.IDToken)
        if err == nil {
            user, err := c.getUseCase.GetUserByEmail(ctx.Request.Context(), userData["email"].(string))
            if err == nil {
                membershipType, err := c.getUseCase.GetUserMembershipType(ctx.Request.Context(), user.ID)
                if err != nil {
                    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get user membership"})
                    return
                }

                token, err := GenerateJWTToken(user)
                if err != nil {
                    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
                    return
                }
                response.Token = token
                response.AuthType = "google"
                response.UserID = user.ID
                response.IsAdmin = membershipType == "admin"
            }
        }
    }
    
    ctx.JSON(http.StatusOK, response)
}

func (c *LoginController) GetAllUsers(ctx *gin.Context) {
	c.getAuthHandler.GetAllUsers(ctx)
}

func (c *LoginController) GetUserByID(ctx *gin.Context) {
	c.getAuthHandler.GetUserByID(ctx)
}

func (c *LoginController) DeleteAccount(ctx *gin.Context) {
    userID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Verificar que el usuario solo se está eliminando a sí mismo
    requestedUserID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
    if err == nil && requestedUserID != 0 && requestedUserID != userID.(int64) {
        ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own account"})
        return
    }

    // Obtener email para eliminar
    userEmail, _ := ctx.Get("userEmail")
    
    err = c.deleteUseCase.DeleteUserByEmail(ctx.Request.Context(), userEmail.(string))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// Función auxiliar para decodificar tokens
func decodeTokenWithoutVerification(idToken string) (map[string]interface{}, error) {
	return nil, nil
}