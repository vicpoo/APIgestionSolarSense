// UpdateSystemNewsController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/application"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/domain/entities"
)

type UpdateSystemNewsController struct {
	updateUseCase *application.UpdateSystemNewsUseCase
}

func NewUpdateSystemNewsController(updateUseCase *application.UpdateSystemNewsUseCase) *UpdateSystemNewsController {
	return &UpdateSystemNewsController{
		updateUseCase: updateUseCase,
	}
}

func (ctrl *UpdateSystemNewsController) Run(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID inválido",
			"error":   err.Error(),
		})
		return
	}

	var newsRequest struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		AuthorID int32  `json:"author_id"`
	}

	if err := c.ShouldBindJSON(&newsRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Datos inválidos",
			"error":   err.Error(),
		})
		return
	}

	news := &entities.SystemNews{
		ID:       int32(id),
		Title:    newsRequest.Title,
		Content:  newsRequest.Content,
		AuthorID: newsRequest.AuthorID,
	}

	updatedNews, err := ctrl.updateUseCase.Run(news)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "No se pudo actualizar la noticia del sistema",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedNews)
}
