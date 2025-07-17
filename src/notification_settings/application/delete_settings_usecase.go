// src/notification_settings/application/delete_settings_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/notification_settings/domain"
)

type DeleteSettingsUseCase struct {
    repo domain.SettingsRepository
}

func NewDeleteSettingsUseCase(repo domain.SettingsRepository) *DeleteSettingsUseCase {
    return &DeleteSettingsUseCase{repo: repo}
}

func (uc *DeleteSettingsUseCase) DeleteSettings(ctx context.Context, userID int) error {
    return uc.repo.Delete(ctx, userID)
}