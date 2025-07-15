//api/src/sensors/application/sensor_service.go
package application

import (
    "context"

    "github.com/vicpoo/apigestion-solar-go/src/sensors/domain"
)

type SensorService struct {
    repo domain.SensorRepository
}

func NewSensorService(repo domain.SensorRepository) *SensorService {
    return &SensorService{repo: repo}
}

func (s *SensorService) CreateSensor(ctx context.Context, sensor *domain.Sensor) error {
    return s.repo.Create(ctx, sensor)
}

func (s *SensorService) GetSensor(ctx context.Context, id int) (*domain.Sensor, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *SensorService) GetUserSensors(ctx context.Context, userID int) ([]domain.Sensor, error) {
    return s.repo.GetByUserID(ctx, userID)
}

func (s *SensorService) UpdateSensor(ctx context.Context, sensor *domain.Sensor) error {
    return s.repo.Update(ctx, sensor)
}

func (s *SensorService) DeleteSensor(ctx context.Context, id int) error {
    return s.repo.Delete(ctx, id)
}
