// src/sensor_thresholds/application/put_threshold_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/domain"
)

type PutThresholdUseCase struct {
    repo domain.ThresholdRepository
}

func NewPutThresholdUseCase(repo domain.ThresholdRepository) *PutThresholdUseCase {
    return &PutThresholdUseCase{repo: repo}
}

func (uc *PutThresholdUseCase) UpdateThreshold(ctx context.Context, threshold *domain.SensorThreshold) error {
    return uc.repo.Update(ctx, threshold)
}