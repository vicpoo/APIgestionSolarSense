// src/sensor_readings/application/get_reading_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/domain"
)

type GetReadingUseCase struct {
    repo domain.ReadingRepository
}

func NewGetReadingUseCase(repo domain.ReadingRepository) *GetReadingUseCase {
    return &GetReadingUseCase{repo: repo}
}

func (uc *GetReadingUseCase) GetReadings(ctx context.Context, sensorID int, limit int) ([]domain.SensorReading, error) {
    return uc.repo.GetBySensorID(ctx, sensorID, limit)
}

func (uc *GetReadingUseCase) GetLatestReading(ctx context.Context, sensorID int) (*domain.SensorReading, error) {
    return uc.repo.GetLatestBySensorID(ctx, sensorID)
}