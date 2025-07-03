package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/login/application"

	"github.com/vicpoo/apigestion-solar-go/src/core"
)

type AuthRouter struct {
	handlers *application.AuthHandlers
}

func NewAuthRouter(handlers *application.AuthHandlers) *AuthRouter {
	return &AuthRouter{handlers: handlers}
}

func (r *AuthRouter) SetupRoutes(router *gin.Engine) {
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/email/register", r.handlers.RegisterEmail)
		authGroup.POST("/email/login", r.handlers.LoginEmail)
		authGroup.POST("/google", r.handlers.GoogleAuth)
	}
}

func InitAuthRoutes(router *gin.Engine) {
	db := core.GetBD()
	repo := NewAuthRepository(db)
	service := application.NewAuthService(repo)
	handlers := application.NewAuthHandlers(service)
	
	authRouter := NewAuthRouter(handlers)
	authRouter.SetupRoutes(router)
}