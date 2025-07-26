// src/alerts/application/alert_service.go
package application

import (
	"context"
	"fmt"
	"log"

	"github.com/vicpoo/apigestion-solar-go/src/alerts/domain"
	"github.com/vicpoo/apigestion-solar-go/src/email"
	nsdomain "github.com/vicpoo/apigestion-solar-go/src/notification_settings/domain"
	udomain "github.com/vicpoo/apigestion-solar-go/src/login/domain"
)

type AlertService struct {
	repo            domain.AlertRepository
	userRepo        udomain.AuthRepository
	settingsRepo    nsdomain.SettingsRepository
	emailService    *email.EmailService
}

func NewAlertService(
	repo domain.AlertRepository,
	userRepo udomain.AuthRepository,
	settingsRepo nsdomain.SettingsRepository,
	emailService *email.EmailService,
) *AlertService {
	return &AlertService{
		repo:         repo,
		userRepo:     userRepo,
		settingsRepo: settingsRepo,
		emailService: emailService,
	}
}

func (s *AlertService) CreateAlert(ctx context.Context, alert *domain.Alert) error {
	err := s.repo.Create(ctx, alert)
	if err != nil {
		return err
	}

	// Obtener usuario dueño del sensor
	user, err := s.userRepo.GetBySensorID(ctx, alert.SensorID)
	if err != nil {
		log.Printf("Error getting user for sensor %d: %v", alert.SensorID, err)
		return nil // No fallar la creación por error de notificación
	}

	// Convertir user.ID de int64 a int para GetByUserID
	userID := int(user.ID)

	// Obtener configuración de notificaciones
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil || settings == nil {
		log.Printf("Error getting notification settings for user %d: %v", userID, err)
		return nil
	}

	// Enviar email si está habilitado
	if settings.EmailAlerts {
		subject := fmt.Sprintf("Alerta de %s", alert.AlertType)
		body := fmt.Sprintf("Se ha generado una alerta:\n\n%s\n\nTipo: %s\nSensor ID: %d\nFecha: %s",
			alert.Message, alert.AlertType, alert.SensorID, alert.TriggeredAt.Format("2006-01-02 15:04:05"))

		err = s.emailService.SendAlertEmail(user.Email, subject, body)
		if err != nil {
			log.Printf("Error sending email to %s: %v", user.Email, err)
		} else {
			// Marcar alerta como enviada
			err = s.repo.MarkAsSent(ctx, alert.ID)
			if err != nil {
				log.Printf("Error marking alert %d as sent: %v", alert.ID, err)
			}
		}
	}

	return nil
}

func (s *AlertService) GetAlert(ctx context.Context, id int) (*domain.Alert, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AlertService) GetSensorAlerts(ctx context.Context, sensorID int, limit int) ([]domain.Alert, error) {
	return s.repo.GetBySensorID(ctx, sensorID, limit)
}

func (s *AlertService) GetUnsent(ctx context.Context) ([]domain.Alert, error) {
	return s.repo.GetUnsent(ctx)
}

func (s *AlertService) MarkAlertAsSent(ctx context.Context, id int) error {
	return s.repo.MarkAsSent(ctx, id)
}

func (s *AlertService) UpdateAlert(ctx context.Context, alert *domain.Alert) error {
	return s.repo.Update(ctx, alert)
}

func (s *AlertService) DeleteAlert(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}