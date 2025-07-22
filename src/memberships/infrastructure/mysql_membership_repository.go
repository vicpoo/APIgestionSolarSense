// src/memberships/infrastructure/mysql_membership_repository.go
package infrastructure

import (
    "context"
    "database/sql"
    "github.com/gin-gonic/gin"
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
    // Iniciar transacción
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 1. Obtener el rol actual para registrar el cambio
    var currentType string
    err = tx.QueryRowContext(ctx, "SELECT type FROM memberships WHERE user_id = ?", membership.UserID).Scan(&currentType)
    if err != nil && err != sql.ErrNoRows {
        return err
    }

    // 2. Actualizar la membresía
    query := `INSERT INTO memberships 
        (user_id, type, extra_storage) 
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE
        type = VALUES(type),
        extra_storage = VALUES(extra_storage)`
    
    _, err = tx.ExecContext(ctx, query,
        membership.UserID,
        membership.Type,
        membership.ExtraStorage,
    )
    if err != nil {
        return err
    }

    // 3. Registrar el cambio en membership_changes (opcional, sin changed_by)
    if currentType != "" {
        _, err = tx.ExecContext(ctx,
            `INSERT INTO membership_changes 
             (user_id, old_role, new_role) 
             VALUES (?, ?, ?)`,
            membership.UserID,
            currentType,
            membership.Type,
        )
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}
func (r *MySQLMembershipRepository) UpgradeToPremium(ctx context.Context, userID int) error {
    return r.changeMembershipType(ctx, userID, "premium")
}

func (r *MySQLMembershipRepository) DowngradeToFree(ctx context.Context, userID int) error {
    return r.changeMembershipType(ctx, userID, "free")
}

func (r *MySQLMembershipRepository) changeMembershipType(ctx context.Context, userID int, newType string) error {
    // Obtener el ID del admin del contexto Gin
    var changedBy int = 1 // Valor por defecto (admin principal)
    
    if ginCtx, ok := ctx.(*gin.Context); ok {
        if claims, exists := ginCtx.Get("userClaims"); exists {
            if claimsMap, ok := claims.(map[string]interface{}); ok {
                if id, ok := claimsMap["user_id"].(float64); ok {
                    changedBy = int(id)
                }
            }
        }
    }

    // Iniciar transacción
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 1. Obtener el rol actual
    var currentType string
    err = tx.QueryRowContext(ctx, "SELECT type FROM memberships WHERE user_id = ?", userID).Scan(&currentType)
    if err != nil && err != sql.ErrNoRows {
        return err
    }

    // 2. Actualizar la membresía
    extraStorage := 0
    if newType == "premium" {
        extraStorage = 1
    }
    
    _, err = tx.ExecContext(ctx, 
        "UPDATE memberships SET type = ?, extra_storage = ? WHERE user_id = ? AND user_id != 1",
        newType, extraStorage, userID)
    if err != nil {
        return err
    }

    // 3. Registrar el cambio
    if currentType != "" { // Solo si existía un registro previo
        _, err = tx.ExecContext(ctx,
            `INSERT INTO membership_changes 
             (user_id, old_role, new_role, changed_by) 
             VALUES (?, ?, ?, ?)`,
            userID,
            currentType,
            newType,
            changedBy,
        )
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (r *MySQLMembershipRepository) Delete(ctx context.Context, userID int) error {
    query := `DELETE FROM memberships WHERE user_id = ? AND user_id != 1`
    _, err := r.db.ExecContext(ctx, query, userID)
    return err
}


func (r *MySQLMembershipRepository) GetAllUsers(ctx context.Context) ([]*domain.UserWithMembership, error) {
    query := `
        SELECT 
            u.id, 
            u.email, 
            u.display_name, 
            u.photo_url,
            u.provider,
            u.is_active,
            m.type AS membership_type,
            m.extra_storage,
            m.created_at AS membership_since
        FROM 
            users u
        LEFT JOIN 
            memberships m ON u.id = m.user_id
        ORDER BY u.id`

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*domain.UserWithMembership
    for rows.Next() {
        var user domain.UserWithMembership
        var photoURL sql.NullString // Usamos sql.NullString para campos que pueden ser NULL
        
        err := rows.Scan(
            &user.ID,
            &user.Email,
            &user.DisplayName,
            &photoURL,    // Escaneamos a sql.NullString
            &user.Provider,
            &user.IsActive,
            &user.MembershipType,
            &user.ExtraStorage,
            &user.MembershipSince,
        )
        if err != nil {
            return nil, err
        }
        
        // Convertimos sql.NullString a *string
        if photoURL.Valid {
            user.PhotoURL = &photoURL.String
        } else {
            user.PhotoURL = nil
        }
        
        users = append(users, &user)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}