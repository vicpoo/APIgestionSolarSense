// src/notification_settings/application/post_settings_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/domain"
)

type PostSettingsUseCase struct {
    repo domain.SettingsRepository
}

func NewPostSettingsUseCase(repo domain.SettingsRepository) *PostSettingsUseCase {
    return &PostSettingsUseCase{repo: repo}
}

func (uc *PostSettingsUseCase) CreateSettings(ctx context.Context, settings *domain.NotificationSettings) error {
    return uc.repo.Update(ctx, settings) // Reutilizamos Update para creación
}