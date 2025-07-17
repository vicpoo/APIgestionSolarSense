// src/sessions/domain/session_repository.go
package domain

import (
    "context"
    "time"
)

type SessionRepository interface {
    Create(ctx context.Context, session *Session) error
    GetByToken(ctx context.Context, token string) (*Session, error)
    UpdateExpiry(ctx context.Context, token string, newExpiry time.Time) error
    Invalidate(ctx context.Context, token string) error
    Delete(ctx context.Context, id int) error
}