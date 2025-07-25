//api/src/login/application/get_auth_usecase.go

package application

import (
	"context"

	"github.com/vicpoo/apigestion-solar-go/src/login/domain"
)

type GetAuthUseCase struct {
	Repo domain.AuthRepository // Cambiado de repo a Repo (exportado)
}

func NewGetAuthUseCase(repo domain.AuthRepository) *GetAuthUseCase {
	return &GetAuthUseCase{Repo: repo} // Actualizado para usar Repo
}


func (uc *GetAuthUseCase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, _, err := uc.Repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}


func (uc *GetAuthUseCase) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return uc.Repo.GetAllUsers(ctx)
}

func (uc *GetAuthUseCase) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	return uc.Repo.GetUserByID(ctx, userID)
}

func (uc *GetAuthUseCase) GetUserMembershipType(ctx context.Context, userID int64) (string, error) {
	return uc.Repo.GetUserMembershipType(ctx, userID)
}
