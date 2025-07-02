// UpdateSystemNewsUseCase.go
package application

import (
	repositories "github.com/vicpoo/apigestion-solar-go/src/system_news/domain"
	"github.com/vicpoo/apigestion-solar-go/src/system_news/domain/entities"
)

type UpdateSystemNewsUseCase struct {
	repo repositories.ISystemNews
}

func NewUpdateSystemNewsUseCase(repo repositories.ISystemNews) *UpdateSystemNewsUseCase {
	return &UpdateSystemNewsUseCase{repo: repo}
}

func (uc *UpdateSystemNewsUseCase) Run(news *entities.SystemNews) (*entities.SystemNews, error) {
	err := uc.repo.Update(news)
	if err != nil {
		return nil, err
	}
	return news, nil
}
