// api/src/login/domain/repository.go
package domain

import "context"

type AuthRepository interface {
	CreateUserWithEmail(ctx context.Context, email, username, passwordHash string) (int64, error)
	FindUserByEmail(ctx context.Context, email string) (*User, string, error)
	FindUserByID(ctx context.Context, id int64) (*User, string, error)
	UpdateLastLogin(ctx context.Context, userID int64) error
	UpdateUserEmail(ctx context.Context, currentEmail, newEmail string) error
	UpdatePassword(ctx context.Context, userID int64, newPasswordHash string) error
	UpsertGoogleUser(ctx context.Context, userData map[string]interface{}) error
	DeleteUserByEmail(ctx context.Context, email string) error
	GetAllUsers(ctx context.Context) ([]*User, error)
	GetUserByID(ctx context.Context, userID int64) (*User, error)
	GetUserMembershipType(ctx context.Context, userID int64) (string, error)
	GetBySensorID(ctx context.Context, sensorID int) (*User, error) // Añadido este método
	EmailExists(ctx context.Context, email string) (bool, error)
	GetBasicUserInfo(ctx context.Context, email string) (*User, error)
	    UpdateDisplayName(ctx context.Context, userID int64, displayName string) error
    UpdateUserEmailById(ctx context.Context, userID int64, newEmail string) error
    UpdateUsername(ctx context.Context, userID int64, newUsername string) error
	   GetGoogleUserByUID(ctx context.Context, uid string) (*User, error)
    GetAllGoogleUsers(ctx context.Context) ([]*User, error)
}