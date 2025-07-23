package infrastructure

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/memberships/application"
)

type UpdateUserHandler struct {
	useCase *application.UpdateUserUseCase
}

func NewUpdateUserHandler(useCase *application.UpdateUserUseCase) *UpdateUserHandler {
	return &UpdateUserHandler{useCase: useCase}
}

func (h *UpdateUserHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var request struct {
		Email    *string `json:"email,omitempty"`
		Username *string `json:"username,omitempty"`
		Password *string `json:"password,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCase.UpdateUser(c.Request.Context(), userID, request.Email, request.Username, request.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}