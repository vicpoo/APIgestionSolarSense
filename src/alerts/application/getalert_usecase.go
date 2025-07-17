// src/alerts/application/getalert_usecase.go
package application

import (
	"context"

	"github.com/vicpoo/apigestion-solar-go/src/alerts/domain"
)

type GetAlertUseCase struct {
	repo domain.AlertRepository
}

func NewGetAlertUseCase(repo domain.AlertRepository) *GetAlertUseCase {
	return &GetAlertUseCase{repo: repo}
}

func (uc *GetAlertUseCase) GetAlert(ctx context.Context, id int) (*domain.Alert, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *GetAlertUseCase) GetSensorAlerts(ctx context.Context, sensorID int, limit int) ([]domain.Alert, error) {
	return uc.repo.GetBySensorID(ctx, sensorID, limit)
}

func (uc *GetAlertUseCase) GetUnsentAlerts(ctx context.Context) ([]domain.Alert, error) {
	return uc.repo.GetUnsent(ctx)
}
