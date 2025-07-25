//src/login/infrastructure/jwt_service.go
package infrastructure

import (
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/vicpoo/apigestion-solar-go/src/core"
    "github.com/vicpoo/apigestion-solar-go/src/login/domain"
)

const tokenDuration = 48 * time.Hour

func GenerateJWTToken(user *domain.User) (string, error) {
    claims := core.JWTClaims{
        UserID:   user.ID,
        Email:    user.Email,
        Username: user.Username,
        AuthType: user.AuthType,
        IsAdmin:  user.IsAdmin,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(core.JwtSecret))
}

func ValidateJWTToken(tokenString string) (*core.JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &core.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(core.JwtSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*core.JWTClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrInvalidKey
}