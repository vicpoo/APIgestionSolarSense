//api/src/login/application/get_auth_usecase.go

package application

import (
	"context"
	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
)

type GetAuthUseCase struct {
	repo domain.AuthRepository
}

func NewGetAuthUseCase(repo domain.AuthRepository) *GetAuthUseCase {
	return &GetAuthUseCase{repo: repo}
}

func (uc *GetAuthUseCase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, _, err := uc.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}