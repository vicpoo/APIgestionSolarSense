// CreateSystemNewsUseCase.go
package application

import (
	repositories "github.com/vicpoo/apigestion-solar-go/src/system_news/domain"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/domain/entities"
)

type CreateSystemNewsUseCase struct {
	repo repositories.ISystemNews
}

func NewCreateSystemNewsUseCase(repo repositories.ISystemNews) *CreateSystemNewsUseCase {
	return &CreateSystemNewsUseCase{repo: repo}
}

func (uc *CreateSystemNewsUseCase) Run(news *entities.SystemNews) (*entities.SystemNews, error) {
	err := uc.repo.Save(news)
	if err != nil {
		return nil, err
	}
	return news, nil
}
