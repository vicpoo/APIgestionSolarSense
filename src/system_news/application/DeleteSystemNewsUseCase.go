// DeleteSystemNewsUseCase.go
package application

import repositories "github.com/vicpoo/apigestion-solar-go/src/system_news/domain"

type DeleteSystemNewsUseCase struct {
	repo repositories.ISystemNews
}

func NewDeleteSystemNewsUseCase(repo repositories.ISystemNews) *DeleteSystemNewsUseCase {
	return &DeleteSystemNewsUseCase{repo: repo}
}

func (uc *DeleteSystemNewsUseCase) Run(id int32) error {
	err := uc.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
