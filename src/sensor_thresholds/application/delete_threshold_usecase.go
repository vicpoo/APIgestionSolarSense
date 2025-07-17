// src/sensor_thresholds/application/delete_threshold_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/domain"
)

type DeleteThresholdUseCase struct {
    repo domain.ThresholdRepository
}

func NewDeleteThresholdUseCase(repo domain.ThresholdRepository) *DeleteThresholdUseCase {
    return &DeleteThresholdUseCase{repo: repo}
}

func (uc *DeleteThresholdUseCase) DeleteThreshold(ctx context.Context, sensorID int) error {
    return uc.repo.Delete(ctx, sensorID)
}