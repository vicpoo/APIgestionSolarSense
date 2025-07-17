// src/sensor_readings/application/reading_service.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/domain"
)

type ReadingService struct {
    repo domain.ReadingRepository
}

func NewReadingService(repo domain.ReadingRepository) *ReadingService {
    return &ReadingService{repo: repo}
}

func (s *ReadingService) CreateReading(ctx context.Context, reading *domain.SensorReading) error {
    return s.repo.Create(ctx, reading)
}

func (s *ReadingService) GetReadings(ctx context.Context, sensorID int, limit int) ([]domain.SensorReading, error) {
    return s.repo.GetBySensorID(ctx, sensorID, limit)
}

func (s *ReadingService) GetLatestReading(ctx context.Context, sensorID int) (*domain.SensorReading, error) {
    return s.repo.GetLatestBySensorID(ctx, sensorID)
}

func (s *ReadingService) UpdateReading(ctx context.Context, reading *domain.SensorReading) error {
    return s.repo.Update(ctx, reading)
}

func (s *ReadingService) DeleteReading(ctx context.Context, id int) error {
    return s.repo.Delete(ctx, id)
}