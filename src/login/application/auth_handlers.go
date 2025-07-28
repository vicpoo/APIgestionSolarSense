// api/src/login/application/auth_handlers.go
package application

import (
	
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
)

type AuthHandlers struct {
	service domain.AuthService
}

func NewAuthHandlers(service domain.AuthService) *AuthHandlers {
	return &AuthHandlers{service: service}
}

func (h *AuthHandlers) RegisterEmail(c *gin.Context) (*domain.AuthResponse, error) {
    // 1. Verificar el Content-Type
    if c.ContentType() != "application/json" {
        return nil, errors.New("content-type must be application/json")
    }

    // 2. Leer el cuerpo manualmente
    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading request body: %v", err)
    }
    defer c.Request.Body.Close()

    // 3. Parsear manualmente el JSON
    var creds domain.UserCredentials
    if err := json.Unmarshal(body, &creds); err != nil {
        return nil, fmt.Errorf("invalid JSON format: %v", err)
    }

    // 4. Validaciones manuales
    if creds.Email == "" {
        return nil, errors.New("email is required")
    }
    if creds.Password == "" {
        return nil, errors.New("password is required")
    }
    if creds.Username == "" {
        return nil, errors.New("username is required")
    }

    // 5. Llamar al servicio
    return h.service.RegisterWithEmail(c.Request.Context(), creds)
}
func (h *AuthHandlers) LoginEmail(c *gin.Context) (*domain.AuthResponse, error) {
	var creds domain.UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		return nil, errors.New("invalid request payload")
	}

	response, err := h.service.LoginWithEmail(c.Request.Context(), creds)
	if err != nil {
		return nil, err
	}

	return response, nil
}
// En auth_handlers.go
func (h *AuthHandlers) GoogleAuth(c *gin.Context) (*domain.AuthResponse, error) {
    var request struct {
        IDToken string `json:"idToken"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        return nil, errors.New("invalid request payload")
    }

    response, err := h.service.AuthenticateWithGoogle(c.Request.Context(), request.IDToken)
    if err != nil {
        return nil, err
    }

    return response, nil
}
func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, domain.AuthResponse{Error: message})
}

func (h *AuthHandlers) UpdateUserProfile(c *gin.Context) (*domain.AuthResponse, error) {
    // 1. Verificar el Content-Type
    if c.ContentType() != "application/json" {
        return nil, errors.New("content-type must be application/json")
    }

    // 2. Leer el cuerpo manualmente
    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading request body: %v", err)
    }
    defer c.Request.Body.Close()

    // 3. Parsear manualmente el JSON
    var updateData struct {
        Email       string `json:"email,omitempty"`
        Username    string `json:"username,omitempty"`
        DisplayName string `json:"display_name,omitempty"`
        AuthType    string `json:"auth_type"` // Necesario para saber qué campos actualizar
        UserID      int64  `json:"user_id"`   // Necesario para identificar al usuario
    }
    
    if err := json.Unmarshal(body, &updateData); err != nil {
        return nil, fmt.Errorf("invalid JSON format: %v", err)
    }

    // 4. Validar datos según el tipo de autenticación
    if updateData.AuthType == "email" {
        if updateData.Email == "" && updateData.Username == "" && updateData.DisplayName == "" {
            return nil, errors.New("at least one field must be provided for email users (email, username or display_name)")
        }
    } else if updateData.AuthType == "google" {
        if updateData.DisplayName == "" {
            return nil, errors.New("display_name is required for google users")
        }
        // Limpiar otros campos que no deberían actualizarse
        updateData.Email = ""
        updateData.Username = ""
    } else {
        return nil, errors.New("invalid auth_type")
    }

    // 5. Llamar al servicio
    return h.service.UpdateUserProfile(c.Request.Context(), updateData.UserID, updateData.Email, updateData.Username, updateData.DisplayName, updateData.AuthType)
}