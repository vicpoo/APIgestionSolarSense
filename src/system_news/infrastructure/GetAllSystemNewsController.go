// src/system_news/infrastructure/GetAllSystemNewsController.go
package infrastructure

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/application"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/domain/entities"
)

type GetAllSystemNewsController struct {
	getAllUseCase *application.GetAllSystemNewsUseCase
}

func NewGetAllSystemNewsController(getAllUseCase *application.GetAllSystemNewsUseCase) *GetAllSystemNewsController {
	return &GetAllSystemNewsController{
		getAllUseCase: getAllUseCase,
	}
}

func (ctrl *GetAllSystemNewsController) Run(c *gin.Context) {
    authType, _ := c.Get("authType")
    
    var news []entities.SystemNews
    var err error
    
    if authType == "email" {
        // Usuarios de email pueden ver todas las noticias
        news, err = ctrl.getAllUseCase.Run()
    } else {
        // Usuarios de Google solo ven noticias recientes (últimas 7 días)
        recentNews, err := ctrl.getAllUseCase.Run()
        if err == nil {
            for _, n := range recentNews {
                if time.Since(n.CreatedAt) <= 7*24*time.Hour {
                    news = append(news, n)
                }
            }
        }
    }
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "No se pudieron obtener las noticias del sistema",
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, news)
}


