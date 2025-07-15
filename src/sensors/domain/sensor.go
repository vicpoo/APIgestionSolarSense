//api/src/sensors/domain/sensor.go
package domain

import "time"

type Sensor struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Name      string    `json:"name"`
    Type      string    `json:"type"` // Enum: "INA219", "DS18B20", etc.
    Location  string    `json:"location"`
    CreatedAt time.Time `json:"created_at"`
}