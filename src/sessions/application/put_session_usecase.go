// src/sessions/application/put_session_usecase.go
package application

import (
    "context"
    "time"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/domain"
)

type PutSessionUseCase struct {
    repo domain.SessionRepository
}

func NewPutSessionUseCase(repo domain.SessionRepository) *PutSessionUseCase {
    return &PutSessionUseCase{repo: repo}
}

func (uc *PutSessionUseCase) RefreshSession(ctx context.Context, token string, newExpiry time.Time) error {
    return uc.repo.UpdateExpiry(ctx, token, newExpiry)
}