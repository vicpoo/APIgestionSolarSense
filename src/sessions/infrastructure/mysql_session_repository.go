//src/sessions/infrastructure/persistence/mysql_session_repository.go

package infrastructure

import (
    "context"
    "database/sql"
    "github.com/vicpoo/apigestion-solar-go/src/core"
    "github.com/vicpoo/apigestion-solar-go/src/sessions/domain"

)

type MySQLSessionRepository struct {
    db *sql.DB
}

func NewMySQLSessionRepository() domain.SessionRepository {
    return &MySQLSessionRepository{db: core.GetBD()}
}

func (r *MySQLSessionRepository) Create(ctx context.Context, session *domain.Session) error {
    query := `INSERT INTO sessions 
        (user_id, session_token, firebase_uid, expires_at, is_active) 
        VALUES (?, ?, ?, ?, ?)`
    _, err := r.db.ExecContext(ctx, query,
        session.UserID,
        session.SessionToken,
        session.FirebaseUID,
        session.ExpiresAt,
        session.IsActive,
    )
    return err
}

func (r *MySQLSessionRepository) GetByToken(ctx context.Context, token string) (*domain.Session, error) {
    query := `SELECT id, user_id, session_token, firebase_uid, created_at, expires_at, is_active 
              FROM sessions 
              WHERE session_token = ? AND is_active = 1 AND expires_at > NOW()`
    row := r.db.QueryRowContext(ctx, query, token)
    
    var session domain.Session
    err := row.Scan(
        &session.ID,
        &session.UserID,
        &session.SessionToken,
        &session.FirebaseUID,
        &session.CreatedAt,
        &session.ExpiresAt,
        &session.IsActive,
    )
    if err != nil {
        return nil, err
    }
    return &session, nil
}

func (r *MySQLSessionRepository) Invalidate(ctx context.Context, token string) error {
    query := `UPDATE sessions SET is_active = 0 WHERE session_token = ?`
    _, err := r.db.ExecContext(ctx, query, token)
    return err
}