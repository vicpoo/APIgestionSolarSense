// src/sensors/application/put_sensor_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/domain"
)

type PutSensorUseCase struct {
    repo domain.SensorRepository
}

func NewPutSensorUseCase(repo domain.SensorRepository) *PutSensorUseCase {
    return &PutSensorUseCase{repo: repo}
}

func (uc *PutSensorUseCase) UpdateSensor(ctx context.Context, sensor *domain.Sensor) error {
    return uc.repo.Update(ctx, sensor)
}