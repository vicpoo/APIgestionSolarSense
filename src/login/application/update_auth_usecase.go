// api/src/login/application/update_auth_usecase.go
package application

import (
	"context"
	"errors"
	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
	"golang.org/x/crypto/bcrypt"
)

type UpdateAuthUseCase struct {
	repo domain.AuthRepository
}

func NewUpdateAuthUseCase(repo domain.AuthRepository) *UpdateAuthUseCase {
	return &UpdateAuthUseCase{repo: repo}
}

func (uc *UpdateAuthUseCase) UpdateUserEmail(ctx context.Context, currentEmail, newEmail string) error {
	return uc.repo.UpdateUserEmail(ctx, currentEmail, newEmail)
}

func (uc *UpdateAuthUseCase) UpdatePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error {
	_, passwordHash, err := uc.repo.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("could not hash password")
	}
	
	return uc.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}