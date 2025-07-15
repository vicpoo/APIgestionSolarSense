//api/alerts/domain/alert.go

package domain

import "time"

type Alert struct {
    ID          int       `json:"id"`
    SensorID    int       `json:"sensor_id"`
    Message     string    `json:"message"`
    TriggeredAt time.Time `json:"triggered_at"`
    AlertType   string    `json:"alert_type"` // "lluvia", "umbral", "sistema"
    IsSent      bool      `json:"is_sent"`
}