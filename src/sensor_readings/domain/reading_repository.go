// api/src/sensor_readings/domain/reading_repository.go
package domain

import (
    "context"
)

type ReadingRepository interface {
    Create(ctx context.Context, reading *SensorReading) error
    GetBySensorID(ctx context.Context, sensorID int, limit int) ([]SensorReading, error)
    GetLatestBySensorID(ctx context.Context, sensorID int) (*SensorReading, error)
    Update(ctx context.Context, reading *SensorReading) error
    Delete(ctx context.Context, id int) error
}