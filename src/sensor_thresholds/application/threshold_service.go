// api/src/sensor_thresholds/application/threshold_service.go
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

func (s *ThresholdService) CreateThreshold(ctx context.Context, threshold *domain.SensorThreshold) error {
    return s.repo.Create(ctx, threshold)
}

func (s *ThresholdService) UpdateThreshold(ctx context.Context, threshold *domain.SensorThreshold) error {
    return s.repo.Update(ctx, threshold)
}

func (s *ThresholdService) DeleteThreshold(ctx context.Context, sensorID int) error {
    return s.repo.Delete(ctx, sensorID)
}