// src/notification_settings/application/settings_service.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/domain"
)

type SettingsService struct {
    repo domain.SettingsRepository
}

func NewSettingsService(repo domain.SettingsRepository) *SettingsService {
    return &SettingsService{repo: repo}
}

func (s *SettingsService) GetSettings(ctx context.Context, userID int) (*domain.NotificationSettings, error) {
    return s.repo.GetByUserID(ctx, userID)
}

func (s *SettingsService) UpdateSettings(ctx context.Context, settings *domain.NotificationSettings) error {
    return s.repo.Update(ctx, settings)
}

func (s *SettingsService) DeleteSettings(ctx context.Context, userID int) error {
    return s.repo.Delete(ctx, userID)
}