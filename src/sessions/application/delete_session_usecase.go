// src/sessions/application/delete_session_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/domain"
)

type DeleteSessionUseCase struct {
    repo domain.SessionRepository
}

func NewDeleteSessionUseCase(repo domain.SessionRepository) *DeleteSessionUseCase {
    return &DeleteSessionUseCase{repo: repo}
}

func (uc *DeleteSessionUseCase) InvalidateSession(ctx context.Context, token string) error {
    return uc.repo.Invalidate(ctx, token)
}

func (uc *DeleteSessionUseCase) Delete(ctx context.Context, id int) error {
    return uc.repo.Delete(ctx, id)
}