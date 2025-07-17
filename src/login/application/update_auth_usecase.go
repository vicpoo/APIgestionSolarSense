// api/src/login/application/update_auth_usecase.go
package application

import (
	"context"
	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
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