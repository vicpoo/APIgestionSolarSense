// api/src/login/application/auth_service.go
package application

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	repo domain.AuthRepository
}

func NewAuthService(repo domain.AuthRepository) domain.AuthService {
	return &AuthServiceImpl{repo: repo}
}

func (s *AuthServiceImpl) RegisterWithEmail(ctx context.Context, creds domain.UserCredentials) (*domain.AuthResponse, error) {
	// Validaciones adicionales
	if len(creds.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	if len(creds.Username) < 3 {
		return nil, errors.New("username must be at least 3 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("could not hash password")
	}

	userID, err := s.repo.CreateUserWithEmail(ctx, creds.Email, creds.Username, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	return &domain.AuthResponse{
		Success:  true,
		Message:  "User registered successfully",
		UserID:   userID,
		Email:    creds.Email,
		Username: creds.Username,
	}, nil
}
func (s *AuthServiceImpl) LoginWithEmail(ctx context.Context, creds domain.UserCredentials) (*domain.AuthResponse, error) {
	if creds.Email == "" || creds.Password == "" {
		return nil, errors.New("email and password are required")
	}

	user, passwordHash, err := s.repo.FindUserByEmail(ctx, creds.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, errors.New("could not update last login")
	}

	token, err := domain.GenerateJWTToken(user)
	if err != nil {
		return nil, errors.New("could not generate token")
	}

	return &domain.AuthResponse{
		Success:  true,
		Message:  "Login successful",
		Token:    token,
		AuthType: "email",
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}, nil
}

func (s *AuthServiceImpl) AuthenticateWithGoogle(ctx context.Context, idToken string) (*domain.AuthResponse, error) {
	userData, err := decodeTokenWithoutVerification(idToken)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if err := s.repo.UpsertGoogleUser(ctx, userData); err != nil {
		return nil, errors.New("could not save user data")
	}

	return &domain.AuthResponse{
		Success: true,
		Message: "Authentication successful",
	}, nil
}

func decodeTokenWithoutVerification(idToken string) (map[string]interface{}, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token claims")
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims")
	}

	userData := map[string]interface{}{
		"uid":         claims["sub"],
		"email":       getClaimValue(claims, "email", ""),
		"displayName": getClaimValue(claims, "name", ""),
		"photoURL":    getClaimValue(claims, "picture", ""),
	}

	return userData, nil
}

func getClaimValue(claims map[string]interface{}, key string, defaultValue string) string {
	if value, ok := claims[key]; ok {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return defaultValue
}

func (s *AuthServiceImpl) UpdateGoogleUserProfile(ctx context.Context, userID int64, displayName string) (*domain.AuthResponse, error) {
    // Obtener el usuario actual para verificar que es de Google
    user, err := s.repo.GetUserByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    // Verificar que el auth_type sea google
    if user.AuthType != "google" {
        return nil, errors.New("only google users can be updated with this endpoint")
    }

    // Actualizar solo el display_name
    if displayName != "" {
        err = s.repo.UpdateDisplayName(ctx, userID, displayName)
        if err != nil {
            return nil, fmt.Errorf("could not update display_name: %w", err)
        }
    }

    return &domain.AuthResponse{
        Success: true,
        Message: "Google user profile updated successfully",
        UserID:  userID,
    }, nil
}

// Implementación correcta de GetUserByEmail
func (s *AuthServiceImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
    user, _, err := s.repo.FindUserByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// Implementación de UpdateGoogleUserProfileByEmail
func (s *AuthServiceImpl) UpdateGoogleUserProfileByEmail(ctx context.Context, email, displayName string) (*domain.AuthResponse, error) {
    // Obtener el usuario por email
    user, err := s.repo.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    // Verificar que el auth_type sea google
    if user.AuthType != "google" {
        return nil, errors.New("only google users can be updated with this endpoint")
    }

    // Actualizar solo el display_name
    if displayName != "" {
        err = s.repo.UpdateDisplayName(ctx, user.ID, displayName)
        if err != nil {
            return nil, fmt.Errorf("could not update display_name: %w", err)
        }
    }

    return &domain.AuthResponse{
        Success: true,
        Message: "Google user profile updated successfully",
        UserID:  user.ID,
    }, nil
}
func (s *AuthServiceImpl) UpdateUserProfile(ctx context.Context, userID int64, email, username, displayName, authType string) (*domain.AuthResponse, error) {
    // Obtener el usuario actual para verificar el auth_type
    user, err := s.repo.GetUserByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    // Verificar que el auth_type coincida
    if user.AuthType != authType {
        return nil, errors.New("auth_type mismatch")
    }

    // Actualizar campos comunes
    if displayName != "" {
        err = s.repo.UpdateDisplayName(ctx, userID, displayName)
        if err != nil {
            return nil, fmt.Errorf("could not update display_name: %w", err)
        }
    }

    // Actualizar campos específicos para email
    if authType == "email" {
        if email != "" {
            err = s.repo.UpdateUserEmailById(ctx, userID, email)
            if err != nil {
                return nil, fmt.Errorf("could not update email: %w", err)
            }
        }

        if username != "" {
            err = s.repo.UpdateUsername(ctx, userID, username)
            if err != nil {
                return nil, fmt.Errorf("could not update username: %w", err)
            }
        }
    }

    return &domain.AuthResponse{
        Success: true,
        Message: "User profile updated successfully",
        UserID:  userID,
    }, nil
}