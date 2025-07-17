//src/alerts/domain/alert_repository.go

package domain

import (
    "context"
)

type AlertRepository interface {
    Create(ctx context.Context, alert *Alert) error
    GetByID(ctx context.Context, id int) (*Alert, error)
    GetBySensorID(ctx context.Context, sensorID int, limit int) ([]Alert, error)
    GetUnsent(ctx context.Context) ([]Alert, error)
    MarkAsSent(ctx context.Context, id int) error
    Update(ctx context.Context, alert *Alert) error
    Delete(ctx context.Context, id int) error
}