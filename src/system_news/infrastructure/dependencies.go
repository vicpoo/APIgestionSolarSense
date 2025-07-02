// dependencies.go
package infrastructure

import (
	"github.com/vicpoo/apigestion-solar-go/src/system_news/application"
)

func InitSystemNewsDependencies() (
	*CreateSystemNewsController,
	*GetSystemNewsByIdController,
	*UpdateSystemNewsController,
	*DeleteSystemNewsController,
	*GetAllSystemNewsController,
) {
	repo := NewMySQLSystemNewsRepository()

	createUseCase := application.NewCreateSystemNewsUseCase(repo)
	getByIdUseCase := application.NewGetSystemNewsByIdUseCase(repo)
	updateUseCase := application.NewUpdateSystemNewsUseCase(repo)
	deleteUseCase := application.NewDeleteSystemNewsUseCase(repo)
	getAllUseCase := application.NewGetAllSystemNewsUseCase(repo)

	createController := NewCreateSystemNewsController(createUseCase)
	getByIdController := NewGetSystemNewsByIdController(getByIdUseCase)
	updateController := NewUpdateSystemNewsController(updateUseCase)
	deleteController := NewDeleteSystemNewsController(deleteUseCase)
	getAllController := NewGetAllSystemNewsController(getAllUseCase)

	return createController, getByIdController, updateController, deleteController, getAllController
}
