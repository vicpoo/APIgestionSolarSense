// src/memberships/infrastructure/membership_routes.go
package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/memberships/application"
)

func InitMembershipRoutes(router *gin.Engine) {
	repo := NewMySQLMembershipRepository()

	getUseCase := application.NewGetMembershipUseCase(repo)
	postUseCase := application.NewPostMembershipUseCase(repo)
	putUseCase := application.NewPutMembershipUseCase(repo)
	deleteUseCase := application.NewDeleteMembershipUseCase(repo)
	registerUseCase := application.NewRegisterUseCase(repo)

	getHandler := NewGetMembershipHandler(getUseCase)
	postHandler := NewPostMembershipHandler(postUseCase)
	putHandler := NewPutMembershipHandler(putUseCase)
	deleteHandler := NewDeleteMembershipHandler(deleteUseCase)
	registerHandler := NewRegisterHandler(registerUseCase)

	controller := NewMembershipController(
		getHandler, 
		postHandler, 
		putHandler, 
		deleteHandler,
		registerHandler,
	)

	membershipGroup := router.Group("/api/memberships")
	{
		membershipGroup.POST("/fix-missing", controller.FixMissingMemberships) // Nueva ruta
		membershipGroup.GET("", controller.GetAllUsers)
		membershipGroup.GET("/user/:user_id", controller.GetUserMembership)
		membershipGroup.PUT("/user/:user_id", controller.CreateOrUpdate)
		membershipGroup.POST("/user/:user_id/upgrade", controller.UpgradeToPremium)
		membershipGroup.POST("/user/:user_id/downgrade", controller.DowngradeToFree)
		membershipGroup.POST("/register", controller.RegisterUser)
	}
}