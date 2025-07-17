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
	
	// Crear todos los casos de uso y servicios
	authService := application.NewAuthService(repo)
	authHandlers := application.NewAuthHandlers(authService)
	getUseCase := application.NewGetAuthUseCase(repo)
	updateUseCase := application.NewUpdateAuthUseCase(repo)
	deleteUseCase := application.NewDeleteAuthUseCase(repo)
	
	// Crear el controlador unificado
	loginController := NewLoginController(authHandlers, getUseCase, updateUseCase, deleteUseCase)
	
	// Configurar rutas
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/email/register", loginController.RegisterEmail)
		authGroup.POST("/email/login", loginController.LoginEmail)
		authGroup.GET("/email", loginController.GetUserByEmail)
		authGroup.PUT("/email", loginController.UpdateUserEmail)
		authGroup.DELETE("/email", loginController.DeleteUserByEmail)
		authGroup.POST("/google", loginController.GoogleAuth)
	}
}