// src/sessions/application/get_session_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/domain"
)

type GetSessionUseCase struct {
    repo domain.SessionRepository
}

func NewGetSessionUseCase(repo domain.SessionRepository) *GetSessionUseCase {
    return &GetSessionUseCase{repo: repo}
}

func (uc *GetSessionUseCase) ValidateSession(ctx context.Context, token string) (*domain.Session, error) {
    return uc.repo.GetByToken(ctx, token)
}