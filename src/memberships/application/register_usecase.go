// src/memberships/application/register_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
    "golang.org/x/crypto/bcrypt"
)

type RegisterUseCase struct {
    repo domain.MembershipRepository
}

func NewRegisterUseCase(repo domain.MembershipRepository) *RegisterUseCase {
    return &RegisterUseCase{repo: repo}
}

func (uc *RegisterUseCase) RegisterUser(ctx context.Context, email, username, password string) (int64, error) {
    if email == "" || username == "" || password == "" {
        return 0, domain.ErrInvalidCredentials
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return 0, domain.ErrPasswordHashing
    }

    return uc.repo.RegisterUser(ctx, email, username, string(hashedPassword))
}