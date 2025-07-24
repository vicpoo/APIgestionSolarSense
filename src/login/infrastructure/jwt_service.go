//src/login/infrastructure/jwt_service.go
package infrastructure

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
)

const (
	jwtSecret = "TuSuperSecretKeySegura123!@#" // Cambia esto en producción
	tokenDuration = 24 * time.Hour
)

type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	AuthType string `json:"auth_type"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(user *domain.User, membershipType string) (string, error) {
	isAdmin := membershipType == "admin"
	
	claims := JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		AuthType: user.AuthType,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}