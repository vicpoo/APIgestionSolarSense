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
        // Endpoints públicos
        authGroup.POST("/email/register", loginController.RegisterEmail)
        authGroup.POST("/email/login", loginController.LoginEmail)
        authGroup.POST("/google", loginController.GoogleAuth)
        authGroup.GET("/public/users/:id", loginController.GetPublicUserInfo)
        authGroup.PUT("/user/actualizar", loginController.UpdateUserProfile)
            authGroup.GET("/google/users/:uid", loginController.GetGoogleUserByUID)
        authGroup.GET("/google/users", loginController.GetAllGoogleUsers)
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
            admin.GET("/users", loginController.GetAllUsers)
        }
    }
}