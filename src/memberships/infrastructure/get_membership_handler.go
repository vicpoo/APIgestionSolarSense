// src/memberships/infrastructure/get_membership_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/application"
)

type GetMembershipHandler struct {
    useCase *application.GetMembershipUseCase
}

func NewGetMembershipHandler(useCase *application.GetMembershipUseCase) *GetMembershipHandler {
    return &GetMembershipHandler{useCase: useCase}
}

func (h *GetMembershipHandler) GetUserMembership(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    membership, err := h.useCase.GetUserMembership(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if membership == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Membership not found"})
        return
    }

    c.JSON(http.StatusOK, membership)
}

// Nuevo método para obtener todos los usuarios
func (h *GetMembershipHandler) GetAllUsers(c *gin.Context) {
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
        c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can list all users"})
        return
    }

    users, err := h.useCase.GetAllUsers(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, users)
}