//api\src\sensor_thresholds\domain\threshold_repository.go

package domain

import (
    "context"
 
)

type SettingsRepository interface {
    GetByUserID(ctx context.Context, userID int) (*NotificationSettings, error)
    Update(ctx context.Context, settings *NotificationSettings) error
}