// src/memberships/infrastructure/post_membership_handler.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/memberships/application"
)

type PostMembershipHandler struct {
	useCase *application.PostMembershipUseCase
}

func NewPostMembershipHandler(useCase *application.PostMembershipUseCase) *PostMembershipHandler {
	return &PostMembershipHandler{useCase: useCase}
}

func (h *PostMembershipHandler) UpgradeToPremium(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.useCase.UpgradeToPremium(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PostMembershipHandler) DowngradeToFree(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.useCase.DowngradeToFree(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Nuevo método para manejar la corrección de membresías faltantes
func (h *PostMembershipHandler) FixMissingMemberships(c *gin.Context) {
	if err := h.useCase.FixMissingMemberships(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Missing memberships fixed"})
}