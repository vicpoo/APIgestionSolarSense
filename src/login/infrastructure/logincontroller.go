// api/src/login/infrastructure/logincontroller.go
package infrastructure

import (
	"log"
	"time"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vicpoo/apigestion-solar-go/src/core"
	"github.com/vicpoo/apigestion-solar-go/src/login/application"
	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
	"golang.org/x/crypto/bcrypt"
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
    claims, exists := ctx.Get("jwtClaims")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    jwtClaims, ok := claims.(*core.JWTClaims)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims"})
        return
    }

    user, err := c.getUseCase.GetUserByID(ctx.Request.Context(), jwtClaims.UserID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

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

    log.Printf("Intento de login con email: '%s'", creds.Email)

    // Buscar usuario
    user, passwordHash, err := c.getUseCase.Repo.FindUserByEmail(ctx.Request.Context(), creds.Email)
    if err != nil {
        log.Printf("Error al buscar usuario: %v", err)
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    log.Printf("Usuario encontrado. Comparando contraseñas...")

    // Verificar contraseña
    err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password))
    if err != nil {
        log.Printf("Error al comparar contraseñas: %v", err)
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Generar token
    claims := core.JWTClaims{
        UserID:   user.ID,
        Email:    user.Email,
        Username: user.Username,
        AuthType: user.AuthType,
        IsAdmin:  user.IsAdmin,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(core.JwtSecret))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "token":   tokenString,
        "user": gin.H{
            "id":       user.ID,
            "email":    user.Email,
            "username": user.Username,
            "is_admin": user.IsAdmin,
        },
    })
}
func (c *LoginController) RegisterEmail(ctx *gin.Context) {
    response, err := c.authHandlers.RegisterEmail(ctx)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
            "type": "validation_error", // Para mejor manejo en frontend
        })
        return
    }

    ctx.JSON(http.StatusCreated, response)
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
    users, err := c.getUseCase.GetAllUsers(ctx.Request.Context())
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users: " + err.Error()})
        return
    }

    // Formatear la respuesta
    var response []gin.H
    for _, user := range users {
        response = append(response, gin.H{
            "id":         user.ID,
            "email":      user.Email,
            "username":   user.Username,
            "auth_type":  user.AuthType,
            "is_active":  user.IsActive,
            "created_at": user.CreatedAt,
            "last_login": user.LastLogin,
            "photo_url":  user.PhotoURL,
        })
    }

    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "users":   response,
    })
}
func (c *LoginController) GetUserByID(ctx *gin.Context) {
    // Verificar permisos de admin
    isAdmin, exists := ctx.Get("isAdmin")
    if !exists || !isAdmin.(bool) {
        ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
        return
    }

    userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    user, err := c.getUseCase.GetUserByID(ctx.Request.Context(), userID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "user": gin.H{
            "id":         user.ID,
            "email":      user.Email,
            "username":   user.Username,
            "auth_type":  user.AuthType,
            "is_active":  user.IsActive,
            "created_at": user.CreatedAt,
            "last_login": user.LastLogin,
            "photo_url": user.PhotoURL,
        },
    })
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