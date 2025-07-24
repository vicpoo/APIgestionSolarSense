// src/system_news/infrastructure/system_news_routes.go
package infrastructure

import (
	"github.com/gin-gonic/gin"
)

type SystemNewsRouter struct {
	engine *gin.Engine
}

func NewSystemNewsRouter(engine *gin.Engine) *SystemNewsRouter {
	return &SystemNewsRouter{
		engine: engine,
	}
}

func (router *SystemNewsRouter) Run() {
	// Inicializar dependencias
	createController, getByIdController, updateController, deleteController, getAllController := InitSystemNewsDependencies()

	// Grupo de rutas para system_news
	newsGroup := router.engine.Group("/system-news")
	{
		newsGroup.POST("/", createController.Run)
		newsGroup.GET("/:id", getByIdController.Run)
		newsGroup.PUT("/:id", updateController.Run)
		newsGroup.DELETE("/:id", deleteController.Run)
		newsGroup.GET("/", getAllController.Run)
	}
}
