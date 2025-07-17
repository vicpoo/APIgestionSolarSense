// src/sensor_readings/application/put_reading_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/domain"
)

type PutReadingUseCase struct {
    repo domain.ReadingRepository
}

func NewPutReadingUseCase(repo domain.ReadingRepository) *PutReadingUseCase {
    return &PutReadingUseCase{repo: repo}
}

func (uc *PutReadingUseCase) UpdateReading(ctx context.Context, reading *domain.SensorReading) error {
    return uc.repo.Update(ctx, reading)
}