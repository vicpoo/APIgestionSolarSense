//api/src/sensors/infrastructure/mysql_sensor_repository.go

package infrastructure

import (
    "context"
    "database/sql"
    "github.com/vicpoo/apigestion-solar-go/src/core"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/domain"
   
)

type MySQLSensorRepository struct {
    db *sql.DB
}

func NewMySQLSensorRepository() domain.SensorRepository {
    return &MySQLSensorRepository{db: core.GetBD()}
}

func (r *MySQLSensorRepository) Create(ctx context.Context, sensor *domain.Sensor) error {
    query := `INSERT INTO sensors (user_id, name, type, location) VALUES (?, ?, ?, ?)`
    res, err := r.db.ExecContext(ctx, query, sensor.UserID, sensor.Name, sensor.Type, sensor.Location)
    if err != nil {
        return err
    }
    id, _ := res.LastInsertId()
    sensor.ID = int(id)
    return nil
}

func (r *MySQLSensorRepository) GetByID(ctx context.Context, id int) (*domain.Sensor, error) {
    query := `SELECT id, user_id, name, type, location, created_at FROM sensors WHERE id = ?`
    row := r.db.QueryRowContext(ctx, query, id)
    var sensor domain.Sensor
    err := row.Scan(&sensor.ID, &sensor.UserID, &sensor.Name, &sensor.Type, &sensor.Location, &sensor.CreatedAt)
    return &sensor, err
}

func (r *MySQLSensorRepository) GetByUserID(ctx context.Context, userID int) ([]domain.Sensor, error) {
    query := `SELECT id, user_id, name, type, location, created_at FROM sensors WHERE user_id = ?`
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sensors []domain.Sensor
    for rows.Next() {
        var sensor domain.Sensor
        if err := rows.Scan(&sensor.ID, &sensor.UserID, &sensor.Name, &sensor.Type, &sensor.Location, &sensor.CreatedAt); err != nil {
            return nil, err
        }
        sensors = append(sensors, sensor)
    }
    return sensors, nil
}

func (r *MySQLSensorRepository) Update(ctx context.Context, sensor *domain.Sensor) error {
    query := `UPDATE sensors SET name = ?, type = ?, location = ? WHERE id = ?`
    _, err := r.db.ExecContext(ctx, query, sensor.Name, sensor.Type, sensor.Location, sensor.ID)
    return err
}

func (r *MySQLSensorRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM sensors WHERE id = ?`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}