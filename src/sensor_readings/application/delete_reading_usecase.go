// src/sensor_readings/application/delete_reading_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/domain"
)

type DeleteReadingUseCase struct {
    repo domain.ReadingRepository
}

func NewDeleteReadingUseCase(repo domain.ReadingRepository) *DeleteReadingUseCase {
    return &DeleteReadingUseCase{repo: repo}
}

func (uc *DeleteReadingUseCase) DeleteReading(ctx context.Context, id int) error {
    return uc.repo.Delete(ctx, id)
}