// DeleteSystemNewsController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/application"
)

type DeleteSystemNewsController struct {
	deleteUseCase *application.DeleteSystemNewsUseCase
}

func NewDeleteSystemNewsController(deleteUseCase *application.DeleteSystemNewsUseCase) *DeleteSystemNewsController {
	return &DeleteSystemNewsController{
		deleteUseCase: deleteUseCase,
	}
}

func (ctrl *DeleteSystemNewsController) Run(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID inválido",
			"error":   err.Error(),
		})
		return
	}

	errDelete := ctrl.deleteUseCase.Run(int32(id))
	if errDelete != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "No se pudo eliminar la noticia",
			"error":   errDelete.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Noticia eliminada exitosamente",
	})
}
