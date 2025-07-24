// src/system_news/infrastructure/GetSystemNewsByIdController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/application"
)

type GetSystemNewsByIdController struct {
	getByIdUseCase *application.GetSystemNewsByIdUseCase
}

func NewGetSystemNewsByIdController(getByIdUseCase *application.GetSystemNewsByIdUseCase) *GetSystemNewsByIdController {
	return &GetSystemNewsByIdController{
		getByIdUseCase: getByIdUseCase,
	}
}

func (ctrl *GetSystemNewsByIdController) Run(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID inválido",
			"error":   err.Error(),
		})
		return
	}

	news, err := ctrl.getByIdUseCase.Run(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "No se pudo obtener la noticia del sistema",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, news)
}
