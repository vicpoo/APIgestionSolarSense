// src/memberships/infrastructure/delete_membership_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/application"
)

type DeleteMembershipHandler struct {
    useCase *application.DeleteMembershipUseCase
}

func NewDeleteMembershipHandler(useCase *application.DeleteMembershipUseCase) *DeleteMembershipHandler {
    return &DeleteMembershipHandler{useCase: useCase}
}

func (h *DeleteMembershipHandler) DeleteMembership(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    if err := h.useCase.DeleteMembership(c.Request.Context(), userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}