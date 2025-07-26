//api/src/reports/infrastructure/mysql_report_repository.go

package infrastructure

import (
    "context"
    "database/sql"
    "github.com/vicpoo/apigestion-solar-go/src/core"
    "github.com/vicpoo/apigestion-solar-go/src/reports/domain"
   
)

type MySQLReportRepository struct {
    db *sql.DB
}

func NewMySQLReportRepository() domain.ReportRepository {
    return &MySQLReportRepository{db: core.GetBD()}
}

func (r *MySQLReportRepository) Create(ctx context.Context, report *domain.Report) error {
    query := `INSERT INTO reports 
        (user_id, sensor_id, file_name, storage_path, generated_from, generated_to, format) 
        VALUES (?, ?, ?, ?, ?, ?, ?)`
    _, err := r.db.ExecContext(ctx, query,
        report.UserID,
        report.SensorID,
        report.FileName,
        report.StoragePath,
        report.GeneratedFrom,
        report.GeneratedTo,
        report.Format,
    )
    return err
}

func (r *MySQLReportRepository) GetByID(ctx context.Context, id int) (*domain.Report, error) {
    query := `SELECT id, user_id, sensor_id, file_name, storage_path, 
                     generated_from, generated_to, created_at, format 
              FROM reports WHERE id = ?`
    row := r.db.QueryRowContext(ctx, query, id)
    
    var report domain.Report
    err := row.Scan(
        &report.ID,
        &report.UserID,
        &report.SensorID,
        &report.FileName,
        &report.StoragePath,
        &report.GeneratedFrom,
        &report.GeneratedTo,
        &report.CreatedAt,
        &report.Format,
    )
    if err != nil {
        return nil, err
    }
    return &report, nil
}

func (r *MySQLReportRepository) GetByUserID(ctx context.Context, userID int) ([]domain.Report, error) {
    query := `SELECT id, user_id, sensor_id, file_name, storage_path, 
                     generated_from, generated_to, created_at, format 
              FROM reports WHERE user_id = ?`
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var reports []domain.Report
    for rows.Next() {
        var report domain.Report
        err := rows.Scan(
            &report.ID,
            &report.UserID,
            &report.SensorID,
            &report.FileName,
            &report.StoragePath,
            &report.GeneratedFrom,
            &report.GeneratedTo,
            &report.CreatedAt,
            &report.Format,
        )
        if err != nil {
            return nil, err
        }
        reports = append(reports, report)
    }
    return reports, nil
}

func (r *MySQLReportRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM reports WHERE id = ?`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}

func (r *MySQLReportRepository) Update(ctx context.Context, report *domain.Report) error {
    query := `UPDATE reports SET 
        user_id = ?,
        sensor_id = ?,
        file_name = ?,
        storage_path = ?,
        generated_from = ?,
        generated_to = ?,
        format = ?
        WHERE id = ?`
    
    _, err := r.db.ExecContext(ctx, query,
        report.UserID,
        report.SensorID,
        report.FileName,
        report.StoragePath,
        report.GeneratedFrom,
        report.GeneratedTo,
        report.Format,
        report.ID,
    )
    return err
}


func (r *MySQLReportRepository) GetSensorReadingsByDate(ctx context.Context, date string) ([]domain.SensorReading, error) {
    query := `
        SELECT id, sensor_id, temperature, humidity, pressure, voltage, current, recorded_at, potencia
        FROM sensor_readings
        WHERE DATE(recorded_at) = ?`
    rows, err := r.db.QueryContext(ctx, query, date)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var readings []domain.SensorReading
    for rows.Next() {
        var reading domain.SensorReading
        err := rows.Scan(
            &reading.ID,
            &reading.SensorID,
            &reading.Temperature,
            &reading.Humidity,
            &reading.Pressure,
            &reading.Voltage,
            &reading.Current,
            &reading.RecordedAt,
            &reading.Potencia,
        )
        if err != nil {
            return nil, err
        }
        readings = append(readings, reading)
    }
    return readings, nil
}