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
    
    store := cookie.NewStore([]byte("secret"))
    router.Use(sessions.Sessions("mysession", store))
    
    authGroup := router.Group("/api/auth")
    {
        // Endpoints públicos
        authGroup.POST("/email/register", loginController.RegisterEmail)
        authGroup.POST("/email/login", loginController.LoginEmail)
        authGroup.POST("/google", loginController.GoogleAuth)
        
        // Nuevo endpoint público para obtener usuarios
        authGroup.GET("/users", loginController.GetAllUsers) // <- Sin middlewares
        
        // Endpoints protegidos para usuarios normales
        private := authGroup.Group("")
        private.Use(core.AuthMiddleware())
        {
            private.GET("/me", loginController.GetCurrentUser)
            private.PUT("/email", loginController.UpdateUserEmail)
            private.PUT("/password", loginController.UpdatePassword)
            private.DELETE("/account", loginController.DeleteAccount)
        }
        
        // Endpoints de admin (protegidos)
        admin := authGroup.Group("")
        admin.Use(core.AuthMiddleware())
        admin.Use(core.AdminMiddleware())
        {
            admin.GET("/users/:id", loginController.GetUserByID)
        }
    }
}