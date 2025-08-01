// src/memberships/infrastructure/put_membership_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/application"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
)

type PutMembershipHandler struct {
    useCase *application.PutMembershipUseCase
}

func NewPutMembershipHandler(useCase *application.PutMembershipUseCase) *PutMembershipHandler {
    return &PutMembershipHandler{useCase: useCase}
}

func (h *PutMembershipHandler) CreateOrUpdate(c *gin.Context) {
    // Obtener el user_id de la URL
    userID, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var membership domain.Membership
    if err := c.ShouldBindJSON(&membership); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    membership.UserID = userID

    // Forzar el tipo a premium
    membership.Type = "premium"
    membership.ExtraStorage = true

    if err := h.useCase.CreateOrUpdate(c, &membership); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}