//api/src/sensor_readings/infrastructure/mysql_reading_repository.go

package infrastructure

import (
	"context"
	"database/sql"

	"github.com/vicpoo/apigestion-solar-go/src/core"
	"github.com/vicpoo/apigestion-solar-go/src/sensor_readings/domain"
)

type MySQLReadingRepository struct {
	db *sql.DB
}

func NewMySQLReadingRepository() domain.ReadingRepository {
	return &MySQLReadingRepository{db: core.GetBD()}
}

func (r *MySQLReadingRepository) Create(ctx context.Context, reading *domain.SensorReading) error {
	query := `INSERT INTO sensor_readings 
        (sensor_id, temperature, humidity, pressure, voltage, current) 
        VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		reading.SensorID,
		reading.Temperature,
		reading.Humidity,
		reading.Pressure,
		reading.Voltage,
		reading.Current,
	)
	return err
}

func (r *MySQLReadingRepository) GetBySensorID(ctx context.Context, sensorID int, limit int) ([]domain.SensorReading, error) {
	query := `SELECT id, sensor_id, temperature, humidity, pressure, voltage, current, recorded_at 
              FROM sensor_readings 
              WHERE sensor_id = ? 
              ORDER BY recorded_at DESC 
              LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, sensorID, limit)
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
		)
		if err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}
	return readings, nil
}

func (r *MySQLReadingRepository) GetLatestBySensorID(ctx context.Context, sensorID int) (*domain.SensorReading, error) {
	query := `SELECT id, sensor_id, temperature, humidity, pressure, voltage, current, recorded_at 
              FROM sensor_readings 
              WHERE sensor_id = ? 
              ORDER BY recorded_at DESC 
              LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, sensorID)

	var reading domain.SensorReading
	err := row.Scan(
		&reading.ID,
		&reading.SensorID,
		&reading.Temperature,
		&reading.Humidity,
		&reading.Pressure,
		&reading.Voltage,
		&reading.Current,
		&reading.RecordedAt,
	)
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

func (r *MySQLReadingRepository) Update(ctx context.Context, reading *domain.SensorReading) error {
    query := `UPDATE sensor_readings SET 
        sensor_id = ?,
        temperature = ?,
        humidity = ?,
        pressure = ?,
        voltage = ?,
        current = ?
        WHERE id = ?`
    
    _, err := r.db.ExecContext(ctx, query,
        reading.SensorID,
        reading.Temperature,
        reading.Humidity,
        reading.Pressure,
        reading.Voltage,
        reading.Current,
        reading.ID,
    )
    return err
}

func (r *MySQLReadingRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM sensor_readings WHERE id = ?`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}
