// src/alerts/application/deletealert_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/domain"
)

type DeleteAlertUseCase struct {
    repo domain.AlertRepository
}

func NewDeleteAlertUseCase(repo domain.AlertRepository) *DeleteAlertUseCase {
    return &DeleteAlertUseCase{repo: repo}
}

func (uc *DeleteAlertUseCase) DeleteAlert(ctx context.Context, id int) error {
    return uc.repo.Delete(ctx, id)
}