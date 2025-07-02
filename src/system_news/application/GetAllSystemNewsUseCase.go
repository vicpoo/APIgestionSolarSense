// GetAllSystemNewsUseCase.go
package application

import (
	repositories "github.com/vicpoo/apigestion-solar-go/src/system_news/domain"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/domain/entities"
)

type GetAllSystemNewsUseCase struct {
	repo repositories.ISystemNews
}

func NewGetAllSystemNewsUseCase(repo repositories.ISystemNews) *GetAllSystemNewsUseCase {
	return &GetAllSystemNewsUseCase{repo: repo}
}

func (uc *GetAllSystemNewsUseCase) Run() ([]entities.SystemNews, error) {
	news, err := uc.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return news, nil
}
