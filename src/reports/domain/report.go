//api/src/reports/report.go

package domain

import "time"

type Report struct {
    ID            int       `json:"id"`
    UserID        int       `json:"user_id"`
    SensorID      *int      `json:"sensor_id,omitempty"`
    FileName      string    `json:"file_name"`
    StoragePath   string    `json:"storage_path"`
    GeneratedFrom *time.Time `json:"generated_from,omitempty"`
    GeneratedTo   *time.Time `json:"generated_to,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
    Format        string    `json:"format"` // "PDF" o "Excel"
}