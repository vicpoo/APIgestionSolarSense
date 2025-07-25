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
        
        // Endpoints que antes estaban protegidos
        authGroup.GET("/me", loginController.GetCurrentUser)
        authGroup.PUT("/email", loginController.UpdateUserEmail)
        authGroup.PUT("/password", loginController.UpdatePassword)
        authGroup.DELETE("/account", loginController.DeleteAccount)
        authGroup.DELETE("/account/:id", loginController.DeleteAccount)
        
        // Endpoints de admin
        authGroup.GET("/users", loginController.GetAllUsers)
        authGroup.GET("/users/:id", loginController.GetUserByID)
    }
}