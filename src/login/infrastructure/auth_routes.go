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
	
	// Crear el handler para obtener usuarios
	getAuthHandler := NewGetAuthHandler(getUseCase)
	
	// Crear el controlador unificado con todos los parámetros necesarios
	loginController := NewLoginController(
		authHandlers, 
		getUseCase, 
		updateUseCase, 
		deleteUseCase,
		getAuthHandler, // Nuevo parámetro añadido
	)
	
	// Configurar rutas
	authGroup := router.Group("/api/auth")
	{
		// Autenticación por email
		authGroup.POST("/email/register", loginController.RegisterEmail)
		authGroup.POST("/email/login", loginController.LoginEmail)
		authGroup.GET("/email", loginController.GetUserByEmail)
		authGroup.PUT("/email", loginController.UpdateUserEmail)
		authGroup.DELETE("/email", loginController.DeleteUserByEmail)
		
		// Autenticación por Google
		authGroup.POST("/google", loginController.GoogleAuth)
		
		// Nuevos endpoints para obtener usuarios
		authGroup.GET("/users", loginController.GetAllUsers)
		authGroup.GET("/users/:user_id", loginController.GetUserByID)
	}
}