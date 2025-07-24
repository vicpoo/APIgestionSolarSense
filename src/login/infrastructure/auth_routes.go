// api/src/login/infrastructure/auth_routes.go
package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/core"
	"github.com/vicpoo/apigestion-solar-go/src/login/application"
)

func InitAuthRoutes(router *gin.Engine) {
	db := core.GetBD()
	repo := NewAuthRepository(db)
	
	authService := application.NewAuthService(repo)
	authHandlers := application.NewAuthHandlers(authService)
	getUseCase := application.NewGetAuthUseCase(repo)
	updateUseCase := application.NewUpdateAuthUseCase(repo)
	deleteUseCase := application.NewDeleteAuthUseCase(repo)
	getAuthHandler := NewGetAuthHandler(getUseCase)
	
	loginController := NewLoginController(
		authHandlers, 
		getUseCase, 
		updateUseCase, 
		deleteUseCase,
		getAuthHandler,
	)
	
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/email/register", loginController.RegisterEmail)
		authGroup.POST("/email/login", loginController.LoginEmail)
		authGroup.POST("/google", loginController.GoogleAuth)
		
		protected := authGroup.Group("")
		protected.Use(AuthMiddleware())
		{
			emailGroup := protected.Group("")
			emailGroup.Use(EmailUserMiddleware())
			{
				emailGroup.PUT("/email", loginController.UpdateUserEmail)
				emailGroup.PUT("/password", loginController.UpdatePassword)
			}
			
			protected.GET("/me", loginController.GetCurrentUser)
			protected.DELETE("/account", loginController.DeleteAccount)
			protected.DELETE("/account/:id", loginController.DeleteAccount)
		}
		
		adminGroup := protected.Group("")
		adminGroup.Use(AdminMiddleware())
		{
			adminGroup.GET("/users", loginController.GetAllUsers)
			adminGroup.GET("/users/:id", loginController.GetUserByID)
		}
	}
}