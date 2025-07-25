// api/src/login/infrastructure/get_auth_handler.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/login/application"
)

type GetAuthHandler struct {
	useCase *application.GetAuthUseCase
}

func NewGetAuthHandler(useCase *application.GetAuthUseCase) *GetAuthHandler {
	return &GetAuthHandler{useCase: useCase}
}

func (h *GetAuthHandler) GetAllUsers(c *gin.Context) {
	users, err := h.useCase.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *GetAuthHandler) GetUserByID(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.useCase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
