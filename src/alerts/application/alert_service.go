//src/alerts/application/alert_service.go

package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/domain"
   
)

type AlertService struct {
    repo domain.AlertRepository
}

func NewAlertService(repo domain.AlertRepository) *AlertService {
    return &AlertService{repo: repo}
}

func (s *AlertService) CreateAlert(ctx context.Context, alert *domain.Alert) error {
    return s.repo.Create(ctx, alert)
}

func (s *AlertService) GetAlert(ctx context.Context, id int) (*domain.Alert, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *AlertService) GetSensorAlerts(ctx context.Context, sensorID int, limit int) ([]domain.Alert, error) {
    return s.repo.GetBySensorID(ctx, sensorID, limit)
}

func (s *AlertService) GetUnsentAlerts(ctx context.Context) ([]domain.Alert, error) {
    return s.repo.GetUnsent(ctx)
}

func (s *AlertService) MarkAlertAsSent(ctx context.Context, id int) error {
    return s.repo.MarkAsSent(ctx, id)
}