// src/notification_settings/application/put_settings_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/domain"
)

type PutSettingsUseCase struct {
    repo domain.SettingsRepository
}

func NewPutSettingsUseCase(repo domain.SettingsRepository) *PutSettingsUseCase {
    return &PutSettingsUseCase{repo: repo}
}

func (uc *PutSettingsUseCase) UpdateSettings(ctx context.Context, settings *domain.NotificationSettings) error {
    return uc.repo.Update(ctx, settings)
}