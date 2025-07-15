//api\src\sensor_thresholds\domain\threshold_repository.go

package infrastructure

import (
    "context"
    "database/sql"
    "github.com/vicpoo/apigestion-solar-go/src/core"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/domain"

)

type MySQLThresholdRepository struct {
    db *sql.DB
}

func NewMySQLThresholdRepository() domain.ThresholdRepository {
    return &MySQLThresholdRepository{db: core.GetBD()}
}

func (r *MySQLThresholdRepository) GetBySensorID(ctx context.Context, sensorID int) (*domain.SensorThreshold, error) {
    query := `SELECT * FROM sensor_thresholds WHERE sensor_id = ?`
    row := r.db.QueryRowContext(ctx, query, sensorID)
    
    var threshold domain.SensorThreshold
    err := row.Scan(
        &threshold.ID,
        &threshold.SensorID,
        &threshold.MinTemperature,
        &threshold.MaxTemperature,
        &threshold.MinHumidity,
        &threshold.MaxHumidity,
        &threshold.MinPressure,
        &threshold.MaxPressure,
        &threshold.MinVoltage,
        &threshold.MaxVoltage,
        &threshold.MinCurrent,
        &threshold.MaxCurrent,
        &threshold.ConfiguredBy,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &threshold, nil
}

func (r *MySQLThresholdRepository) Upsert(ctx context.Context, threshold *domain.SensorThreshold) error {
    query := `INSERT INTO sensor_thresholds (
        sensor_id, min_temperature, max_temperature, min_humidity, max_humidity,
        min_pressure, max_pressure, min_voltage, max_voltage, min_current, max_current, configured_by
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    ON DUPLICATE KEY UPDATE
        min_temperature = VALUES(min_temperature),
        max_temperature = VALUES(max_temperature),
        min_humidity = VALUES(min_humidity),
        max_humidity = VALUES(max_humidity),
        min_pressure = VALUES(min_pressure),
        max_pressure = VALUES(max_pressure),
        min_voltage = VALUES(min_voltage),
        max_voltage = VALUES(max_voltage),
        min_current = VALUES(min_current),
        max_current = VALUES(max_current),
        configured_by = VALUES(configured_by)`
    
    _, err := r.db.ExecContext(ctx, query,
        threshold.SensorID,
        threshold.MinTemperature,
        threshold.MaxTemperature,
        threshold.MinHumidity,
        threshold.MaxHumidity,
        threshold.MinPressure,
        threshold.MaxPressure,
        threshold.MinVoltage,
        threshold.MaxVoltage,
        threshold.MinCurrent,
        threshold.MaxCurrent,
        threshold.ConfiguredBy,
    )
    return err
}