//api\src\notification_settings\infrastructure\mysql_settings_repository.go

package infrastructure

import (
    "context"
    "database/sql"
    "github.com/vicpoo/apigestion-solar-go/src/core"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/domain"

)

type MySQLSettingsRepository struct {
    db *sql.DB
}

func NewMySQLSettingsRepository() domain.SettingsRepository {
    return &MySQLSettingsRepository{db: core.GetBD()}
}

func (r *MySQLSettingsRepository) GetByUserID(ctx context.Context, userID int) (*domain.NotificationSettings, error) {
    query := `SELECT id, user_id, email_alerts, push_notifications FROM notification_settings WHERE user_id = ?`
    row := r.db.QueryRowContext(ctx, query, userID)
    
    var settings domain.NotificationSettings
    err := row.Scan(
        &settings.ID,
        &settings.UserID,
        &settings.EmailAlerts,
        &settings.PushNotifications,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &settings, nil
}

func (r *MySQLSettingsRepository) Update(ctx context.Context, settings *domain.NotificationSettings) error {
    query := `INSERT INTO notification_settings 
        (user_id, email_alerts, push_notifications) 
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE
        email_alerts = VALUES(email_alerts),
        push_notifications = VALUES(push_notifications)`
    
    _, err := r.db.ExecContext(ctx, query,
        settings.UserID,
        settings.EmailAlerts,
        settings.PushNotifications,
    )
    return err
}

func (r *MySQLSettingsRepository) Delete(ctx context.Context, userID int) error {
    query := `DELETE FROM notification_settings WHERE user_id = ?`
    _, err := r.db.ExecContext(ctx, query, userID)
    return err
}