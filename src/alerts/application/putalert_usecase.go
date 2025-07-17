// src/alerts/application/putalert_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/domain"
)

type PutAlertUseCase struct {
    repo domain.AlertRepository
}

func NewPutAlertUseCase(repo domain.AlertRepository) *PutAlertUseCase {
    return &PutAlertUseCase{repo: repo}
}

func (uc *PutAlertUseCase) UpdateAlert(ctx context.Context, alert *domain.Alert) error {
    return uc.repo.Update(ctx, alert)
}