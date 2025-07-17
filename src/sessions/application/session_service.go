// src/sessions/application/session_service.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/domain"
)

type SessionService struct {
    repo domain.SessionRepository
}

func NewSessionService(repo domain.SessionRepository) *SessionService {
    return &SessionService{repo: repo}
}

func (s *SessionService) CreateSession(ctx context.Context, session *domain.Session) error {
    return s.repo.Create(ctx, session)
}

func (s *SessionService) ValidateSession(ctx context.Context, token string) (*domain.Session, error) {
    return s.repo.GetByToken(ctx, token)
}

func (s *SessionService) InvalidateSession(ctx context.Context, token string) error {
    return s.repo.Invalidate(ctx, token)
}

func (s *SessionService) DeleteSession(ctx context.Context, id int) error {
    return s.repo.Delete(ctx, id)
}