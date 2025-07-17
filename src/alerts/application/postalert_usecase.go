// src/alerts/application/postalert_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/domain"
)

type PostAlertUseCase struct {
    repo domain.AlertRepository
}

func NewPostAlertUseCase(repo domain.AlertRepository) *PostAlertUseCase {
    return &PostAlertUseCase{repo: repo}
}

func (uc *PostAlertUseCase) CreateAlert(ctx context.Context, alert *domain.Alert) error {
    return uc.repo.Create(ctx, alert)
}

func (uc *PostAlertUseCase) MarkAlertAsSent(ctx context.Context, id int) error {
    return uc.repo.MarkAsSent(ctx, id)
}