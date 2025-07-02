// GetAllSystemNewsController.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/application"
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
	news, err := ctrl.getAllUseCase.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "No se pudieron obtener las noticias del sistema",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, news)
}


