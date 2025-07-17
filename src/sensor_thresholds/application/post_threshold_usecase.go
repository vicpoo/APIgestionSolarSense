// src/sensor_thresholds/application/post_threshold_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/domain"
)

type PostThresholdUseCase struct {
    repo domain.ThresholdRepository
}

func NewPostThresholdUseCase(repo domain.ThresholdRepository) *PostThresholdUseCase {
    return &PostThresholdUseCase{repo: repo}
}

func (uc *PostThresholdUseCase) CreateThreshold(ctx context.Context, threshold *domain.SensorThreshold) error {
    return uc.repo.Create(ctx, threshold)
}