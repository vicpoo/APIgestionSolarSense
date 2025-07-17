// src/sensors/application/post_sensor_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/domain"
)

type PostSensorUseCase struct {
    repo domain.SensorRepository
}

func NewPostSensorUseCase(repo domain.SensorRepository) *PostSensorUseCase {
    return &PostSensorUseCase{repo: repo}
}

func (uc *PostSensorUseCase) CreateSensor(ctx context.Context, sensor *domain.Sensor) error {
    return uc.repo.Create(ctx, sensor)
}