// api/src/login/infrastructure/auth_repository.go
package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
)

type AuthRepositoryImpl struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return &AuthRepositoryImpl{db: db}
}

func (r *AuthRepositoryImpl) CreateUserWithEmail(ctx context.Context, email, username, passwordHash string) (int64, error) {
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

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("transaction failed: %w", err)
	}

	return userID, nil
}

func (r *AuthRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*domain.User, string, error) {
	var user domain.User
	var passwordHash string
	
	err := r.db.QueryRowContext(ctx,
		`SELECT ea.user_id, ea.username, ea.password_hash 
		 FROM email_auth ea
		 JOIN users u ON ea.user_id = u.id
		 WHERE ea.email = ? AND u.auth_type = 'email' AND u.is_active = 1`,
		email,
	).Scan(&user.ID, &user.Username, &passwordHash)

	if err != nil {
		return nil, "", err
	}

	return &user, passwordHash, nil
}

func (r *AuthRepositoryImpl) UpdateLastLogin(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET last_login = NOW() WHERE id = ?", userID)
	return err
}

func (r *AuthRepositoryImpl) UpsertGoogleUser(ctx context.Context, userData map[string]interface{}) error {
	query := `INSERT INTO users (uid, email, display_name, photo_url, provider, last_login, auth_type) 
	          VALUES (?, ?, ?, ?, 'google', NOW(), 'google')
	          ON DUPLICATE KEY UPDATE 
	          email = VALUES(email),
	          display_name = VALUES(display_name),
	          photo_url = VALUES(photo_url),
	          last_login = NOW()`

	_, err := r.db.ExecContext(ctx, query,
		userData["uid"],
		userData["email"],
		userData["displayName"],
		userData["photoURL"])
	
	return err
}

func (r *AuthRepositoryImpl) UpdateUserEmail(ctx context.Context, currentEmail, newEmail string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	// Verificar si el nuevo email ya existe
	var count int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE email = ?", newEmail).Scan(&count)
	if err != nil {
		return fmt.Errorf("could not check email existence: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("email already in use")
	}

	// Actualizar en tabla users
	_, err = tx.ExecContext(ctx, "UPDATE users SET email = ? WHERE email = ?", newEmail, currentEmail)
	if err != nil {
		return fmt.Errorf("could not update user email: %w", err)
	}

	// Actualizar en tabla email_auth
	_, err = tx.ExecContext(ctx, "UPDATE email_auth SET email = ? WHERE email = ?", newEmail, currentEmail)
	if err != nil {
		return fmt.Errorf("could not update email_auth: %w", err)
	}

	return tx.Commit()
}

func (r *AuthRepositoryImpl) DeleteUserByEmail(ctx context.Context, email string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	// Primero obtener el user_id para las eliminaciones en cascada
	var userID int64
	err = tx.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Eliminar de email_auth
	_, err = tx.ExecContext(ctx, "DELETE FROM email_auth WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("could not delete from email_auth: %w", err)
	}

	// Eliminar de users (esto debería activar las eliminaciones en cascada para otras tablas)
	_, err = tx.ExecContext(ctx, "DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}

	return tx.Commit()
}