// src/memberships/infrastructure/membership_routes.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/application"
)

// src/memberships/infrastructure/membership_routes.go
func InitMembershipRoutes(router *gin.Engine) {
    repo := NewMySQLMembershipRepository()

    // Crear casos de uso
    getUseCase := application.NewGetMembershipUseCase(repo)
    postUseCase := application.NewPostMembershipUseCase(repo)
    putUseCase := application.NewPutMembershipUseCase(repo)
    deleteUseCase := application.NewDeleteMembershipUseCase(repo)
    registerUseCase := application.NewRegisterUseCase(repo) // Nuevo

    // Crear handlers
    getHandler := NewGetMembershipHandler(getUseCase)
    postHandler := NewPostMembershipHandler(postUseCase)
    putHandler := NewPutMembershipHandler(putUseCase)
    deleteHandler := NewDeleteMembershipHandler(deleteUseCase)
    registerHandler := NewRegisterHandler(registerUseCase) // Nuevo

    // Crear controlador
    controller := NewMembershipController(
        getHandler, 
        postHandler, 
        putHandler, 
        deleteHandler,
        registerHandler, // Nuevo
    )

    // Configurar rutas SIN seguridad
    membershipGroup := router.Group("/api/memberships")
    {
        membershipGroup.POST("/register", controller.RegisterUser) // Nueva ruta
        membershipGroup.GET("", controller.GetAllUsers)
        membershipGroup.GET("/user/:user_id", controller.GetUserMembership)
        membershipGroup.PUT("/user/:user_id", controller.CreateOrUpdate)
        membershipGroup.DELETE("/user/:user_id", controller.DeleteMembership)
        membershipGroup.POST("/user/:user_id/upgrade", controller.UpgradeToPremium)
        membershipGroup.POST("/user/:user_id/downgrade", controller.DowngradeToFree)
    }
}