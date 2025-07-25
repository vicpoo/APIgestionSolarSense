// api/src/login/infrastructure/auth_repository.go
package infrastructure

import (
	"context"
	"database/sql"
	"errors"
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

	// 1. Verificar si el usuario ya existe
	var existingID int64
	err = tx.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", email).Scan(&existingID)
	if err == nil {
		return 0, fmt.Errorf("user with this email already exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("could not check user existence: %w", err)
	}

	// 2. Insertar en users
	res, err := tx.ExecContext(ctx,
		`INSERT INTO users (email, display_name, username, password_hash, auth_type, created_at, last_login, is_active) 
         VALUES (?, ?, ?, ?, 'email', NOW(), NOW(), 1)`,
		email, username, username, passwordHash)
	if err != nil {
		return 0, fmt.Errorf("could not create user: %w", err)
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not get user ID: %w", err)
	}

	// 3. Insertar membresía solo si no existe
	_, err = tx.ExecContext(ctx,
		`INSERT INTO memberships (user_id, type, created_at, updated_at) 
         VALUES (?, 'free', NOW(), NOW())
         ON DUPLICATE KEY UPDATE updated_at = NOW()`, // Maneja el caso duplicado
		userID)
	if err != nil {
		return 0, fmt.Errorf("could not create membership: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("transaction failed: %w", err)
	}

	return userID, nil
}
func (r *AuthRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*domain.User, string, error) {
    var user domain.User
    var passwordHash string
    var isAdmin bool

    // Consulta mejorada con JOIN explícito y manejo de errores
    err := r.db.QueryRowContext(ctx,
        `SELECT u.id, u.email, COALESCE(ea.username, u.username) as username, 
         ea.password_hash, 
         CASE WHEN m.type = 'admin' THEN 1 ELSE 0 END as is_admin,
         u.is_active
         FROM users u
         LEFT JOIN email_auth ea ON u.id = ea.user_id
         LEFT JOIN memberships m ON u.id = m.user_id
         WHERE u.email = ? AND u.auth_type = 'email'`,
        email,
    ).Scan(&user.ID, &user.Email, &user.Username, &passwordHash, &isAdmin, &user.IsActive)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, "", errors.New("user not found")
        }
        return nil, "", fmt.Errorf("database error: %v", err)
    }

    if !user.IsActive {
        return nil, "", errors.New("account is not active")
    }

    if passwordHash == "" {
        return nil, "", errors.New("no password set for this user")
    }

    user.AuthType = "email"
    user.IsAdmin = isAdmin
    return &user, passwordHash, nil
}
func (r *AuthRepositoryImpl) FindUserByID(ctx context.Context, id int64) (*domain.User, string, error) {
	var user domain.User
	var passwordHash string

	err := r.db.QueryRowContext(ctx,
		`SELECT u.id, u.email, ea.username, ea.password_hash 
		 FROM users u
		 JOIN email_auth ea ON u.id = ea.user_id
		 WHERE u.id = ? AND u.auth_type = 'email'`,
		id,
	).Scan(&user.ID, &user.Email, &user.Username, &passwordHash)

	if err != nil {
		return nil, "", err
	}

	return &user, passwordHash, nil
}

func (r *AuthRepositoryImpl) UpdateLastLogin(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET last_login = NOW() WHERE id = ?", userID)
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

func (r *AuthRepositoryImpl) UpdatePassword(ctx context.Context, userID int64, newPasswordHash string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE email_auth SET password_hash = ? WHERE user_id = ?",
		newPasswordHash, userID)
	return err
}

func (r *AuthRepositoryImpl) UpsertGoogleUser(ctx context.Context, userData map[string]interface{}) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO users (uid, email, display_name, photo_url, provider, auth_type, last_login) 
              VALUES (?, ?, ?, ?, 'google', 'google', NOW())
              ON DUPLICATE KEY UPDATE 
              email = VALUES(email),
              display_name = VALUES(display_name),
              photo_url = VALUES(photo_url),
              last_login = NOW()`

	_, err = tx.ExecContext(ctx, query,
		userData["uid"],
		userData["email"],
		userData["displayName"],
		userData["photoURL"])
	if err != nil {
		return fmt.Errorf("could not upsert user: %w", err)
	}

	return tx.Commit()
}

func (r *AuthRepositoryImpl) DeleteUserByEmail(ctx context.Context, email string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	var userID int64
	err = tx.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM email_auth WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("could not delete from email_auth: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}

	return tx.Commit()
}

func (r *AuthRepositoryImpl) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	query := `
        SELECT 
            id, 
            uid, 
            email, 
            display_name, 
            username,
            photo_url, 
            provider, 
            auth_type, 
            is_active, 
            created_at, 
            last_login
        FROM users
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var (
			uid, username, photoURL sql.NullString
			lastLogin               sql.NullTime // Para manejar NULL en last_login
		)

		err := rows.Scan(
			&user.ID,
			&uid,
			&user.Email,
			&user.Username, // display_name como username
			&username,
			&photoURL,
			&user.Provider,
			&user.AuthType,
			&user.IsActive,
			&user.CreatedAt,
			&lastLogin, // Escaneamos como NullTime
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan user: %w", err)
		}

		// Manejar campos NULL
		if uid.Valid {
			user.UID = uid.String
		}
		if photoURL.Valid {
			user.PhotoURL = photoURL.String
		}
		if !username.Valid {
			user.Username = "sin_usuario"
		}
		if lastLogin.Valid {
			user.LastLogin = lastLogin.Time
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}
func (r *AuthRepositoryImpl) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	query := `
        SELECT 
            id, 
            uid, 
            email, 
            display_name, 
            photo_url, 
            provider, 
            auth_type, 
            is_active, 
            created_at, 
            last_login
        FROM users
        WHERE id = ?`

	var user domain.User
	var uid, photoURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&uid,
		&user.Email,
		&user.Username,
		&photoURL,
		&user.Provider,
		&user.AuthType,
		&user.IsActive,
		&user.CreatedAt,
		&user.LastLogin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("could not get user: %w", err)
	}

	if uid.Valid {
		user.UID = uid.String
	}
	if photoURL.Valid {
		user.PhotoURL = photoURL.String
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) GetUserMembershipType(ctx context.Context, userID int64) (string, error) {
	var membershipType string
	err := r.db.QueryRowContext(ctx,
		`SELECT type FROM memberships WHERE user_id = ?`,
		userID,
	).Scan(&membershipType)

	if err != nil {
		if err == sql.ErrNoRows {
			return "free", nil // Valor por defecto si no hay membresía
		}
		return "", err
	}

	return membershipType, nil
}
