// api/src/login/infrastructure/logincontroller.go
package infrastructure

import (
    "log"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
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

func (c *LoginController) GetPublicUserInfo(ctx *gin.Context) {
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
            "id":          user.ID,
            "email":      user.Email,
            "display_name": user.Username,
            "username":    user.Username,
            "photo_url":   user.PhotoURL,
            "provider":    user.Provider,
            "auth_type":   user.AuthType,
            "created_at":  user.CreatedAt,
            "last_login":  user.LastLogin,
            "is_active":   user.IsActive,
        },
    })
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
        "email":     user.Email,
        "username":  user.Username,
        "auth_type": user.AuthType,
        "is_active": user.IsActive,
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

    user, passwordHash, err := c.getUseCase.Repo.FindUserByEmail(ctx.Request.Context(), creds.Email)
    if err != nil {
        log.Printf("Error de autenticación para %s: %v", creds.Email, err)
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password)); err != nil {
        log.Printf("Error en contraseña para usuario %d: %v", user.ID, err)
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    if err := c.getUseCase.Repo.UpdateLastLogin(ctx.Request.Context(), user.ID); err != nil {
        log.Printf("Error actualizando last_login para %d: %v", user.ID, err)
    }

    token, err := domain.GenerateJWTToken(user)
    if err != nil {
        log.Printf("Error generando token para %d: %v", user.ID, err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "token":   token,
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
            "type": "validation_error",
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

    var response []gin.H
    for _, user := range users {
        response = append(response, gin.H{
            "id":         user.ID,
            "email":      user.Email,
            "username":   user.Username,
            "auth_type":  user.AuthType,
            "is_active": user.IsActive,
            "created_at": user.CreatedAt,
            "last_login": user.LastLogin,
            "photo_url": user.PhotoURL,
        })
    }

    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "users":   response,
    })
}

func (c *LoginController) GetUserByID(ctx *gin.Context) {
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
            "is_active": user.IsActive,
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

    requestedUserID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
    if err == nil && requestedUserID != 0 && requestedUserID != userID.(int64) {
        ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own account"})
        return
    }

    userEmail, _ := ctx.Get("userEmail")
    
    err = c.deleteUseCase.DeleteUserByEmail(ctx.Request.Context(), userEmail.(string))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func decodeTokenWithoutVerification(idToken string) (map[string]interface{}, error) {
    return nil, nil
}