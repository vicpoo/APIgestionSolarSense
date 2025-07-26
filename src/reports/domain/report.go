//api/src/reports/report.go
package domain

import "time"

type Report struct {
    ID            int        `json:"id"`
    UserID        int        `json:"user_id"`
    SensorID      *int       `json:"sensor_id,omitempty"`
    FileName      string     `json:"file_name"`
    StoragePath   string     `json:"storage_path"`
    GeneratedFrom *time.Time `json:"generated_from,omitempty"`
    GeneratedTo   *time.Time `json:"generated_to,omitempty"`
    CreatedAt     time.Time  `json:"created_at"`
    Format        string     `json:"format"` // "PDF" o "Excel"
}

// Nuevo struct para las lecturas de sensores
type SensorReading struct {
    ID         int        `json:"id"`
    SensorID   int        `json:"sensor_id"`
    Temperature *float64   `json:"temperature,omitempty"`
    Humidity    *float64   `json:"humidity,omitempty"`
    Pressure    *float64   `json:"pressure,omitempty"`
    Voltage     *float64   `json:"voltage,omitempty"`
    Current     *float64   `json:"current,omitempty"`
    Potencia    *float64   `json:"potencia,omitempty"`
    RecordedAt  time.Time  `json:"recorded_at"`
}

type GenerateReportRequest struct {
    Date             string `json:"date" binding:"required"` // Formato: "YYYY-MM-DD"
    RequestedByEmail string `json:"requested_by_email" binding:"required,email"`
    Format           string `json:"format"` // "PDF" (default) o "Excel"
}