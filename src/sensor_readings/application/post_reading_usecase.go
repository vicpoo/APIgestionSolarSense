// src/sensor_readings/application/post_reading_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/domain"
)

type PostReadingUseCase struct {
    repo domain.ReadingRepository
}

func NewPostReadingUseCase(repo domain.ReadingRepository) *PostReadingUseCase {
    return &PostReadingUseCase{repo: repo}
}

func (uc *PostReadingUseCase) CreateReading(ctx context.Context, reading *domain.SensorReading) error {
    return uc.repo.Create(ctx, reading)
}