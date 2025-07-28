// api/src/login/domain/service.go
package domain

import "context"

type AuthService interface {
    RegisterWithEmail(ctx context.Context, creds UserCredentials) (*AuthResponse, error)
    LoginWithEmail(ctx context.Context, creds UserCredentials) (*AuthResponse, error)
    AuthenticateWithGoogle(ctx context.Context, idToken string) (*AuthResponse, error)
    UpdateUserProfile(ctx context.Context, userID int64, email, username, displayName, authType string) (*AuthResponse, error)
	UpdateGoogleUserProfile(ctx context.Context, userID int64, displayName string) (*AuthResponse, error)
}