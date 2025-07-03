//api/src/login/application/auth_handlers.go

package application

import (
	
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
)

type AuthHandlers struct {
	service domain.AuthService
}

func NewAuthHandlers(service domain.AuthService) *AuthHandlers {
	return &AuthHandlers{service: service}
}

func (h *AuthHandlers) RegisterEmail(c *gin.Context) {
	var creds domain.UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	response, err := h.service.RegisterWithEmail(c.Request.Context(), creds)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandlers) LoginEmail(c *gin.Context) {
	var creds domain.UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	response, err := h.service.LoginWithEmail(c.Request.Context(), creds)
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandlers) GoogleAuth(c *gin.Context) {
	var request struct {
		IDToken string `json:"idToken"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	response, err := h.service.AuthenticateWithGoogle(c.Request.Context(), request.IDToken)
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, domain.AuthResponse{Error: message})
}