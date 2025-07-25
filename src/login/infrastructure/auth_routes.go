// api/src/login/infrastructure/auth_routes.go
package infrastructure

import (
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
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
    
    // Configurar almacenamiento de sesión (opcional si usas JWT)
    store := cookie.NewStore([]byte("secret"))
    router.Use(sessions.Sessions("mysession", store))
    
    authGroup := router.Group("/api/auth")
    {
        authGroup.POST("/email/register", loginController.RegisterEmail)
        authGroup.POST("/email/login", loginController.LoginEmail)
        authGroup.POST("/google", loginController.GoogleAuth)
        
        // Proteger estos endpoints con JWT
        private := authGroup.Group("")
        private.Use(core.AuthMiddleware())
        private.Use(core.EmailUserMiddleware()) // Solo para usuarios de email
        {
            private.GET("/me", loginController.GetCurrentUser)
            private.PUT("/email", loginController.UpdateUserEmail)
            private.PUT("/password", loginController.UpdatePassword)
            private.DELETE("/account", loginController.DeleteAccount)
        }
        
        // Endpoints de admin
        authGroup.GET("/users", loginController.GetAllUsers)
        authGroup.GET("/users/:id", loginController.GetUserByID)
    }
}