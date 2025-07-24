//api/src/login/application/auth_service.go

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
	if creds.Email == "" || creds.Password == "" || creds.Username == "" {
		return nil, errors.New("email, password and username are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("could not hash password")
	}

	_, err = s.repo.CreateUserWithEmail(ctx, creds.Email, creds.Username, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Success: true,
		Message: "User registered successfully",
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

    return &domain.AuthResponse{
        Success: true,
        Message: "Login successful",
        // El token ahora se genera en el controlador
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
        // El token ahora se genera en el controlador
    }, nil
}
func decodeTokenWithoutVerification(idToken string) (map[string]interface{}, error) {
	// Dividir el token JWT en sus partes
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decodificar la parte de los claims (payload)
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token claims")
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims")
	}

	// Extraer datos del usuario
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