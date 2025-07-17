// src/sensors/application/delete_sensor_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/domain"
)

type DeleteSensorUseCase struct {
    repo domain.SensorRepository
}

func NewDeleteSensorUseCase(repo domain.SensorRepository) *DeleteSensorUseCase {
    return &DeleteSensorUseCase{repo: repo}
}

func (uc *DeleteSensorUseCase) DeleteSensor(ctx context.Context, id int) error {
    return uc.repo.Delete(ctx, id)
}