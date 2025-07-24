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
        // Rutas públicas
        authGroup.POST("/email/register", loginController.RegisterEmail)
        authGroup.POST("/email/login", loginController.LoginEmail)
        authGroup.POST("/google", loginController.GoogleAuth)
        
        // Rutas protegidas (solo admin)
        adminGroup := authGroup.Group("")
        adminGroup.Use(AdminMiddleware()) // Aplicar middleware a todas las rutas siguientes
        {
            adminGroup.GET("/email", loginController.GetUserByEmail)
            adminGroup.PUT("/email", loginController.UpdateUserEmail)
            adminGroup.DELETE("/email", loginController.DeleteUserByEmail)
            adminGroup.GET("/users", loginController.GetAllUsers)
            adminGroup.GET("/users/:user_id", loginController.GetUserByID)
        }
    }
}