// src/notification_settings/application/get_settings_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/domain"
)

type GetSettingsUseCase struct {
    repo domain.SettingsRepository
}

func NewGetSettingsUseCase(repo domain.SettingsRepository) *GetSettingsUseCase {
    return &GetSettingsUseCase{repo: repo}
}

func (uc *GetSettingsUseCase) GetSettings(ctx context.Context, userID int) (*domain.NotificationSettings, error) {
    return uc.repo.GetByUserID(ctx, userID)
}