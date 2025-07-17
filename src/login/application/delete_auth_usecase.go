// api/src/login/application/delete_auth_usecase.go
package application

import (
	"context"
	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
)

type DeleteAuthUseCase struct {
	repo domain.AuthRepository
}

func NewDeleteAuthUseCase(repo domain.AuthRepository) *DeleteAuthUseCase {
	return &DeleteAuthUseCase{repo: repo}
}

func (uc *DeleteAuthUseCase) DeleteUserByEmail(ctx context.Context, email string) error {
	return uc.repo.DeleteUserByEmail(ctx, email)
}