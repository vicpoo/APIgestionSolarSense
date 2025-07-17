//api/src/sensor_readings/domain/reading.go

package domain

import "time"

type SensorReading struct {
    ID         int       `json:"id"`
    SensorID   int       `json:"sensor_id"`
    Temperature *float64  `json:"temperature,omitempty"`
    Humidity    *float64  `json:"humidity,omitempty"`
    Pressure    *float64  `json:"pressure,omitempty"`
    Voltage     *float64  `json:"voltage,omitempty"`
    Current     *float64  `json:"current,omitempty"`
    RecordedAt time.Time `json:"recorded_at"`
}