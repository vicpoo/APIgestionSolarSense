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
    // Verificar si el usuario es admin
    claims, exists := c.Get("userClaims")
    if !exists {
        c.JSON(http.StatusForbidden, gin.H{"error": "Authentication required"})
        return
    }
    
    claimsMap, ok := claims.(map[string]interface{})
    if !ok {
        c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user claims"})
        return
    }
    
    role, ok := claimsMap["role"].(string)
    if !ok || role != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can modify memberships"})
        return
    }

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

    // Pasar el contexto Gin directamente
    if err := h.useCase.CreateOrUpdate(c, &membership); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}