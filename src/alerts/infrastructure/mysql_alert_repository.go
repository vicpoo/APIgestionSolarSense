//src/alerts/infrastructure/persistence/mysql_alert_repository.go

package infrastructure

import (
    "context"
    "database/sql"
    "github.com/vicpoo/apigestion-solar-go/src/core"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/domain"
 
)

type MySQLAlertRepository struct {
    db *sql.DB
}

func NewMySQLAlertRepository() domain.AlertRepository {
    return &MySQLAlertRepository{db: core.GetBD()}
}

func (r *MySQLAlertRepository) Create(ctx context.Context, alert *domain.Alert) error {
    query := `INSERT INTO alerts 
        (sensor_id, message, alert_type, is_sent) 
        VALUES (?, ?, ?, ?)`
    _, err := r.db.ExecContext(ctx, query,
        alert.SensorID,
        alert.Message,
        alert.AlertType,
        alert.IsSent,
    )
    return err
}

func (r *MySQLAlertRepository) GetByID(ctx context.Context, id int) (*domain.Alert, error) {
    query := `SELECT id, sensor_id, message, triggered_at, alert_type, is_sent 
              FROM alerts WHERE id = ?`
    row := r.db.QueryRowContext(ctx, query, id)
    
    var alert domain.Alert
    err := row.Scan(
        &alert.ID,
        &alert.SensorID,
        &alert.Message,
        &alert.TriggeredAt,
        &alert.AlertType,
        &alert.IsSent,
    )
    if err != nil {
        return nil, err
    }
    return &alert, nil
}

func (r *MySQLAlertRepository) GetBySensorID(ctx context.Context, sensorID int, limit int) ([]domain.Alert, error) {
    query := `SELECT id, sensor_id, message, triggered_at, alert_type, is_sent 
              FROM alerts WHERE sensor_id = ? 
              ORDER BY triggered_at DESC 
              LIMIT ?`
    rows, err := r.db.QueryContext(ctx, query, sensorID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var alerts []domain.Alert
    for rows.Next() {
        var alert domain.Alert
        err := rows.Scan(
            &alert.ID,
            &alert.SensorID,
            &alert.Message,
            &alert.TriggeredAt,
            &alert.AlertType,
            &alert.IsSent,
        )
        if err != nil {
            return nil, err
        }
        alerts = append(alerts, alert)
    }
    return alerts, nil
}

func (r *MySQLAlertRepository) GetUnsent(ctx context.Context) ([]domain.Alert, error) {
    query := `SELECT id, sensor_id, message, triggered_at, alert_type, is_sent 
              FROM alerts WHERE is_sent = 0`
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var alerts []domain.Alert
    for rows.Next() {
        var alert domain.Alert
        err := rows.Scan(
            &alert.ID,
            &alert.SensorID,
            &alert.Message,
            &alert.TriggeredAt,
            &alert.AlertType,
            &alert.IsSent,
        )
        if err != nil {
            return nil, err
        }
        alerts = append(alerts, alert)
    }
    return alerts, nil
}

func (r *MySQLAlertRepository) MarkAsSent(ctx context.Context, id int) error {
    query := `UPDATE alerts SET is_sent = 1 WHERE id = ?`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}