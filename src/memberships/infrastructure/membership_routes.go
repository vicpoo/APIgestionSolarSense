//src/memberships/infrastructure/api/membership_routes.go

package infrastructure

import (
    "github.com/gin-gonic/gin"
 
    "github.com/vicpoo/apigestion-solar-go/src/memberships/application"
   
)

func InitMembershipRoutes(router *gin.Engine) {
    repo := NewMySQLMembershipRepository()
    service := application.NewMembershipService(repo)
    handlers := NewMembershipHandlers(service)

    membershipGroup := router.Group("/api/memberships")
    {
        membershipGroup.GET("/user/:user_id", handlers.GetUserMembership)
        membershipGroup.PUT("/user/:user_id", handlers.UpdateMembership)
        membershipGroup.POST("/user/:user_id/upgrade", handlers.UpgradeToPremium)
        membershipGroup.POST("/user/:user_id/downgrade", handlers.DowngradeToFree)
    }
}