//api/src/sensors/domain/sensor_repository.go
package domain

import (
    "context"
)

type SensorRepository interface {
    Create(ctx context.Context, sensor *Sensor) error
    GetByID(ctx context.Context, id int) (*Sensor, error)
    GetByUserID(ctx context.Context, userID int) ([]Sensor, error)
    Update(ctx context.Context, sensor *Sensor) error
    Delete(ctx context.Context, id int) error
}