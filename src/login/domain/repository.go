// api/src/login/domain/repository.go
package domain

import "context"

type AuthRepository interface {
	CreateUserWithEmail(ctx context.Context, email, username, passwordHash string) (int64, error)
	FindUserByEmail(ctx context.Context, email string) (*User, string, error)
	UpdateLastLogin(ctx context.Context, userID int64) error
	UpsertGoogleUser(ctx context.Context, userData map[string]interface{}) error
	UpdateUserEmail(ctx context.Context, currentEmail, newEmail string) error
	DeleteUserByEmail(ctx context.Context, email string) error
	GetAllUsers(ctx context.Context) ([]*User, error)
    GetUserByID(ctx context.Context, userID int64) (*User, error)
}