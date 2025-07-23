// src/memberships/infrastructure/mysql_membership_repository.go
package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentType string
	err = tx.QueryRowContext(ctx, "SELECT type FROM memberships WHERE user_id = ?", membership.UserID).Scan(&currentType)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	query := `INSERT INTO memberships 
		(user_id, type, extra_storage, created_at) 
		VALUES (?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
		type = VALUES(type),
		extra_storage = VALUES(extra_storage),
		updated_at = NOW()`
	
	_, err = tx.ExecContext(ctx, query,
		membership.UserID,
		membership.Type,
		membership.ExtraStorage)
	if err != nil {
		return err
	}

	if currentType != "" {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO membership_changes 
			 (user_id, old_role, new_role, changed_at) 
			 VALUES (?, ?, ?, NOW())`,
			membership.UserID,
			currentType,
			membership.Type)
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
	var changedBy int = 1
	
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if claims, exists := ginCtx.Get("userClaims"); exists {
			if claimsMap, ok := claims.(map[string]interface{}); ok {
				if id, ok := claimsMap["user_id"].(float64); ok {
					changedBy = int(id)
				}
			}
		}
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentType string
	err = tx.QueryRowContext(ctx, "SELECT type FROM memberships WHERE user_id = ?", userID).Scan(&currentType)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

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

	if currentType != "" {
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

func (r *MySQLMembershipRepository) UpdateUser(ctx context.Context, userID int, email, username, passwordHash *string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Actualizar usuario en tabla users
	if email != nil || username != nil {
		query := "UPDATE users SET "
		var updates []string
		var args []interface{}

		if email != nil {
			updates = append(updates, "email = ?")
			args = append(args, *email)
		}
		if username != nil {
			updates = append(updates, "display_name = ?")
			args = append(args, *username)
		}

		query += strings.Join(updates, ", ") + " WHERE id = ?"
		args = append(args, userID)

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("could not update user: %w", err)
		}
	}

	// Actualizar credenciales si es usuario de email y se proporcionó password
	if passwordHash != nil {
		var authType string
		err := tx.QueryRowContext(ctx, "SELECT auth_type FROM users WHERE id = ?", userID).Scan(&authType)
		if err != nil {
			return fmt.Errorf("could not get user auth type: %w", err)
		}

		if authType == "email" {
			_, err = tx.ExecContext(ctx,
				"UPDATE email_auth SET password_hash = ? WHERE user_id = ?",
				*passwordHash, userID)
			if err != nil {
				return fmt.Errorf("could not update password: %w", err)
			}
		}
	}

	return tx.Commit()
}



func (r *MySQLMembershipRepository) Delete(ctx context.Context, userID int) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 1. Primero eliminar los registros en membership_changes
    _, err = tx.ExecContext(ctx, "DELETE FROM membership_changes WHERE user_id = ?", userID)
    if err != nil {
        return fmt.Errorf("could not delete membership changes: %w", err)
    }

    // 2. Eliminar la membresía
    _, err = tx.ExecContext(ctx, "DELETE FROM memberships WHERE user_id = ?", userID)
    if err != nil {
        return fmt.Errorf("could not delete membership: %w", err)
    }

    // 3. Eliminar credenciales de email si existe
    _, err = tx.ExecContext(ctx, "DELETE FROM email_auth WHERE user_id = ?", userID)
    if err != nil {
        return fmt.Errorf("could not delete email auth: %w", err)
    }

    // 4. Finalmente eliminar el usuario
    _, err = tx.ExecContext(ctx, "DELETE FROM users WHERE id = ?", userID)
    if err != nil {
        return fmt.Errorf("could not delete user: %w", err)
    }

    return tx.Commit()
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
			COALESCE(m.type, 'free') AS membership_type,
			COALESCE(m.extra_storage, 0) AS extra_storage,
			COALESCE(m.created_at, u.created_at) AS membership_since
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
		var photoURL sql.NullString
		
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.DisplayName,
			&photoURL,
			&user.Provider,
			&user.IsActive,
			&user.MembershipType,
			&user.ExtraStorage,
			&user.MembershipSince,
		)
		if err != nil {
			return nil, err
		}
		
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

func (r *MySQLMembershipRepository) RegisterUser(ctx context.Context, email, username, passwordHash string) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`INSERT INTO users (email, display_name, auth_type, created_at, last_login) 
		 VALUES (?, ?, 'email', NOW(), NOW())`,
		email, username)
	if err != nil {
		return 0, fmt.Errorf("could not create user: %w", err)
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not get user ID: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO email_auth 
		 (user_id, email, password_hash, username, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, NOW(), NOW())`,
		userID, email, passwordHash, username)
	if err != nil {
		return 0, fmt.Errorf("could not save credentials: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO memberships 
		 (user_id, type, extra_storage, created_at) 
		 VALUES (?, 'free', false, NOW())`,
		userID)
	if err != nil {
		return 0, fmt.Errorf("could not create membership: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("transaction failed: %w", err)
	}

	return userID, nil
}

// Nuevo método para corregir membresías faltantes
func (r *MySQLMembershipRepository) FixMissingMemberships(ctx context.Context) error {
	query := `
		INSERT INTO memberships (user_id, type, extra_storage, created_at)
		SELECT id, 'free', 0, NOW() FROM users
		WHERE id NOT IN (SELECT user_id FROM memberships)
		AND id != 1`
	
	_, err := r.db.ExecContext(ctx, query)
	return err
}