//api\src\sensor_thresholds\domain\threshold_repository.go

package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/domain"

)

type ThresholdService struct {
    repo domain.ThresholdRepository
}

func NewThresholdService(repo domain.ThresholdRepository) *ThresholdService {
    return &ThresholdService{repo: repo}
}

func (s *ThresholdService) GetThresholds(ctx context.Context, sensorID int) (*domain.SensorThreshold, error) {
    return s.repo.GetBySensorID(ctx, sensorID)
}

func (s *ThresholdService) SetThresholds(ctx context.Context, threshold *domain.SensorThreshold) error {
    return s.repo.Upsert(ctx, threshold)
}