// src/memberships/infrastructure/membership_routes.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/application"
)

func InitMembershipRoutes(router *gin.Engine) {
    repo := NewMySQLMembershipRepository()

    // Crear casos de uso
    getUseCase := application.NewGetMembershipUseCase(repo)
    postUseCase := application.NewPostMembershipUseCase(repo)
    putUseCase := application.NewPutMembershipUseCase(repo)
    deleteUseCase := application.NewDeleteMembershipUseCase(repo)

    // Crear handlers
    getHandler := NewGetMembershipHandler(getUseCase)
    postHandler := NewPostMembershipHandler(postUseCase)
    putHandler := NewPutMembershipHandler(putUseCase)
    deleteHandler := NewDeleteMembershipHandler(deleteUseCase)

    // Crear controlador
    controller := NewMembershipController(getHandler, postHandler, putHandler, deleteHandler)

    // Configurar rutas SIN seguridad
    membershipGroup := router.Group("/api/memberships")
    {
        membershipGroup.GET("", controller.GetAllUsers)
        membershipGroup.GET("/user/:user_id", controller.GetUserMembership)
        membershipGroup.PUT("/user/:user_id", controller.CreateOrUpdate)
        membershipGroup.DELETE("/user/:user_id", controller.DeleteMembership)
        membershipGroup.POST("/user/:user_id/upgrade", controller.UpgradeToPremium)
        membershipGroup.POST("/user/:user_id/downgrade", controller.DowngradeToFree)
    }
}