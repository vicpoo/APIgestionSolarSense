// src/sensor_thresholds/application/get_threshold_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/domain"
)

type GetThresholdUseCase struct {
    repo domain.ThresholdRepository
}

func NewGetThresholdUseCase(repo domain.ThresholdRepository) *GetThresholdUseCase {
    return &GetThresholdUseCase{repo: repo}
}

func (uc *GetThresholdUseCase) GetThresholds(ctx context.Context, sensorID int) (*domain.SensorThreshold, error) {
    return uc.repo.GetBySensorID(ctx, sensorID)
}