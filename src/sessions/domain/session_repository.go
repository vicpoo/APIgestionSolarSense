//api/src/sessions/domain/session_repository.go

package domain

import (
    "context"

)

type SessionRepository interface {
    Create(ctx context.Context, session *Session) error
    GetByToken(ctx context.Context, token string) (*Session, error)
    Invalidate(ctx context.Context, token string) error
}