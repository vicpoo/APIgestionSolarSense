// api/src/login/domain/entities.go
package domain


import (
	"time"
	"github.com/golang-jwt/jwt/v4"
)

const JwtSecret = "d3c8f2b9e7a14b0f932c0d7d9a7e4f5d6c1a2e8b9f3c4a6e7b1d0f4c9a5e6b7"

type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	AuthType string `json:"auth_type"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(user *User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		AuthType: user.AuthType,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JwtSecret))
}

func ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

type AuthResponse struct {
    Success  bool   `json:"success"`
    Message  string `json:"message,omitempty"`
    Error    string `json:"error,omitempty"`
    Token    string `json:"token,omitempty"`
    AuthType string `json:"auth_type,omitempty"`
    UserID   int64  `json:"user_id,omitempty"`
    Email    string `json:"email,omitempty"`
    Username string `json:"username,omitempty"`
    IsAdmin  bool   `json:"is_admin,omitempty"`
}

type UserCredentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Username string `json:"username,omitempty"`
}

type User struct {
    ID        int64     `json:"id"`
    UID       string    `json:"uid,omitempty"`
    Email     string    `json:"email"`
    Username  string    `json:"username,omitempty"`
    AuthType  string    `json:"auth_type"`
    IsActive  bool      `json:"is_active"`
    IsAdmin   bool      `json:"is_admin"`
    LastLogin time.Time `json:"last_login,omitempty"`
    PhotoURL  string    `json:"photo_url,omitempty"`
    Provider  string    `json:"provider,omitempty"`
    CreatedAt time.Time `json:"created_at,omitempty"`
}