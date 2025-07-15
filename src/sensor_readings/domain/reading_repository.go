//api/src/sensors/application/sensor_service.go
package domain

import (
    "context"

)

type ReadingRepository interface {
    Create(ctx context.Context, reading *SensorReading) error
    GetBySensorID(ctx context.Context, sensorID int, limit int) ([]SensorReading, error)
    GetLatestBySensorID(ctx context.Context, sensorID int) (*SensorReading, error)
}