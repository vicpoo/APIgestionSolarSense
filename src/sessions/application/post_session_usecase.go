// src/sessions/application/post_session_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/domain"
)

type PostSessionUseCase struct {
    repo domain.SessionRepository
}

func NewPostSessionUseCase(repo domain.SessionRepository) *PostSessionUseCase {
    return &PostSessionUseCase{repo: repo}
}

func (uc *PostSessionUseCase) CreateSession(ctx context.Context, session *domain.Session) error {
    return uc.repo.Create(ctx, session)
}