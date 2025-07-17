// src/sensors/application/get_sensor_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/domain"
)

type GetSensorUseCase struct {
    repo domain.SensorRepository
}

func NewGetSensorUseCase(repo domain.SensorRepository) *GetSensorUseCase {
    return &GetSensorUseCase{repo: repo}
}

func (uc *GetSensorUseCase) GetSensor(ctx context.Context, id int) (*domain.Sensor, error) {
    return uc.repo.GetByID(ctx, id)
}

func (uc *GetSensorUseCase) GetUserSensors(ctx context.Context, userID int) ([]domain.Sensor, error) {
    return uc.repo.GetByUserID(ctx, userID)
}