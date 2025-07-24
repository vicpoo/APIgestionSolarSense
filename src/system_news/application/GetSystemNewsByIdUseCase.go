//src/system_news/application GetSystemNewsByIdUseCase.go
package application

import (
	repositories "github.com/vicpoo/apigestion-solar-go/src/system_news/domain"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/domain/entities"
)

type GetSystemNewsByIdUseCase struct {
	repo repositories.ISystemNews
}

func NewGetSystemNewsByIdUseCase(repo repositories.ISystemNews) *GetSystemNewsByIdUseCase {
	return &GetSystemNewsByIdUseCase{repo: repo}
}

func (uc *GetSystemNewsByIdUseCase) Run(id int32) (*entities.SystemNews, error) {
	news, err := uc.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return news, nil
}
