//src/memberships/infrastructure/mysql_membership_repository.go

package infrastructure

import (
    "context"
    "database/sql"
    "github.com/vicpoo/apigestion-solar-go/src/core"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/domain"

)

type MySQLMembershipRepository struct {
    db *sql.DB
}

func NewMySQLMembershipRepository() domain.MembershipRepository {
    return &MySQLMembershipRepository{db: core.GetBD()}
}

func (r *MySQLMembershipRepository) GetByUserID(ctx context.Context, userID int) (*domain.Membership, error) {
    query := `SELECT id, user_id, type, extra_storage, created_at 
              FROM memberships WHERE user_id = ?`
    row := r.db.QueryRowContext(ctx, query, userID)
    
    var membership domain.Membership
    err := row.Scan(
        &membership.ID,
        &membership.UserID,
        &membership.Type,
        &membership.ExtraStorage,
        &membership.CreatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &membership, nil
}

func (r *MySQLMembershipRepository) CreateOrUpdate(ctx context.Context, membership *domain.Membership) error {
    query := `INSERT INTO memberships 
        (user_id, type, extra_storage) 
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE
        type = VALUES(type),
        extra_storage = VALUES(extra_storage)`
    _, err := r.db.ExecContext(ctx, query,
        membership.UserID,
        membership.Type,
        membership.ExtraStorage,
    )
    return err
}

func (r *MySQLMembershipRepository) UpgradeToPremium(ctx context.Context, userID int) error {
    query := `UPDATE memberships SET type = 'premium' WHERE user_id = ?`
    _, err := r.db.ExecContext(ctx, query, userID)
    return err
}

func (r *MySQLMembershipRepository) DowngradeToFree(ctx context.Context, userID int) error {
    query := `UPDATE memberships SET type = 'free', extra_storage = 0 WHERE user_id = ?`
    _, err := r.db.ExecContext(ctx, query, userID)
    return err
}

func (r *MySQLMembershipRepository) Delete(ctx context.Context, userID int) error {
    query := `DELETE FROM memberships WHERE user_id = ?`
    _, err := r.db.ExecContext(ctx, query, userID)
    return err
}