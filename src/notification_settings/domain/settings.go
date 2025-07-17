//api\src\notification_settings\domain\settings.go
package domain

type NotificationSettings struct {
    ID                int  `json:"id"`
    UserID            int  `json:"user_id"`
    EmailAlerts       bool `json:"email_alerts"`
    PushNotifications bool `json:"push_notifications"`
}