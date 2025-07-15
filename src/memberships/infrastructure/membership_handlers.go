//src/memberships/infrastructure/api/membership_handlers.go

package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/application"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
)

type MembershipHandlers struct {
    service *application.MembershipService
}

func NewMembershipHandlers(service *application.MembershipService) *MembershipHandlers {
    return &MembershipHandlers{service: service}
}

func (h *MembershipHandlers) GetUserMembership(c *gin.Context) {
    userID, _ := strconv.Atoi(c.Param("user_id"))
    membership, err := h.service.GetUserMembership(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, membership)
}

func (h *MembershipHandlers) UpdateMembership(c *gin.Context) {
    userID, _ := strconv.Atoi(c.Param("user_id"))
    
    var membership domain.Membership
    if err := c.ShouldBindJSON(&membership); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    membership.UserID = userID
    
    if err := h.service.UpdateMembership(c.Request.Context(), &membership); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}

func (h *MembershipHandlers) UpgradeToPremium(c *gin.Context) {
    userID, _ := strconv.Atoi(c.Param("user_id"))
    if err := h.service.UpgradeToPremium(c.Request.Context(), userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}

func (h *MembershipHandlers) DowngradeToFree(c *gin.Context) {
    userID, _ := strconv.Atoi(c.Param("user_id"))
    if err := h.service.DowngradeToFree(c.Request.Context(), userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}