//api/src/system_news/infrastructure CreateSystemNewsController.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/application"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/domain/entities"
)

type CreateSystemNewsController struct {
	createUseCase *application.CreateSystemNewsUseCase
}

func NewCreateSystemNewsController(useCase *application.CreateSystemNewsUseCase) *CreateSystemNewsController {
	return &CreateSystemNewsController{
		createUseCase: useCase,
	}
}

func (ctrl *CreateSystemNewsController) Run(c *gin.Context) {
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

	news := entities.NewSystemNews(
		newsRequest.Title,
		newsRequest.Content,
		newsRequest.AuthorID,
	)

	createdNews, err := ctrl.createUseCase.Run(news)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "No se pudo crear la noticia",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdNews)
}
