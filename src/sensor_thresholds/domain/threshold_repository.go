// api/src/sensor_thresholds/domain/threshold_repository.go
package domain

import (
    "context"
)

type ThresholdRepository interface {
    GetBySensorID(ctx context.Context, sensorID int) (*SensorThreshold, error)
    Create(ctx context.Context, threshold *SensorThreshold) error
    Update(ctx context.Context, threshold *SensorThreshold) error
    Delete(ctx context.Context, sensorID int) error
}