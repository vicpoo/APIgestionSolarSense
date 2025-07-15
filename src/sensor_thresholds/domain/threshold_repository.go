//api\src\sensor_thresholds\domain\threshold_repository.go

package domain

import (
    "context"

)

type ThresholdRepository interface {
    GetBySensorID(ctx context.Context, sensorID int) (*SensorThreshold, error)
    Upsert(ctx context.Context, threshold *SensorThreshold) error
}